package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

type Bot struct {
	Server    `yaml:"server"`
	Email     string `yaml:"email"`
	Password  string `yaml:"password"`
	Username  string `yaml:"display_name"`
	FirstName string `yaml:"first_name"`
	LastName  string `yaml:"last_name"`
}

func (self Bot) SendMessage(channelName, message, replyToId string) bool {
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

func (self Bot) HandleWebSocketResponse(event *model.WebSocketEvent) {
	self.HandleMessageFromDebugChannel(event)
}

func (self Bot) HandleMessageFromDebugChannel(event *model.WebSocketEvent) {
	//fmt.Println("\t", event.Data["post"].(string))
	if event.Broadcast.ChannelId != self.Server.DebugChannel.Id {
		return
	}
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

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

		// TODO: Configure this using a plugin type system that can be defined from the cli

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			self.SendDebugMessage("Yes I'm running", post.Id)
			return
		}

		// if you see any word matching 'up' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			self.SendDebugMessage("Yes I'm running", post.Id)
			return
		}

		// if you see any word matching 'running' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			self.SendDebugMessage("Yes I'm running", post.Id)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello(?:$|\W)`, post.Message); matched {
			self.SendDebugMessage("Yes I'm running", post.Id)
			return
		}
	}
}
