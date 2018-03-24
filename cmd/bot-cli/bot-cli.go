package main

import (
	"fmt"

	bot "github.com/hackwave/hackwave-mattermost-bot"
	"github.com/mattermost/mattermost-server/model"
)

func main() {
	app := bot.Application{
		Name:    "Hackwave Bot",
		Version: "v0.1.0",
	}

	fmt.Println(app.Name, ":", app.Version)
	fmt.Println("==============================")

	app = app.LoadConfiguration("./bot.yaml")
	fmt.Println("Did yaml file load?")
	fmt.Println("app.Bot.Email is: ", app.Bot.Email)
	fmt.Println("app.Bot.Username is: ", app.Bot.Username)

	//SetupGracefulShutdown()
	app.Bot.Server.HTTPClient = model.NewAPIv4Client(app.Bot.Server.ServerAddress(HTTPServer))

	//MakeSureServerIsRunning()
	app.Bot.Server.Ping()

	//LoginAsTheBotUser()
	app.Bot.Server.Login(app.Bot.Email, app.Bot.Password)

	//UpdateTheBotUserIfNeeded()
	// TODO: Implement this one with the new function made in last attempt

	//FindBotTeam()
	// TODO: Implement with GetTeam()

	// This is an important step.  Lets make sure we use the botTeam
	// for all future web service requests that require a team.
	//client.SetTeamId(botTeam.Id)

	// Lets create a bot channel for logging debug messages into
	//CreateBotDebuggingChannelIfNeeded()
	//SendMsgToDebuggingChannel("_"+SAMPLE_NAME+" has **started** running_", "")

	//// Lets start listening to some channels via the websocket!
	//webSocketClient, err := model.NewWebSocketClient4("ws://localhost:8065", client.AuthToken)
	//if err != nil {
	//	println("We failed to connect to the web socket")
	//	PrintError(err)
	//}

	//webSocketClient.Listen()

	//go func() {
	//	for {
	//		select {
	//		case resp := <-webSocketClient.EventChannel:
	//			HandleWebSocketResponse(resp)
	//		}
	//	}
	//}()

	// You can block forever with
	//select {}
}
