package main

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
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
	fmt.Println("[SERVER] Connecting to server with bot named: ", app.Bot.Username)
	app.HandleSignals()
	app.Bot.Server = app.Bot.Server.Connect(app.Bot.Email, app.Bot.Password)
	app.Bot.UpdateServerProfile()
	app.Bot.SendDebugMessage("[BOT] "+app.Bot.Username+" in DEBUG MODE has joined the channel.", "")

	go func() {
		for {
			select {
			case event := <-app.Bot.Server.WSClient.EventChannel:
				app.Bot.HandleWebSocketResponse(event)
			}
		}
	}()

	shell := ishell.New()
	shell.Println(">>> Opening Chat Interface")
	shell.Println(">>>   Enter manual chat messages that will be posted by the bot")
	shell.Println(">>>   send {message you want to send as the bot}")
	shell.AddCmd(&ishell.Cmd{
		Name: "send",
		Help: "send message",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				app.Bot.SendDebugMessage(strings.Join(c.Args, " "), "")
			} else {
				fmt.Println("[Error] No message provided, nothing sent.")
			}
		},
	})
	shell.Run()
}
