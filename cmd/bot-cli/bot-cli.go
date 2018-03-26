package main

import (
	"math/rand"
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
			Name:  "AnnoyingRobot",
			Regex: `robot`,
			Function: func(post *model.Post) {
				startOfSentence := []string{"meh,", "whatever,", "who cares,", "boring,", "", "", "", ""}
				selectedStart := startOfSentence[rand.Intn(len(startOfSentence)-1)]
				app.Bot.SendMessageToChannelWithId(post.ChannelId, selectedStart+" the japanese version was better", "")
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

				var pluginsString string
				var count int
				for pluginName, _ := range app.Bot.RegexFunctions {
					count++
					if count == len(app.Bot.RegexFunctions) {
						pluginsString = pluginsString + "**" + pluginName + "** "
					} else {
						pluginsString = pluginsString + "**" + pluginName + "**, "
					}
				}

				app.Bot.SendMessageToChannelWithId(post.ChannelId, "*The following plugins are loaded:* "+pluginsString, "")

			},
		})

		app.Bot.Start(app.Debug)
		app.Bot.OpenShell()
	}
	ui.Run(os.Args)
}
