// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

const (
	begCmdValue           = "pls beg"
	postmemeCmdValue      = "pls pm"
	searchCmdValue        = "pls search"
	highlowCmdValue       = "pls hl"
	fishCmdValue          = "pls fish"
	huntCmdValue          = "pls hunt"
	balanceCheckCmdValue  = "pls bal"
	tidepodCmdValue       = "pls use tidepod"
	acceptTidepodCmdValue = "y"
	buyBaseCmdValue       = "pls buy"
	blackjackBaseCmdValue = "pls bj"
	sellBaseCmdValue      = "pls sell"
	shopBaseCmdValue      = "pls shop"
	tradeBaseCmdValue     = "pls trade"
	digCmdValue           = "pls dig"
	workCmdValue          = "pls work"
	triviaCmdValue        = "pls trivia"
	crimeCmdValue         = "pls crime"
	scratchBaseCmdValue   = "pls scratch"
	guessCmdValue         = "pls gtn"
)

func blackjackCmdValue(amount string) string {
	return fmt.Sprintf("%v %v", blackjackBaseCmdValue, amount)
}
func scratchCmdValue(amount string) string {
	return fmt.Sprintf("%v %v", scratchBaseCmdValue, amount)
}
func buyCmdValue(amount, item string) string {
	return fmt.Sprintf("%v %v %v", buyBaseCmdValue, item, amount)
}

func sellCmdValue(amount, item string) string {
	return fmt.Sprintf("%v %v %v", sellBaseCmdValue, item, amount)
}

func shopCmdValue(item string) string {
	return fmt.Sprintf("%v %v", shopBaseCmdValue, item)
}

func tradeCmdValue(itemlist, id string) string {
	// itemlist will always have a trailing space so it is not necessary here
	return fmt.Sprintf("%v %v<@%v>", tradeBaseCmdValue, itemlist, id)
}

func tradeItemListValue(amount, item string) string {
	return fmt.Sprintf("%v %v ", amount, item)
}

func shareCmdValue(amount, id string) string {
	return fmt.Sprintf("%v %v <@%v>", tradeBaseCmdValue, amount, id)
}

