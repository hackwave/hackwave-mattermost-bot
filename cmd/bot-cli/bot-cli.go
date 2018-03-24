package main

import (
	"os"
	"time"

	cli "github.com/hackwave/cli-framework"
	bot "github.com/hackwave/hackwave-mattermost-bot"
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
			Function: func() {
				app.Bot.SendDebugMessage("hello dawg", "")
			},
		})

		app.Bot.Start()
		app.Bot.OpenShell()
	}
	ui.Run(os.Args)
}
