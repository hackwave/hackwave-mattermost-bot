package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/mattermost/mattermost-server/model"
)

type Bot struct {
	Server         `yaml:"server"`
	Shell          *ishell.Shell
	Email          string `yaml:"email"`
	Password       string `yaml:"password"`
	Username       string `yaml:"display_name"`
	FirstName      string `yaml:"first_name"`
	LastName       string `yaml:"last_name"`
	ActiveChannel  *model.Channel
	RegexFunctions map[string]*RegexFunction
}

type RegexFunction struct {
	Name     string
	Regex    string
	Function func(post *model.Post)
}

func (self Bot) Start(debug bool) {
	if debug {
		self.SendDebugMessage("[BOT] "+self.Username+" in DEBUG MODE has joined the channel.", "")
	} else {
		self.SendDebugMessage("[BOT] "+self.Username+" in NORMAL MODE has joined the channel.", "")
	}
	self.UpdateServerProfile()
	self.HandleWS()
}

func (self Bot) OpenShell() {
	self.Shell = ishell.New()
	self.Shell.Println(">>> Opening Chat Interface")
	self.Shell.Println(">>>   Enter manual chat messages that will be posted by the bot")
	self.Shell.Println(">>>   send {message you want to send as the bot}")

	self.ActiveChannel = self.Server.Channels[0].API
	fmt.Println("[CONFIG] Active channel is now:", self.ActiveChannel.Name)

	self.Shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list available channels",
		Func: func(c *ishell.Context) {
			fmt.Println("\nChannels:")
			for _, channel := range self.Server.Channels {
				fmt.Println("  " + channel.Name)
			}
			fmt.Println("\n")
		},
	})
	self.Shell.AddCmd(&ishell.Cmd{
		Name: "channel",
		Help: "select active channel to send messages to",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				selected := false
				for _, channel := range self.Server.Channels {
					if channel.Name == c.Args[0] {
						fmt.Println("[CONFIG] Active channel is now:", channel.Name)
						self.ActiveChannel = channel.API
						selected = true
					}
				}
				if !selected {
					fmt.Println("[CONFIG] Selected channel is not one of the available channels.")
				}
			} else {
				fmt.Println("[Error] No channel supplied to set as active channel.")
			}
		},
	})
	self.Shell.AddCmd(&ishell.Cmd{
		Name: "send",
		Help: "send message",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				self.SendMessage(strings.Join(c.Args, " "), "")
			} else {
				fmt.Println("[Error] No message provided, nothing sent.")
			}
		},
	})
	self.Shell.Run()
}

func (self Bot) AddShellCommand(cmd *ishell.Cmd) {
	self.Shell.AddCmd(cmd)
}

// TODO: The idea for this is that you would use it in shell mode, and allow the user to select
// active channels the shell is watching
func (self Bot) SendMessage(message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	post := &model.Post{
		ChannelId: self.ActiveChannel.Id,
		Message:   message,
		RootId:    replyToId,
	}
	return (self.Server.SendPost(post))
}

func (self Bot) SendMessageToChannelWithId(channelId, message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
		RootId:    replyToId,
	}
	return (self.Server.SendPost(post))
}

func (self Bot) SendMessageToChannelWithName(channelName, message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	// Use Channel Caching
	channel := self.Server.CachedChannels[channelName]
	if channel == nil {
		channel = self.Server.GetChannel(channelName)
		self.Server.CachedChannels[channelName] = channel
	}
	post := &model.Post{
		ChannelId: channel.Id,
		Message:   message,
		RootId:    replyToId,
	}
	return (self.Server.SendPost(post))
}

func (self Bot) SendDebugMessage(message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	if self.Server.DebugChannel != nil {
		post := &model.Post{
			ChannelId: self.Server.DebugChannel.Id,
			Message:   message,
			RootId:    replyToId,
		}
		return (self.Server.SendPost(post))
	} else {
		return false
	}
}

func (self Bot) RegisterHook(regexFunction RegexFunction) Bot {
	if self.RegexFunctions == nil {
		self.RegexFunctions = make(map[string]*RegexFunction)
	}

	self.RegexFunctions[regexFunction.Name] = &regexFunction
	return self
}

func (self Bot) HandleWS() {
	go func() {
		for {
			select {
			case event := <-self.Server.WSClient.EventChannel:
				self.HandleMessageFromChannel(event)
			}
		}
	}()
}

func (self Bot) HandleMessageFromChannel(event *model.WebSocketEvent) {
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}
	//fmt.Println("\t", event.Data["post"].(string))
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		if post.UserId == self.Account.Id {
			return
		}

		// Cache users
		user := self.Server.CachedUsers[post.UserId]
		if user == nil {
			user = self.Server.GetUser(post.UserId)
			self.Server.CachedUsers[post.UserId] = user
		}

		// Use Channel Caching
		channel := self.Server.CachedChannels[post.ChannelId]
		if channel == nil {
			for _, c := range self.Server.Channels {
				if c.API.Id == post.ChannelId {
					self.Server.CachedChannels[post.ChannelId] = c.API
					channel = c.API

				}
			}
		}
		if user != nil && channel != nil {
			fmt.Println("\t(CHAT)["+channel.Name+"]["+user.Username+"]", post.Message)
		}
		for _, regexFunction := range self.RegexFunctions {
			if matched, _ := regexp.MatchString(regexFunction.Regex, post.Message); matched {
				fmt.Println("[BOT] MATCHED regex function:", regexFunction.Name)
				regexFunction.Function(post)
				return
			}
		}
	}
}
