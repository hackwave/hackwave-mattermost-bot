package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	cli "github.com/hackwave/cli-framework"
	bot "github.com/hackwave/hackwave-mattermost-bot"
	"github.com/mattermost/mattermost-server/model"
)

func main() {
	app := bot.Application{
		Name: "Hackwave Bot",
		Version: bot.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}
	ui := cli.NewApp()
	ui.Name = app.Name
	ui.Version = app.Version.ToString()
	ui.Compiled = time.Now()
	ui.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "./bot.yaml",
			Usage:  "Specify the location of the YAML config file `FILE`",
			EnvVar: "HACKBOT_CONFIG_PATH",
		},
	}
	ui.Action = func(c *cli.Context) {
		app = app.Init(c.String("config"))
		app.Bot.Server = app.Bot.Server.Connect(app.Bot.Email, app.Bot.Password)

		app.Bot = app.Bot.RegisterHook(bot.RegexFunction{
			Name:  "HelloResponder",
			Regex: `(?:^|\W)hello(?:$|\W)`,
			Function: func(post *model.Post) {
				app.Bot.SendDebugMessage("hello catfood, you send the following post :"+post.Message, "")
			},
		})

		app.Bot = app.Bot.RegisterHook(bot.RegexFunction{
			Name:  "DiceRoller",
			Regex: `^roll (|\d{1,2})(d|D|)\d{1,3}$`,
			Function: func(post *model.Post) {
				regex := regexp.MustCompile(" \\d{1,2}")
				numberOfDiceMatch := regex.FindStringSubmatch(post.Message)

				fmt.Println("len of diceMatch:", len(numberOfDiceMatch))
				var numberOfDice string
				if len(numberOfDiceMatch) > 0 {
					// Number of dice
					fmt.Println("numberOfDiceMatch [0]:", numberOfDiceMatch[0])
					numberOfDice = numberOfDiceMatch[0]
				}

				regex = regexp.MustCompile("d\\d{1,3}")
				diceMatch := regex.FindStringSubmatch(post.Message)

				fmt.Println("len of diceMatch:", len(diceMatch))
				var dice string
				if len(diceMatch) > 0 {
					// Number of dice
					fmt.Println("diceMatch [0]:", diceMatch[0])
					fmt.Println("diceMatch [0] without d:", diceMatch[0][1:(len(diceMatch[0])-1)])
					dice = diceMatch[0][1:(len(diceMatch[0]) - 1)]
				}

				app.Bot.SendDebugMessage("hello dogfood, you seem to be trying to roll dice :"+post.Message, "")
				app.Bot.SendDebugMessage("  number of dice to roll :"+numberOfDice, "")
				app.Bot.SendDebugMessage("  type of dice to roll d(sides):"+dice, "")
			},
		})

		app.Bot.Start()
		app.Bot.OpenShell()
	}
	ui.Run(os.Args)
}
