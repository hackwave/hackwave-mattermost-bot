package bot

import (
	"os"
	"os/signal"
)

type Application struct {
	Name    string
	Version string
	Bot     `yaml:"bot"`
}

func (self Application) HandleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if self.Bot.Server.WSClient != nil {
				self.Bot.Server.WSClient.Close()
			}

			self.Bot.SendDebugMessage("_"+self.Bot.Username+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}
