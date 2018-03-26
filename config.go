package bot

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func (self Application) LoadConfiguration(configPath string) Application {
	fmt.Println("[CONFIG] Loading configuration from", configPath, "(looking in current working directory).")
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalError(UNABLE_TO_LOCATE_CONFIG_FILE, err)
	}
	err = yaml.Unmarshal(configFile, &self)
	if err != nil {
		FatalError(UNABLE_TO_PARSE_CONFIG_FILE, err)
	}
	if self.Debug {
		fmt.Println("[CONFIG] Bot is in DEBUG mode.")
	} else {
		fmt.Println("[CONFIG] Bot is in NORMAL mode.")
	}
	return self
}
