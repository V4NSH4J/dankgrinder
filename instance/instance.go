// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/dankgrinder/dankgrinder/config"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const fundReqInterval = time.Minute * 10

type Instance struct {
	Client             *discord.Client
	Logger             *logrus.Logger
	ChannelID          string
	WG                 *sync.WaitGroup
	Master             *Instance
	Cluster            []*Instance
	Features           config.Features
	SuspicionAvoidance config.SuspicionAvoidance
	Compat             config.Compat
	Shifts             []config.Shift

	sdlr              *scheduler.Scheduler
	ws                *discord.WSConn
	initialBalance    int
	balance           int

	// auto-trade variables
	totalTradeItems   int
	iteratedItems     int
	currentTradeItems int
	tradeList         string

	lastTradeTime     time.Time
	queuingInstances  int 

	startingTime      time.Time
	lastState         string
	lastBalanceUpdate time.Time
	fatal             chan error
	isClosed          bool
	result            map[string]string
	result2			  string
	result3			  []string
}

func (in *Instance) Start() error {
	if in.Client == nil {
		return fmt.Errorf("no client")
	}
	if in.ChannelID == "" {
		return fmt.Errorf("no channel id")
	}
	if len(in.Shifts) == 0 {
		return fmt.Errorf("no shifts")
	}
	if in.WG == nil {
		return fmt.Errorf("no waitgroup")
	}
	if in.Logger == nil {
		return fmt.Errorf("no logger")
	}
	if in.Master == nil {
		if in.Features.AutoGift.Enable {
			in.Logger.Warnf("nobody to auto-gift to, no master instance available")
		}
		if in.Features.AutoShare.Enable {
			in.Logger.Warnf("nobody to auto-share to, no master instance available")
		}
	}

	in.totalTradeItems = len(in.Features.AutoGift.Items)

	// For now, we assume that in.SuspicionAvoidance, in.Compat and in.Features
	// are correct. They are currently validated in the main function. Ideally,
	// this needs to change in the future.

	in.fatal = make(chan error)
	in.WG.Add(1)
	go func() {
		defer in.WG.Done()
		defer func() {
			in.isClosed = true
		}()
		for {
			for i, shift := range in.Shifts {
				dur := shiftDur(shift)
				in.Logger.WithFields(map[string]interface{}{
					"state":    shift.State,
					"duration": dur,
				}).Infof("starting shift %v", i+1)
				if shift.State == in.lastState {
					in.sleep(dur)
					continue
				}
				in.lastState = shift.State
				if shift.State == config.ShiftStateDormant {
					if in.ws != nil {
						if err := in.ws.Close(); err != nil {
							in.Logger.Errorf("error while closing websocket: %v", err)
						}
					}
					if in.sdlr != nil {
						if err := in.sdlr.Close(); err != nil {
							in.Logger.Errorf("error while closing scheduler: %v", err)
						}
					}
					in.sleep(dur)
					continue
				}
				if err := in.startWS(); err != nil {
					in.Logger.Errorf("instance fatal: error while starting websocket: %v", err)
					return
				}
				if err := in.startSdlr(); err != nil {
					in.Logger.Errorf("instance fatal: error while starting scheduler: %v", err)
					return
				}
				cmds := in.newCmds()
				if in.Features.AutoSell.Enable {
					cmds = append(cmds, in.newAutoSellChain())
				}
				if in.Features.AutoGift.Enable &&
					in.Master != nil &&
					in != in.Master {
					cmds = append(cmds, in.newAutoGiftChain())
				}
				for _, cmd := range cmds {
					in.sdlr.Schedule(cmd)
				}
				in.sleep(dur)
			}
		}
	}()
	if in.Features.AutoShare.Enable && in.Features.AutoShare.Fund && in == in.Master {
		go func() {
			t := time.NewTicker(time.Minute*5 + time.Duration(len(in.Cluster)*in.Compat.Cooldown.Share)*time.Second)
			defer t.Stop()
			for {
				<-t.C
				var totalFunding int
				var fundingCmds []*scheduler.Command
				for _, clusterInstance := range in.Cluster {
					if in == clusterInstance {
						continue
					}
					if !clusterInstance.Features.AutoShare.Enable ||
						clusterInstance.LastBalanceUpdate().IsZero() ||
						clusterInstance.IsClosed() {
						continue
					}
					balance := clusterInstance.Balance()
					if balance >= clusterInstance.Features.AutoShare.MinimumBalance {
						continue
					}
					deficit := clusterInstance.Features.AutoShare.MinimumBalance - balance
					deficit = int(math.Round(float64(deficit) / 0.92)) // Account for 8% tax.
					if totalFunding+deficit > in.balance {
						break
					}
					totalFunding += deficit
					fundingCmds = append(fundingCmds, &scheduler.Command{
						Value:                shareCmdValue(strconv.Itoa(deficit), clusterInstance.Client.User.ID),
						Log:                  "funding",
						Interval:             time.Duration(in.Compat.Cooldown.Share) * time.Second,
						RescheduleAsPriority: true,
					})
				}
				if len(fundingCmds) > 0 {
					in.sdlr.PrioritySchedule(in.newCmdChain(fundingCmds, 0))
				}
			}
		}()
	}
	return nil
}

func (in *Instance) sleep(dur time.Duration) {
	select {
	case err := <-in.fatal:
		in.Logger.Errorf("instance fatal: %v", err)
		runtime.Goexit()
	case <-time.After(dur):
	}
}

func (in *Instance) startSdlr() error {
	in.sdlr = &scheduler.Scheduler{
		Client:             in.Client,
		ChannelID:          in.ChannelID,
		Typing:             &in.SuspicionAvoidance.Typing,
		MessageDelay:       &in.SuspicionAvoidance.MessageDelay,
		Logger:             in.Logger,
		AwaitResumeTimeout: time.Duration(in.Compat.AwaitResponseTimeout) * time.Second,
		FatalHandler: func(ferr error) {
			in.fatal <- fmt.Errorf("scheduler fatal: %v", ferr)
		},
	}
	if err := in.sdlr.Start(); err != nil {
		return fmt.Errorf("error while starting scheduler: %v", err)
	}
	return nil
}

func (in *Instance) startWS() error {
	ws, err := in.Client.NewWSConn(in.router(), in.wsFatalHandler)
	if err != nil {
		return fmt.Errorf("error while starting websocket: %v", err)
	}
	in.ws = ws
	return nil
}

func shiftDur(shift config.Shift) time.Duration {
	if shift.Duration.Base <= 0 {
		return time.Duration(math.MaxInt64)
	}
	d := time.Duration(shift.Duration.Base) * time.Second
	if shift.Duration.Variation > 0 {
		d += time.Duration(rand.Intn(shift.Duration.Variation)) * time.Second
	}
	return d
}

func (in *Instance) wsFatalHandler(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 4004 {
		in.fatal <- fmt.Errorf("websocket closed: authentication failed, try using a new token")
		return
	}
	in.Logger.Errorf("websocket closed: %v", err)

	in.ws, err = in.Client.NewWSConn(in.router(), in.wsFatalHandler)
	if err != nil {
		in.fatal <- fmt.Errorf("error while connecting to websocket: %v", err)
		return
	}
	in.Logger.Infof("reconnected to websocket")
}

func (in *Instance) IsClosed() bool {
	return in.isClosed
}

func (in *Instance) LastBalanceUpdate() time.Time {
	return in.lastBalanceUpdate
}

func (in *Instance) Balance() int {
	return in.balance
}

func (in *Instance) IsActive() bool {
	return in.lastState != config.ShiftStateDormant
}