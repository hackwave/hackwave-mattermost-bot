package bot

import "github.com/mattermost/mattermost-server/model"

type Channel struct {
	API         *model.Channel
	Name        string `yaml:"name"`
	DisplayName string `yaml:"display_name"`
	Description string `yaml:"description"`
	Debug       bool   `yaml:"debug"`
}
