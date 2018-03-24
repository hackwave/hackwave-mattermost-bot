package bot

import (
	"fmt"
	"os"
)

const (
	UNABLE_TO_JOIN_TEAM               = "Failed to join team defined in configuration:"
	UNABLE_TO_CREATE_CHANNEL          = "Unable to create channel:"
	UNABLE_TO_CREATE_OR_JOIN_CHANNEL  = "Failed to create (or) join channel specified in the configuration file:"
	UNABLE_TO_GENERATE_SERVER_ADDRESS = "Unable to generate a functional server address, check your bot.yml configuration."
	UNABLE_TO_LOCATE_CONFIG_FILE      = "Failed to locate config file current working directory:"
	UNABLE_TO_PARSE_CONFIG_FILE       = "Failed to parse the YAML config file:"
	UNABLE_TO_SEND_MESSAGE_TO_CHANNEL = "Failed to send a message to channel."
	UNABLE_TO_PING_SERVER             = "Failed to ping the Matttermost server:"
	UNABLE_TO_LOGIN                   = "Failed to login with the email/password combination provided in the bot.yaml file:"
	UNABLE_TO_UPDATE_PROFILE          = "Unable to update account profile:"
	UNABLE_TO_INIT_WS_CONNECTION      = "Unable to connect to the following WS Server:"
)

func RuntimeError(message string, err error) {
	fmt.Println("[Error]", message, err)
}

func FatalError(message string, err error) {
	fmt.Println("[Fatal Error]", message, err)
	os.Exit(1)
}
