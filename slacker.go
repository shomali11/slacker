package slacker

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker/expression"
	"strings"
)

const (
	space               = " "
	star                = "*"
	tick                = "`"
	dash                = "-"
	underscore          = "_"
	newLine             = "\n"
	invalidToken        = "Invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	noCommandsAvailable = "No commands were setup."
)

// NewClient creates a new client using the Slack API
func NewClient(token string) *Slacker {
	client := slack.New(token)
	rtm := client.NewRTM()

	return &Slacker{Client: client, rtm: rtm}
}

// Slacker contains the Slack API, commands, and handlers
type Slacker struct {
	Client       *slack.Client
	rtm          *slack.RTM
	commands     []*Command
	initHandler  func()
	errorHandler func(err string)
}

// Init handle the event when the bot is first connected
func (s *Slacker) Init(initHandler func()) {
	s.initHandler = initHandler
}

// Err handle when errors are encountered
func (s *Slacker) Err(errorHandler func(err string)) {
	s.errorHandler = errorHandler
}

// Command define a new command and append it to the list of existing commands
func (s *Slacker) Command(usage string, description string, handler func(request *Request, response *Response)) {
	s.commands = append(s.commands, NewCommand(usage, description, handler))
}

// Listen receives events from Slack and each is handled as needed
func (s *Slacker) Listen() error {
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			if s.initHandler == nil {
				continue
			}
			go s.initHandler()

		case *slack.MessageEvent:
			if s.isFromBot(event) {
				continue
			}

			if !s.isBotMentioned(event) && !s.isDirectMessage(event) {
				continue
			}
			go s.handleMessage(event)

		case *slack.RTMError:
			if s.errorHandler == nil {
				continue
			}
			go s.errorHandler(event.Error())

		case *slack.InvalidAuthEvent:
			return errors.New(invalidToken)
		}
	}
	return nil
}

func (s *Slacker) sendMessage(text string, channel string) {
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channel))
}

func (s *Slacker) isFromBot(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return event.User == info.User.ID
}

func (s *Slacker) isBotMentioned(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return strings.Contains(event.Text, fmt.Sprintf(userMentionFormat, info.User.ID))
}

func (s *Slacker) isDirectMessage(event *slack.MessageEvent) bool {
	return strings.HasPrefix(event.Channel, directChannelMarker)
}

func (s *Slacker) isHelpRequest(event *slack.MessageEvent) bool {
	return strings.Contains(strings.ToLower(event.Text), helpCommand)
}

func (s *Slacker) handleHelp(channel string) {
	if len(s.commands) == 0 {
		s.sendMessage(underscore+noCommandsAvailable+underscore, channel)
		return
	}

	helpMessage := empty
	for _, command := range s.commands {
		tokens := strings.Split(command.usage, space)
		for _, token := range tokens {
			if expression.IsParameter(token) {
				helpMessage += tick + token[1:len(token)-1] + tick + space
			} else {
				helpMessage += star + token + star + space
			}
		}
		helpMessage += dash + space + underscore + command.description + underscore + newLine
	}

	s.sendMessage(helpMessage, channel)
}

func (s *Slacker) handleMessage(event *slack.MessageEvent) {
	if s.isHelpRequest(event) {
		s.handleHelp(event.Channel)
		return
	}

	for _, cmd := range s.commands {
		isMatch, parameters := cmd.Match(event.Text)
		if !isMatch {
			continue
		}

		cmd.Execute(NewRequest(event, parameters), NewResponse(event.Channel, s.rtm))
		return
	}
}
