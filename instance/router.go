package instance

import (

	"fmt"
	"math/rand"
	"regexp"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const DMID = "270904126974590976"

var exp = struct {
	search,
	huntEvent,
	hl,
	bal,
	gift,
	shop,
	blackjack,
	blackjackBal,
	digEventScramble,
	digEventRetype,
	digEventFTB,
	workEventReverse,
	workEventRetype,
	workEventScramble,
	workEventSoccer,
	workEventHangman,
	workEventRepeat,
	workEventColor,
	workEventColor2,
	fishEventScramble,
	fishEventFTB,
	fishEventReverse,
	fishEventRetype,
	fishCatch,
	fishCatch2,
	trivia,
	guess,
	guessHint,
	workEventEmoji,
	shopEvent,
	event,
	insufficientCoins *regexp.Regexp
}{
	search:            regexp.MustCompile(`Pick from the list below and type the name in chat\.\s\x60(.+)\x60,\s\x60(.+)\x60,\s\x60(.+)\x60`),
	huntEvent:         regexp.MustCompile(`Dodge the Fireball\n(\s*)+<:Dragon:861390869696741396>\n(\s*)<:FireBall:883714770748964864>\n(\s*):levitate:`),
	hl:                regexp.MustCompile(`I just chose a secret number between 1 and 100.\nIs the secret number \*higher\* or \*lower\* than \*\*(.+)\*\*.`),
	bal:               regexp.MustCompile(`\*\*Wallet\*\*: \x60?⏣?\s?([0-9,]+)\x60?`),
	event:             regexp.MustCompile(`Attack the boss by clicking \x60(.+)\x60`),
	gift:              regexp.MustCompile(`[a-zA-Z\s]* \(([0-9,]+) owned\)`),
	shop:              regexp.MustCompile(`pls shop ([a-zA-Z\s]+)`),
	blackjack:         regexp.MustCompile(`\x60[♥♦♠♣] ([0-9]{1,2}|[JQKA])\x60`),
	blackjackBal:      regexp.MustCompile(`(You now have|You have) (\*\*)?(⏣\s)?(\*\*)?([0-9,]+)(\*\*)?(\sstill)?\.`),
	digEventScramble:  regexp.MustCompile(`Quickly unscramble the word to uncover what's in the dirt! in the next 15\sseconds\s\x60(.+)\x60`),
	digEventRetype:    regexp.MustCompile(`Quickly re-type the phrase to uncover what's in the dirt! in the next 15 seconds\nType\s\x60(.+)\x60`),
	digEventFTB:       regexp.MustCompile(`Quickly guess the missing word to uncover what's in the dirt in the next 15 seconds!\n\x60(.+)\x60`),
	workEventReverse:  regexp.MustCompile(`\*\*Work for (.+)\*\* - Reverse - Type the following word backwards.\n\x60(.+)\x60`),
	workEventRetype:   regexp.MustCompile(`\*\*Work for (.+)\*\* - Retype - Retype the following phrase below.\nType\s\x60(.+)\x60`),
	workEventScramble: regexp.MustCompile(`\*\*Work for (.+)\*\* - Scramble - The following word is scrambled, you need to try and unscramble it to reveal the original word.\n\x60(.+)\x60`),
	workEventSoccer:   regexp.MustCompile(`Hit the ball!\n:goal::goal::goal:\n(\s*):levitate:`),
	workEventHangman:  regexp.MustCompile(`\*\*Work for (.+)\*\* - Hangman - Find the missing __word__ in the following sentence:\n\x60(.+)\x60`),
	workEventRepeat:   regexp.MustCompile(`Remember words order!\n\x60(.+)\x60\n\x60(.+)\x60\n\x60(.+)\x60\n\x60(.+)\x60\n\x60(.+)\x60`),
	workEventColor:    regexp.MustCompile(`\*\*Work for (.+)\*\* - Color Match - Match the color to the selected word.\n<:(.+):[\d]+>\s\x60(.+)\x60\n<:(.+):[\d]+>\s\x60(.+)\x60\n<:(.+):[\d]+>\s\x60(.+)\x60`),
	workEventColor2:   regexp.MustCompile(`What color was next to the word \x60(.+)\x60`),
	fishEventScramble: regexp.MustCompile(`the fish is too strong! Quickly unscramble the word to catch it in the next 15 seconds\n\x60(.+)\x60`),
	fishEventFTB:      regexp.MustCompile(`the fish is too strong! Quickly guess the missing word to catch it in the next 15 seconds!\n\x60(.+)\x60`),
	fishEventReverse:  regexp.MustCompile(`the fish is too strong! Quickly reverse the word to catch it in the next 10 seconds!.\n\x60(.+)\x60`),
	fishEventRetype:   regexp.MustCompile(`the fish is too strong! Quickly re-type the phrase to catch it in the next 15 seconds\nType\s\x60(.+)\x60`),
	trivia:            regexp.MustCompile(`\*\*(.+)\*\*\n\*You have \d\d seconds to answer`),
	guess:             regexp.MustCompile(`not this time, \x60(.+)\x60 attempts left and \x60(.+)\x60 (hint|hints) left.`),
	guessHint:         regexp.MustCompile(`Your last number \(\*\*(.+)\*\*\) was too (.+)\nYou\'ve got \x60(.+)\x60 attempts left and \x60(.+)\x60 (hint|hints) left.`),
	workEventEmoji:    regexp.MustCompile(`\*\*Work for (.+)\*\* - Emoji Match - Look at the emoji closely!\n(.+)`),
	fishCatch:         regexp.MustCompile(`Catch the fish!\n(\s*)<:(.+):[\d]+>\n:bucket::bucket::bucket:`),
	fishCatch2:        regexp.MustCompile(`Catch the fish!\n<:(.+):[\d]+>\n:bucket::bucket::bucket:`),
	shopEvent:         regexp.MustCompile(`What is the \*\*(.+)\*\* of this item?`),
	insufficientCoins: regexp.MustCompile(`You have ([\d,]+) coins, you can't give them ([\d,]+) \(\+ ⏣ ([\d,]+) tax\)`),
}

var numFmt = message.NewPrinter(language.English)

func (in *Instance) pm(msg discord.Message) {
	i := rand.Intn(5)
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Actionrow: 1,
		Button:    i + 1,
		Message:   msg,
		Log:       "Responding to post meme randomly",
	})
}

