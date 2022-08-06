// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"strings"
	"strconv"
	"time"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) gift(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	
	// increment items iterated
	in.iteratedItems++
	
	if trigger == nil || !strings.Contains(trigger.Value, shopBaseCmdValue) {
		return
	}
	if in == in.Master {
		in.sdlr.Resume()
		return
	}

	if !in.Master.IsActive() {
		in.Logger.Errorf("gift failed - master is dormant")
		in.sdlr.Resume()
		return
	}

	if !exp.gift.Match([]byte(msg.Embeds[0].Title)) || !exp.shop.Match([]byte(trigger.Value)) {
		in.sdlr.Resume()
		return
	}

	amount := strings.Replace(exp.gift.FindStringSubmatch(msg.Embeds[0].Title)[1], ",", "", -1)
	item := exp.shop.FindStringSubmatch(trigger.Value)[1]

	if amount == "0" {
		in.Logger.Infof("no items of type %v", item)
		in.sdlr.Resume()
		return
	}

	giftChainEnd := in.iteratedItems == in.totalTradeItems

	// append items to the list in the format for trade command 
	in.tradeList += tradeItemListValue(amount, item)

  	// store amount of items in current item list 
	in.currentTradeItems++

	// If less then max amount of items and not at the end of gift chain wait until
	// later iteration to send 
	if in.currentTradeItems < in.Features.Trade.MaxItems && !giftChainEnd {
		in.Logger.Infof("added %v %v to gift list", amount, item)
		in.sdlr.Resume()
		return
	} 
	
	if giftChainEnd {
		// reset counter when iteration is completed
		in.iteratedItems = 0
	}
	
	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: tradeCmdValue(in.tradeList, in.Master.Client.User.ID),
		Log:   "gifting items - starting trade",
		AwaitResume: true,
	})
	
	in.tradeList = ""
	in.currentTradeItems = 0
}

func (in *Instance) confirmTrade(msg discord.Message) {
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Actionrow: 1,
		Button: 2,
		Message: msg,
		Log: "gifting items - accepting trade as sender",
		AwaitResume: true,
	})
}

func (in *Instance) confirmTradeAsMaster(msg discord.Message) {
	if !in.Master.IsActive() {
		// Ensure that the master is active before trying to click button
		in.Logger.Errorf("master is dormant, sleeping for 5mins to wait for trade to timeout")
		time.Sleep(5 * 60 * 1000 * time.Millisecond)
		in.sdlr.Resume()
		return
	}
	
	in.Master.queuingInstances++;
	delay := int64(in.Master.Features.Trade.Delay)
	// If master has traded less then delay time
	// wait until delay time is reached and create a 'queue' of instances
	elapsedTime := int64(time.Since(in.Master.lastTradeTime))

	// convert to ms
	elapsedTime /= (int64(time.Millisecond)/int64(time.Nanosecond))

	if elapsedTime < delay {
		// sleep until elapsedTime == delay
		// if there is other instances waiting in queue, wait for them as well
		sleepTime := delay*(int64(in.Master.queuingInstances) - 1) + delay - elapsedTime
		in.Logger.Infof("waiting for turn for master to accept: sleeping for %vms (%v in queue)", sleepTime, in.Master.queuingInstances)
		// convert sleep time back to nanoseconds
		time.Sleep(time.Duration(sleepTime * (int64(time.Millisecond)/int64(time.Nanosecond))))
	}

	in.Master.queuingInstances--;

	// After waiting for this instance's turn to trade, check if master is active again
	if !in.Master.IsActive() {
		in.Logger.Errorf("master is dormant, sleeping for 5mins to wait for trade to timeout")
		time.Sleep(5 * 60 * 1000 * time.Millisecond)
		in.sdlr.Resume()
		return
	}

	// If trade request mentioning master is sent, priority schedule a click on accept
	in.Master.sdlr.PrioritySchedule(&scheduler.Command{
		Actionrow: 1,
		Button: 2,
		Message: msg,
		Log: "gifting items - accepting trade as master",
	})

	// update the last time that master has traded
	in.Master.lastTradeTime = time.Now()

	// now that master has accepted, resume 
	in.sdlr.Resume()
}

func (in *Instance) shareWithTax(msg discord.Message) {
	match := exp.insufficientCoins.FindStringSubmatch(msg.Content)

	currentAmount, err := strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	if err != nil {
		in.Logger.Errorf("error while reading sending amount: %v", err)
		in.sdlr.Resume()
		return
	}

	attemptedAmount, err := strconv.Atoi(strings.Replace(match[2], ",", "", -1))

	if err != nil {
		in.Logger.Errorf("error while reading sending amount: %v", err)
		in.sdlr.Resume()
		return
	}

	// If amount of coins that is attempted to send is more then the amount 
	// in wallet, something went wrong and we log an error 
	if attemptedAmount > currentAmount {
		in.Logger.Errorf("attempted to send more coins then current amount: %v, attempted to send: %v", currentAmount, attemptedAmount)
		in.sdlr.Resume()
		return
	}

	tax, err:= strconv.Atoi(strings.Replace(match[3], ",", "", -1))

	if err != nil {
		in.Logger.Errorf("error while reading tax amount: %v", err)
		in.sdlr.Resume()
		return
	}
	amount := strconv.Itoa(attemptedAmount - tax)

	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: shareCmdValue(amount, in.Master.Client.User.ID),
		Log:   "re-trading coins to account for tax",
		AwaitResume: true,
	})
}