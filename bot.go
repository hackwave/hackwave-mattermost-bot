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
	Function func()
}

func (self Bot) Start() {
	self.SendDebugMessage("[BOT] "+self.Username+" in DEBUG MODE has joined the channel.", "")
	self.UpdateServerProfile()
	self.HandleWS()
}

func (self Bot) OpenShell() {
	self.Shell = ishell.New()
	self.Shell.Println(">>> Opening Chat Interface")
	self.Shell.Println(">>>   Enter manual chat messages that will be posted by the bot")
	self.Shell.Println(">>>   send {message you want to send as the bot}")
	self.Shell.AddCmd(&ishell.Cmd{
		Name: "send",
		Help: "send message",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				self.SendDebugMessage(strings.Join(c.Args, " "), "")
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

func (self Bot) SendMessage(message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	post := &model.Post{
		ChannelId: self.ActiveChannel.Id,
		Message:   message,
		RootId:    replyToId,
	}
	return (self.Server.SendPost(post))
}

func (self Bot) SendMessageToChannel(channelName, message, replyToId string) bool {
	fmt.Println("\t(CHAT)[", self.Username, "]", message)
	channel := self.Server.GetChannel(channelName)
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
				self.HandleMessageFromDebugChannel(event)
			}
		}
	}()
}

func (self Bot) HandleMessageFromDebugChannel(event *model.WebSocketEvent) {
	//fmt.Println("\t", event.Data["post"].(string))
	if event.Broadcast.ChannelId != self.Server.DebugChannel.Id {
		return
	}
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	// TODO: Bot is not filtering out own messages

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		if post.UserId == self.Account.Id {
			return
		}

		// Cache users
		var user *model.User
		if self.Server.Users[post.UserId] == nil {
			user = self.Server.GetUser(post.UserId)
			self.Server.Users[post.UserId] = user
		} else {
			user = self.Server.Users[post.UserId]
		}

		fmt.Println("\t(CHAT)[", user.Username, "]", post.Message)

		fmt.Println("[BOT] Checking", len(self.RegexFunctions), "regex functions.")
		for _, regexFunction := range self.RegexFunctions {
			if matched, _ := regexp.MatchString(regexFunction.Regex, post.Message); matched {
				fmt.Println("[BOT] MATCHED regex function:", regexFunction.Name)
				regexFunction.Function()
				return
			}
		}
	}
}
