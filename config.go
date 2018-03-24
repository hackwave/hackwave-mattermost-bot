package bot

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Application struct {
	Name    string
	Version string
	Bot     `yaml:"bot"`
}

func (self Application) LoadConfiguration(configPath string) Application {
	fmt.Println("[CONFIG] Loading configuration from", configPath, "(looking in current working directory).")

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalError(UNABLE_TO_LOCATE_CONFIG_FILE, err)
	}
	fmt.Println("CONTENTS OF configFile:", configFile)
	err = yaml.Unmarshal(configFile, self.Bot)
	if err != nil {
		fmt.Println("yaml tried to marshal and failed: ", err)
		FatalError(UNABLE_TO_PARSE_CONFIG_FILE, err)
	}
	fmt.Println("yaml umarshaled")
	fmt.Println("bot:")
	fmt.Println("  bot.Name", self.Bot.Name)

	return self
}
