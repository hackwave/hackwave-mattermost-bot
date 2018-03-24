package main

import (
	"fmt"
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
}
