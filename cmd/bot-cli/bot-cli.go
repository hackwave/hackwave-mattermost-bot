package main

import (
	"fmt"

	bot "github.com/hackwave/hackwave-mattermost-bot"
)

func main() {
	app := bot.Application{
		Name:    "Hackwave Bot",
		Version: "v0.1.0",
	}

	fmt.Println(app.Name, ":", app.Version)
	fmt.Println("==============================")
	app = app.LoadConfiguration("./bot.yaml")
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