func (in *Instance) automodBypass(_ discord.Message) {
	in.sdlr.Logger.Errorf("Instance %v is blacklisted, Please remove it from config to prevent a bot ban", in.Client.User)
	in.ws.Close()
	in.sdlr.Close()
}

// clean removes all characters except for ASCII characters [32, 126] (basically
// all keys you would find on a US keyboard).
func clean(s string) string {
	var result string
	allowedChars := regexp.MustCompile(`[\x20-\x7E]`)
	for _, char := range s {
		if allowedChars.MatchString(string(char)) {
			result += string(char)
		}
	}
	return result
}

func (in *Instance)printMessage(msg discord.Message) {
	fmt.Println("%v", msg)

}

func (in *Instance) router() *discord.MessageRouter {
	rtr := &discord.MessageRouter{}
	// hunting.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.huntEvent).
		Handler(in.huntEvent)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.huntEvent).
		EventType(discord.EventNameMessageUpdate).
		Handler(in.huntEvent)

	// When a hunt is completed without any events, Dank Memer will
	// reference the original command. If there are events it will mention. This
	// can therefore be used to differentiate between the two.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		RespondsTo(in.Client.User.ID).
		Handler(in.huntEnd)
	// Automod Ban Bypass
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("stop trying to run commands, you're blacklisted. Do this too much and you'll get a full ban.").
		Handler(in.automodBypass)


	// fishing without event
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		RespondsTo(in.Client.User.ID).
		Handler(in.fishEnd)
	// Catch the fish
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishCatch).
		Handler(in.fishCatch)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishCatch2).
		Handler(in.fishCatch)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		EventType(discord.EventNameMessageUpdate).
		ContentMatchesExp(exp.fishCatch2).
		Handler(in.fishCatch)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishCatch).
		EventType(discord.EventNameMessageUpdate).
		Handler(in.fishCatch)
	// fishing with reversing
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishEventReverse).
		Mentions(in.Client.User.ID).
		Handler(in.fishEventReverse)
	// fishing with scramble
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishEventScramble).
		Mentions(in.Client.User.ID).
		Handler(in.fishEventScramble)
	// fishing with fill blank
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishEventFTB).
		Mentions(in.Client.User.ID).
		Handler(in.fishEventFTB)
	// fishing with retyping
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.fishEventRetype).
		Mentions(in.Client.User.ID).
		Handler(in.fishEventRetype)
	// Scratch
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("You can scratch **").
		Handler(in.scratch)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("You can scratch **").
		EventType(discord.EventNameMessageUpdate).
		Handler(in.scratch)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("You can scratch **0** more fields").
		EventType(discord.EventNameMessageUpdate).
		Handler(in.scratchEnd)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("You abandoned your Scratch-Off, SHAME!").
		EventType(discord.EventNameMessageUpdate).
		Handler(in.scratchEnd)

	//Digging Without Event
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		RespondsTo(in.Client.User.ID).
		Handler(in.digEnd)
	//Digging With Scramble
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.digEventScramble).
		Mentions(in.Client.User.ID).
		Handler(in.digEventScramble)
	//Digging with Retype
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.digEventRetype).
		Mentions(in.Client.User.ID).
		Handler(in.digEventRetype)
	//Digging with Fill in the blanks
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.digEventFTB).
		Mentions(in.Client.User.ID).
		Handler(in.digEventFTB)
	//Working Reversing string
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventReverse).
		RespondsTo(in.Client.User.ID).
		Handler(in.workEventReverse)
	//Working Retyping string
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventRetype).
		RespondsTo(in.Client.User.ID).
		Handler(in.workEventRetype)
	//Working Scramble solve
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventScramble).
		RespondsTo(in.Client.User.ID).
		Handler(in.workEventScramble)
	//Working Soccer
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventSoccer).
		Handler(in.workEventSoccer)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventSoccer).
		EventType(discord.EventNameMessageUpdate).
		Handler(in.workEventSoccer)
	//Working Hangman
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventHangman).
		RespondsTo(in.Client.User.ID).
		Handler(in.workEventHangman)
	//Working Repeat + response
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventRepeat).
		Handler(in.workEventRepeat)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("Click the buttons in correct order!").
		EventType(discord.EventNameMessageUpdate).
		Handler(in.workEventRetype2)
	//Working Emoji + response
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventEmoji).
		Handler(in.workEventEmoji)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("What was the emoji?").
		EventType(discord.EventNameMessageUpdate).
		Handler(in.workEventEmoji2)
	//Working Color + response
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventColor).
		Handler(in.workEventColor)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.workEventColor2).
		EventType(discord.EventNameMessageUpdate).
		Handler(in.workEventColor2)
	//Working Don't have a job
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("You don't currently have a job to work at").
		RespondsTo(in.Client.User.ID).
		Handler(in.WorkEnd)
	//Working Recently resigned
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("You recently resigned from your old job.").
		RespondsTo(in.Client.User.ID).
		Handler(in.WorkEnd)
	// Worked Recently
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("You need to wait").
		RespondsTo(in.Client.User.ID).
		Handler(in.WorkEnd)
	//Work promotion
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("You never fail to amaze me").
		Mentions(in.Client.User.ID).
		Handler(in.workPromotion)
	// Trivia
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		RespondsTo(in.Client.User.ID).
		EmbedContains("seconds to answer").
		Handler(in.trivia)

	// Postmeme.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("Hopefully people will like it and give you some").
		Handler(in.pm)

	// Global events.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(false).
		ContentMatchesExp(exp.event).
		Handler(in.event)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(false).
		ContentContains("Berries and Cream").
		Handler(in.event)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedMatchesExp(exp.shopEvent).
		Handler(in.shopEvent)
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EventType(discord.EventNameMessageUpdate).
		EmbedMatchesExp(exp.shopEvent).
		Handler(in.shopEvent)

	// Search.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("**Where do you want to search?**").
		RespondsTo(in.Client.User.ID).
		Handler(in.search)
	// Gtn
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("You've got 4 attempts to try and guess my random number between").
		RespondsTo(in.Client.User.ID).
		Handler(in.guessTheNumber)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.guess).
		RespondsTo(in.Client.User.ID).
		Handler(in.guessTheNumber)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentMatchesExp(exp.guessHint).
		RespondsTo(in.Client.User.ID).
		Handler(in.guessTheNumber)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("Good stuff, you got the number right. I was thinking").
		RespondsTo(in.Client.User.ID).
		Handler(in.gtnEnd)

	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("Unlucky, you ran out of attempts to guess the number").
		RespondsTo(in.Client.User.ID).
		Handler(in.gtnEnd)

	// Highlow.
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		HasEmbeds(true).
		EmbedContains("I just chose a secret number between 1 and 100.").
		RespondsTo(in.Client.User.ID).
		Handler(in.hl)
	// Crime
	rtr.NewRoute().
		Channel(in.ChannelID).
		Author(DMID).
		ContentContains("**What crime do you want to commit?**").
		RespondsTo(in.Client.User.ID).
		Handler(in.crime)

	// Balance report.
	if in.Features.BalanceCheck.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.balanceCheck)
	}

	// Auto-buy laptop.
	if in.Features.AutoBuy.Laptop {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("oi you need to buy a laptop in the shop to post memes").
			RespondsTo(in.Client.User.ID).
			Handler(in.abLaptop)
	}

	// Auto-buy fishing pole.
	if in.Features.AutoBuy.FishingPole {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a fishing pole").
			RespondsTo(in.Client.User.ID).
			Handler(in.abFishingPole)
	}

	// Auto-buy hunting rifle.
	if in.Features.AutoBuy.HuntingRifle {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a hunting rifle").
			RespondsTo(in.Client.User.ID).
			Handler(in.abHuntingRifle)
	}

	// Auto-buy shovel.
	if in.Features.AutoBuy.Shovel {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't have a shovel").
			RespondsTo(in.Client.User.ID).
			Handler(in.abShovel)
	}

	// Auto-gift
	if in.Features.AutoGift.Enable &&
		in.Master != nil &&
		in != in.Master {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			Handler(in.gift)
	}

	// Auto Accept trade as both instance and master
	if (in.Features.AutoShare.Enable || in.Features.AutoGift.Enable) &&
	in.Master != nil &&
	in != in.Master {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			RespondsTo(in.Client.User.ID).
			EmbedContains("Continue trade?").
			Handler(in.confirmTrade)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			RespondsTo(in.Client.User.ID).
			Mentions(in.Master.Client.User.ID).
			EmbedContains("Continue trade?").
			Handler(in.confirmTradeAsMaster)
	}

	// Auto-share (detect if share fails and retry)
	if in.Features.AutoShare.Enable && 
	in.Master != nil &&
	in != in.Master {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			RespondsTo(in.Client.User.ID).
			ContentMatchesExp(exp.insufficientCoins).
			Handler(in.shareWithTax)
	}

	// Auto-tidepod
	if in.Features.AutoTidepod.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("There's a high chance you'll injure yourself from the tidepod").
			Handler(in.tidepod)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("Eating a tidepod is just dumb and stupid.").
			Handler(in.tidepodDeath)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You lost **all of your coins**.").
			Handler(in.tidepodDeath)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			ContentContains("You don't own this item??").
			Handler(in.abTidepod)

	}

	// Auto-blackjack
	if in.Features.AutoBlackjack.Enable {
		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			AuthorNameContains("blackjack game").
			Handler(in.blackjack)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			EventType(discord.EventNameMessageUpdate).
			ContentContains("Type `h` to **hit**, type `s` to **stand**, or type `e` to **end** the game.").
			Handler(in.blackjack)

		rtr.NewRoute().
			Channel(in.ChannelID).
			Author(DMID).
			HasEmbeds(true).
			EventType(discord.EventNameMessageUpdate).
			Handler(in.blackjackEnd)
	}

	return rtr
}
