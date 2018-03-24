package bot

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

type ProtocolType int

const (
	HTTPServer ProtocolType = iota
	WSServer
)

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

func (self Server) Connect(login, password string) Server {
	self.HTTPClient = model.NewAPIv4Client(self.ServerAddress(HTTPServer))
	self.Ping()
	self.Login(login, password)

	return self
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

func (self Server) GetTeam(teamName string) (team *model.Team) {
	team, response := self.HTTPClient.GetTeamByName(teamName, "")
	if response.Error != nil {
		RuntimeError(fmt.Sprintf(UNABLE_TO_JOIN_TEAM, teamName, ":"), response.Error)
	}
	return team
}

func (self Server) GetChannel(teamId, channelName string) (channel *model.Channel) {
	channel, response := self.HTTPClient.GetChannelByName(channelName, teamId, "")
	if response.Error != nil {
		RuntimeError(UNABLE_TO_CREATE_OR_JOIN_CHANNEL, response.Error)
	}
	return channel
}

func (self Server) CreateChannel(channel *model.Channel) bool {
	if _, response := self.HTTPClient.CreateChannel(channel); response.Error != nil {
		RuntimeError(UNABLE_TO_CREATE_CHANNEL, response.Error)
		return false
	} else {
		return true
	}
}

func (self Server) SendPost(post *model.Post) bool {
	if _, response := self.HTTPClient.CreatePost(post); response.Error != nil {
		RuntimeError(UNABLE_TO_SEND_MESSAGE_TO_CHANNEL, response.Error)
		return false
	} else {
		return true
	}
}

func (self Server) Ping() bool {
	if pingResponse, response := self.HTTPClient.GetOldClientConfig(""); response.Error != nil {
		FatalError(UNABLE_TO_PING_SERVER, response.Error)
		return false
	} else {
		fmt.Println("[SERVER] Server responded top ping, is running Mattermost version:" + pingResponse["Version"])
		return true
	}
}

func (self Server) Login(email, password string) (account *model.User) {
	account, response := self.HTTPClient.Login(email, password)
	if response.Error != nil {
		FatalError(UNABLE_TO_LOGIN, response.Error)
	}
	return account
}

func (self Server) UpdateAccount(account *model.User) (savedAccount *model.User) {
	savedAccount, response := self.HTTPClient.UpdateUser(account)
	if response.Error != nil {
		FatalError(UNABLE_TO_UPDATE_PROFILE, response.Error)
	}
	return savedAccount
}
