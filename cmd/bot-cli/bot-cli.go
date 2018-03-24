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
	fmt.Println("[SERVER] Connecting to server with bot named: ", app.Bot.Username)
	app.HandleSignals()
	app.Bot.Server = app.Bot.Server.Connect(app.Bot.Email, app.Bot.Password)
	app.Bot.UpdateServerProfile()
	app.Bot.SendDebugMessage("[BOT] "+app.Bot.Username+" in DEBUG MODE has joined the channel.", "")

	go func() {
		for {
			select {
			case response := <-app.Bot.Server.WSClient.EventChannel:
				app.Bot.HandleWebSocketResponse(response)
			}
		}
	}()
	// Hold open function forever
	select {}
}
