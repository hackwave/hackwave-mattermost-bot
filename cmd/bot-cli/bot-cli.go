package main

import (
	"os"
	"time"

	cli "github.com/hackwave/cli-framework"
	bot "github.com/hackwave/hackwave-mattermost-bot"
	"github.com/hackwave/hackwave-mattermost-bot/plugins/dice"
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
			Name:  dice.PLUGIN_NAME,
			Regex: dice.PLUGIN_REGEX,
			Function: func(post *model.Post) {
				diceResult := dice.ParseDiceCommand(post.Message)
				if app.Debug {
					app.Bot.SendDebugMessage(diceResult, "")
				} else {
					app.Bot.SendMessageToChannelWithId(post.ChannelId, diceResult, "")
				}
			},
		})

		app.Bot = app.Bot.RegisterHook(bot.RegexFunction{
			Name:  "Help",
			Regex: `^help$`,
			Function: func(post *model.Post) {
				app.Bot.SendMessageToChannelWithId(post.ChannelId, "**Hello, my name is "+app.Bot.Username+" and I'm a often malfunctioning bot**", "")
				app.Bot.SendMessageToChannelWithId(post.ChannelId, " Below is a list of commands I respond to:", "")
				app.Bot.SendMessageToChannelWithId(post.ChannelId, "", "")
				app.Bot.SendMessageToChannelWithId(post.ChannelId, dice.PLUGIN_HELP_COMMAND, "")
				app.Bot.SendMessageToChannelWithId(post.ChannelId, dice.PLUGIN_HELP_TEXT, "")
				app.Bot.SendMessageToChannelWithId(post.ChannelId, "", "")
			},
		})

		app.Bot.Start()
		app.Bot.OpenShell()
	}
	ui.Run(os.Args)
}