// commands returns a command pointer slice with all commands that should be
// executed periodically. It contains all commands as configured.
func (in *Instance) newCmds() []*scheduler.Command {
	var cmds []*scheduler.Command
	if in.Features.Commands.Beg {
		cmds = append(cmds, &scheduler.Command{
			Value:    begCmdValue,
			Interval: time.Duration(in.Compat.Cooldown.Beg) * time.Second,
		})
	}
	if in.Features.Commands.Postmeme {
		cmds = append(cmds, &scheduler.Command{
			Value:       postmemeCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Postmeme) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Search {
		cmds = append(cmds, &scheduler.Command{
			Value:       searchCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Search) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Crime {
		cmds = append(cmds, &scheduler.Command{
			Value:       crimeCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Crime) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Highlow {
		cmds = append(cmds, &scheduler.Command{
			Value:       highlowCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Highlow) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Fish {
		cmds = append(cmds, &scheduler.Command{
			Value:       fishCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Fish) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Hunt {
		cmds = append(cmds, &scheduler.Command{
			Value:       huntCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Hunt) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Guess {
		cmds = append(cmds, &scheduler.Command{
			Value:       guessCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Guess) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.BalanceCheck.Enable {
		cmds = append(cmds, &scheduler.Command{
			Value:    balanceCheckCmdValue,
			Interval: time.Duration(in.Features.BalanceCheck.Interval) * time.Second,
		})
	}
	if in.Features.AutoTidepod.Enable {
		cmds = append(cmds, &scheduler.Command{
			Value:       tidepodCmdValue,
			Interval:    time.Duration(in.Features.AutoTidepod.Interval) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Dig {
		cmds = append(cmds, &scheduler.Command{
			Value:       digCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Dig) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.AutoBlackjack.Enable {
		cmds = append(cmds, in.newAutoBlackjackCmd())
	}
	if in.Features.Scratch.Enable {
		cmds = append(cmds, in.newScratchCmd())
	}
	if in.Features.Commands.Work {
		cmds = append(cmds, &scheduler.Command{
			Value:       workCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Work) * time.Second,
			AwaitResume: true,
		})
	}
	if in.Features.Commands.Trivia {
		cmds = append(cmds, &scheduler.Command{
			Value:       triviaCmdValue,
			Interval:    time.Duration(in.Compat.Cooldown.Trivia) * time.Second,
			AwaitResume: true,
		})
	}

	for _, cmd := range in.Features.CustomCommands {
		// cmd.Value and cmd.Amount are not checked for correct values here
		// because they were checked when the application started using
		// cfg.Validate().
		cmds = append(cmds, &scheduler.Command{
			Value:    cmd.Value,
			Interval: time.Duration(cmd.Interval) * time.Second,
			Amount:   uint(cmd.Amount),
			CondFunc: func() bool {
				return cmd.PauseBelowBalance == 0 || in.balance >= cmd.PauseBelowBalance
			},
		})
	}
	return cmds
}

func (in *Instance) newAutoSellChain() *scheduler.Command {
	var cmds []*scheduler.Command
	for _, item := range in.Features.AutoSell.Items {
		cmds = append(cmds, &scheduler.Command{
			Value:    sellCmdValue("max", item),
			Interval: time.Duration(in.Compat.Cooldown.Sell) * time.Second,
		})
	}
	return in.newCmdChain(
		cmds,
		time.Duration(in.Features.AutoSell.Interval)*time.Second,
	)
}

func (in *Instance) newAutoGiftChain() *scheduler.Command {
	var cmds []*scheduler.Command
	for _, item := range in.Features.AutoGift.Items {
		cmds = append(cmds, &scheduler.Command{
			Value:       shopCmdValue(item),
			Interval:    time.Duration(in.Compat.Cooldown.Gift) * time.Second,
			AwaitResume: true,
		})
	}
	return in.newCmdChain(
		cmds,
		time.Duration(in.Features.AutoGift.Interval)*time.Second,
	)
}

func (in *Instance) newAutoBlackjackCmd() *scheduler.Command {
	cmd := &scheduler.Command{
		Value:    blackjackCmdValue(strconv.Itoa(in.Features.AutoBlackjack.Amount)),
		Interval: time.Duration(in.Compat.Cooldown.Blackjack) * time.Second,
		CondFunc: func() bool {
			correctBalance := in.Features.AutoBlackjack.PauseBelowBalance == 0 || in.balance >= in.Features.AutoBlackjack.PauseBelowBalance
			return correctBalance && in.balance < in.Features.AutoBlackjack.PauseAboveBalance
		},
		AwaitResume:          true,
		RescheduleAsPriority: in.Features.AutoBlackjack.Priority,
	}
	if in.Features.AutoBlackjack.Amount == 0 {
		cmd.Value = blackjackCmdValue("max")
	}
	return cmd
}
func (in *Instance) newScratchCmd() *scheduler.Command {
	cmd := &scheduler.Command{
		Value:                scratchCmdValue(strconv.Itoa(in.Features.Scratch.Amount)),
		Interval:             time.Duration(in.Compat.Cooldown.Scratch) * time.Second,
		AwaitResume:          true,
		RescheduleAsPriority: in.Features.Scratch.Priority,
	}
	if in.Features.Scratch.Amount == 0 {
		cmd.Value = scratchCmdValue("max")
	}
	return cmd
}

func (in *Instance) newCmdChain(cmds []*scheduler.Command, chainInterval time.Duration) *scheduler.Command {
	for i := 0; i < len(cmds); i++ {
		if i != 0 {
			cmds[i-1].Next = cmds[i]
		}
		if i == len(cmds)-1 {
			cmds[i].Next = cmds[0]
			cmds[i].Interval = chainInterval
		}
	}
	return cmds[0]
}
