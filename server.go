package bot

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mattermost/mattermost-server/model"
)

type ProtocolType int

const (
	HTTPServer ProtocolType = iota
	WSServer
)

type Server struct {
	SeverType      string `yaml:"server_type"`
	Host           string `yaml:"host"`
	HTTPS          bool   `yaml:"https"`
	WSS            bool   `yaml:"wss"`
	Account        *model.User
	HTTPClient     *model.Client4
	WSClient       *model.WebSocketClient
	TeamName       string      `yaml:"team"`
	Team           *model.Team `yaml:"-"`
	DebugChannel   *model.Channel
	Channels       []Channel `yaml:"channels"`
	CachedUsers    map[string]*model.User
	CachedChannels map[string]*model.Channel
}

func (self Server) GetDebugChannel() *model.Channel {
	var debugChannelName string
	for _, channel := range self.Channels {
		if channel.Debug {
			debugChannelName = channel.Name
		}
	}
	if debugChannelName == "" {
		debugChannelName = "bots"
	}

	debugChannel, response := self.HTTPClient.GetChannelByName(debugChannelName, self.Team.Id, "")
	if response.Error != nil {
		if self.CreateChannel(debugChannelName, "Bot Testing", "Bot debug channel used for testing bots") {
			debugChannel, response := self.HTTPClient.GetChannelByName(debugChannelName, self.Team.Id, "")
			if response.Error != nil {
				FatalError("[Fatal Error] Failed to join newly created debug channel", nil)
			} else {
				self.DebugChannel = debugChannel
			}
		} else {
			FatalError("[Fatal Error] Failed to create debug channel", nil)
			return nil
		}
	} else {
		self.DebugChannel = debugChannel
	}
	return debugChannel

}

func (self Server) JoinChannels() Server {
	for _, channel := range self.Channels {
		channel.API = self.GetChannel(channel.Name)
	}
	return self
}

func (self Server) Connect(login, password string) Server {
	self.HTTPClient = model.NewAPIv4Client(self.ServerAddress(HTTPServer))
	self.Ping()
	fmt.Println("[SERVER] Logging into server as:", login)
	self.Account = self.Login(login, password)
	fmt.Println("[SERVER] Logged into to server and obtained AuthToken:", self.HTTPClient.AuthToken)
	self.Team = self.GetTeam(self.TeamName)
	fmt.Println("[SERVER] Team ID:", self.Team.Id)
	fmt.Println("[SERVER] Connecting to wss address:", self.ServerAddress(WSServer))
	self.WSClient, _ = model.NewWebSocketClient4(self.ServerAddress(WSServer), self.HTTPClient.AuthToken)
	self.WSClient.Listen()
	self.DebugChannel = self.GetDebugChannel()
	fmt.Println("[SERVER] Using the following channel for debugging:", self.DebugChannel.Name)
	self.CachedUsers = make(map[string]*model.User)
	self.CachedChannels = make(map[string]*model.Channel)
	var updatedChannels []Channel
	for _, channel := range self.Channels {
		channel.API = self.GetChannel(channel.Name)
		updatedChannels = append(updatedChannels, channel)
	}
	self.Channels = updatedChannels
	self.HandleSignals()
	fmt.Println("[SERVER] Connecting to server with bot named: ", login)
	return (self.JoinChannels())
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

func (self Server) GetUser(userId string) *model.User {
	user, response := self.HTTPClient.GetUser(userId, "")
	if response.Error != nil {
		RuntimeError("Unable to get user", response.Error)
	}
	return user
}

func (self Server) GetTeam(teamName string) (team *model.Team) {
	team, response := self.HTTPClient.GetTeamByName(teamName, "")
	if response.Error != nil {
		RuntimeError(fmt.Sprintf(UNABLE_TO_JOIN_TEAM, teamName, ":"), response.Error)
	}
	return team
}

func (self Server) GetChannel(channelName string) (channel *model.Channel) {
	channel, response := self.HTTPClient.GetChannelByName(channelName, self.Team.Id, "")
	if response.Error != nil {
		RuntimeError(UNABLE_TO_CREATE_OR_JOIN_CHANNEL, response.Error)
	}
	return channel
}

func (self Server) CreateChannel(channelName, displayName, channelPurpose string) bool {
	channel := &model.Channel{
		Name:        channelName,
		DisplayName: displayName,
		Purpose:     channelPurpose,
		Type:        model.CHANNEL_OPEN,
		TeamId:      self.Team.Id,
	}
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

func (self Server) HandleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if self.WSClient != nil {
				self.WSClient.Close()
			}
			os.Exit(0)
		}
	}()
}
