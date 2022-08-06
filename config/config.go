// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ShiftStateActive  = "active"
	ShiftStateDormant = "dormant"
)

type Config struct {
	Clusters           map[string]Cluster `yaml:"clusters"`
	Shifts             []Shift            `yaml:"shifts"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
}

type Cluster struct {
	Master    Instance   `yaml:"master"`
	Instances []Instance `yaml:"instances"`
}

type Instance struct {
	Token              string             `yaml:"token"`
	ChannelID          string             `yaml:"channel_id"`
	Features           Features           `yaml:"-"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"-"`
	Shifts             []Shift            `yaml:"-"`
}

type override struct {
	Clusters map[string]struct {
		Master struct {
			Shifts             yaml.Node `yaml:"shifts"`
			Features           yaml.Node `yaml:"features"`
			SuspicionAvoidance yaml.Node `yaml:"suspicion_avoidance"`
		} `yaml:"master"`
		Instances []struct {
			Shifts             yaml.Node `yaml:"shifts"`
			Features           yaml.Node `yaml:"features"`
			SuspicionAvoidance yaml.Node `yaml:"suspicion_avoidance"`
		} `yaml:"instances"`
	} `yaml:"clusters"`
}

type Compat struct {
	AllowedSearches      []string `yaml:"allowed_searches"`
	AllowedCrimes        []string `yaml:"allowed_crimes"`
	Cooldown             Cooldown `yaml:"cooldown"`
	AwaitResponseTimeout int      `yaml:"await_response_timeout"`
	AllowedScrambles     []string `yaml:"allowed_scrambles"`
	DigCancel            []string `yaml:"dig_cancel"`
	AllowedFTB           []string `yaml:"allowed_ftb"`
	WorkCancel           []string `yaml:"work_cancel"`
	AllowedScramblesWork []string `yaml:"allowed_scrambles_work"`
	AllowedHangman       []string `yaml:"allowed_hangman"`
	AllowedScramblesFish []string `yaml:"allowed_scrambles_fish"`
	AllowedFishFTB       []string `yaml:"allowed_fish_ftb"`
	FishCancel           []string `yaml:"fish_cancel"`
	SearchMode           int      `yaml:"search_mode"`
	CrimeMode            int      `yaml:"crime_mode"`
}

type Cooldown struct {
	Beg       int `yaml:"beg"`
	Fish      int `yaml:"fish"`
	Hunt      int `yaml:"hunt"`
	Postmeme  int `yaml:"postmeme"`
	Search    int `yaml:"search"`
	Highlow   int `yaml:"highlow"`
	Blackjack int `yaml:"blackjack"`
	Sell      int `yaml:"sell"`
	Gift      int `yaml:"gift"`
	Share     int `yaml:"share"`
	Dig       int `yaml:"dig"`
	Work      int `yaml:"work"`
	Trivia    int `yaml:"trivia"`
	Crime     int `yaml:"crime"`
	Scratch   int `yaml:"scratch"`
	Guess     int `yaml:"guess"`
}

type Features struct {
	Commands           Commands        `yaml:"commands"`
	CustomCommands     []CustomCommand `yaml:"custom_commands"`
	AutoBuy            AutoBuy         `yaml:"auto_buy"`
	AutoSell           AutoSell        `yaml:"auto_sell"`
	AutoGift           AutoGift        `yaml:"auto_gift"`
	AutoBlackjack      AutoBlackjack   `yaml:"auto_blackjack"`
	AutoShare          AutoShare       `yaml:"auto_share"`
	AutoTidepod        AutoTidepod     `yaml:"auto_tidepod"`
	BalanceCheck       BalanceCheck    `yaml:"balance_check"`
	LogToFile          bool            `yaml:"log_to_file"`
	VerboseLogToStdout bool            `yaml:"verbose_log_to_stdout"`
	Debug              bool            `yaml:"debug"`
	Scratch            Scratch         `yaml:"scratch"`
	Trade          	   Trade           `yaml:"trade"`
}
type Scratch struct {
	Enable   bool `yaml:"enable"`
	Amount   int  `yaml:"amount"`
	Priority bool `yaml:"priority"`
}

type BalanceCheck struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

type AutoTidepod struct {
	Enable              bool `yaml:"enable"`
	Interval            int  `yaml:"interval"`
	BuyLifesaverOnDeath bool `yaml:"buy_lifesaver_on_death"`
}

type AutoBlackjack struct {
	Enable            bool                         `yaml:"enable"`
	Priority          bool                         `yaml:"priority"`
	Amount            int                          `yaml:"amount"`
	PauseBelowBalance int                          `yaml:"pause_below_balance"`
	PauseAboveBalance int 						   `yaml:"pause_above_balance"`
	LogicTable        map[string]map[string]string `yaml:"logic_table"`
}

type AutoShare struct {
	Enable         bool `yaml:"enable"`
	Fund           bool `yaml:"fund"`
	MaximumBalance int  `yaml:"maximum_balance"`
	MinimumBalance int  `yaml:"minimum_balance"`
}

type AutoGift struct {
	Enable   bool     `yaml:"enable"`
	Interval int      `yaml:"interval"`
	Items    []string `yaml:"items"`
}

type Trade struct {
	MaxItems  int     `yaml:"max_items"`
	Delay     int     `yaml:"delay"`
}

type CustomCommand struct {
	Value             string `yaml:"value"`
	Interval          int    `yaml:"interval"`
	Amount            int    `yaml:"amount"`
	PauseBelowBalance int    `yaml:"pause_below_balance"`
}

type AutoBuy struct {
	FishingPole  bool `yaml:"fishing_pole"`
	HuntingRifle bool `yaml:"hunting_rifle"`
	Laptop       bool `yaml:"laptop"`
	Shovel       bool `yaml:"shovel"`
}

type AutoSell struct {
	Enable   bool     `yaml:"enable"`
	Interval int      `yaml:"interval"`
	Items    []string `yaml:"items"`
}

type Commands struct {
	Beg      bool `yaml:"beg"`
	Postmeme bool `yaml:"postmeme"`
	Search   bool `yaml:"search"`
	Highlow  bool `yaml:"highlow"`
	Fish     bool `yaml:"fish"`
	Hunt     bool `yaml:"hunt"`
	Dig      bool `yaml:"dig"`
	Work     bool `yaml:"work"`
	Trivia   bool `yaml:"trivia"`
	Crime    bool `yaml:"crime"`
	Guess    bool `yaml:"guess"`
}

type SuspicionAvoidance struct {
	Typing           Typing           `yaml:"typing"`
	MessageDelay     MessageDelay     `yaml:"message_delay"`
	ButtonPressDelay ButtonPressDelay `yaml:"button_press"`
}

type ButtonPressDelay struct {
	Base      int `yaml:"base"`
	Variation int `yaml:"variation"`
}

type Typing struct {
	Base      int `yaml:"base"`      // A base duration in milliseconds.
	Speed     int `yaml:"speed"`     // Speed in keystrokes per minute.
	Variation int `yaml:"variation"` // A random value in milliseconds from [0,n) added to the base.
}

type MessageDelay struct {
	Base      int `yaml:"base"`      // A base duration in milliseconds.
	Variation int `yaml:"variation"` // A random value in milliseconds from [0,n) added to the base.
}


// Shift indicates an application state (active or dormant) for a duration.
type Shift struct {
	State    string   `yaml:"state"`
	Duration Duration `yaml:"duration"`
}

// Duration is not related to a time.Duration. It is a structure used in a Shift
// type.
type Duration struct {
	Base      int `yaml:"base"`      // A base duration in seconds.
	Variation int `yaml:"variation"` // A random value in seconds from [0,n) added to the base.
}

// Load loads the config from the expected path.
func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}
	defer f.Close()

	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while decoding config: %v", err)
	}

	if _, err = f.Seek(0, 0); err != nil {
		return Config{}, fmt.Errorf("error while seeking back to beginning of file: %v", err)
	}
	var ovr override
	if err = yaml.NewDecoder(f).Decode(&ovr); err != nil {
		return Config{}, fmt.Errorf("error while decoding config override: %v", err)
	}

	if len(cfg.Clusters) != len(ovr.Clusters) {
		panic("amount of instances not equal to the amount of override configs")
	}

	for ck, cluster := range ovr.Clusters {
		for i, instance := range append(cluster.Instances, cluster.Master) {
			features, suspicionAvoidance, shifts := cfg.Features, cfg.SuspicionAvoidance, cfg.Shifts
			if instance.Features.Kind != 0 {
				if err = instance.Features.Decode(&features); err != nil {
					return Config{}, fmt.Errorf(
						"clusters[%v].instances[%v].features error while decoding: %v",
						ck,
						i,
						err,
					)
				}
			}
			if instance.SuspicionAvoidance.Kind != 0 {
				if err = instance.SuspicionAvoidance.Decode(&suspicionAvoidance); err != nil {
					return Config{}, fmt.Errorf(
						"clusters[%v].instances[%v].suspicion_avoidance error while decoding: %v",
						ck,
						i,
						err,
					)
				}
			}
			if instance.Shifts.Kind != 0 {
				if err = instance.Shifts.Decode(&shifts); err != nil {
					return Config{}, fmt.Errorf(
						"clusters[%v].instances[%v].shifts error while decoding: %v",
						ck,
						i,
						err,
					)
				}
			}
			if i == len(cluster.Instances) { // If this is the master instance
				// Workaround. If done similar to the else case, a cannot assign
				// compiler error is given.
				cluster := cfg.Clusters[ck]
				cluster.Master.Features = features
				cluster.Master.SuspicionAvoidance = suspicionAvoidance
				cluster.Master.Shifts = shifts
				cfg.Clusters[ck] = cluster
			} else {
				cfg.Clusters[ck].Instances[i].Features = features
				cfg.Clusters[ck].Instances[i].SuspicionAvoidance = suspicionAvoidance
				cfg.Clusters[ck].Instances[i].Shifts = shifts
			}
		}
	}

	return cfg, nil
}
