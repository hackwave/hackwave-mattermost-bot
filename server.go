package bot

import "github.com/mattermost/mattermost-server/model"

type Server struct {
	SeverType    string `yaml:"server_type"`
	Host         string `yaml:"host"`
	HTTPS        bool   `yaml:"https"`
	WSS          bool   `yaml:"wss"`
	HTTPClient   *model.Client4
	WSClient     *model.WebSocketClient
	DebugChannel *model.Channel
	Channels     []Channel `yaml:"channels"`
}

func (self Server) ServerAddress(protocolType ProtocolType) string {
	if protocolType == HTTPServer {
		if self.HTTPS {
			return "https://" + self.Host
		} else {
			return "http://" + self.Host
		}
	} else if protocolType == WSServer {
		if self.WSS {
			return "wss://" + self.Host
		} else {
			return "ws://" + self.Host
		}
	}
	FatalError(UNABLE_TO_GENERATE_SERVER_ADDRESS, nil)
	// This never happens because FatalError calls os.Exit(1)
	return ""
}
