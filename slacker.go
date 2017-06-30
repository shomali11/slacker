package slacker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
)

const (
	space               = " "
	dash                = "-"
	newLine             = "\n"
	invalidToken        = "Invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	codeMessageFormat   = "`%s`"
	boldMessageFormat   = "*%s*"
	italicMessageFormat = "_%s_"
	noCommandsAvailable = "No botCommands were setup."
)

// NewClient creates a new client using the Slack API
func NewClient(token string) *Slacker {
	client := slack.New(token)
	rtm := client.NewRTM()

	return &Slacker{Client: client, rtm: rtm}
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	Client         *slack.Client
	rtm            *slack.RTM
	botCommands    []*BotCommand
	initHandler    func()
	errorHandler   func(err string)
	defaultHandler func(request *Request, response ResponseWriter)
}

// Init handle the event when the bot is first connected
func (s *Slacker) Init(initHandler func()) {
	s.initHandler = initHandler
}

// Err handle when errors are encountered
func (s *Slacker) Err(errorHandler func(err string)) {
	s.errorHandler = errorHandler
}

// Default handle when none of the commands are matched
func (s *Slacker) Default(defaultHandler func(request *Request, response ResponseWriter)) {
	s.defaultHandler = defaultHandler
}

// Command define a new command and append it to the list of existing commands
func (s *Slacker) Command(usage string, description string, handler func(request *Request, response ResponseWriter)) {
	s.botCommands = append(s.botCommands, NewBotCommand(usage, description, handler))
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
	if len(s.botCommands) == 0 {
		s.sendMessage(fmt.Sprintf(italicMessageFormat, noCommandsAvailable), channel)
		return
	}

	helpMessage := empty
	for _, command := range s.botCommands {
		tokens := strings.Split(command.usage, space)
		for _, token := range tokens {
			if commander.IsParameter(token) {
				helpMessage += fmt.Sprintf(codeMessageFormat, token[1:len(token)-1]) + space
			} else {
				helpMessage += fmt.Sprintf(boldMessageFormat, token) + space
			}
		}
		helpMessage += dash + space + fmt.Sprintf(italicMessageFormat, command.description) + newLine
	}

	s.sendMessage(helpMessage, channel)
}

func (s *Slacker) handleMessage(event *slack.MessageEvent) {
	if s.isHelpRequest(event) {
		s.handleHelp(event.Channel)
		return
	}

	response := NewResponse(event.Channel, s.rtm)
	for _, cmd := range s.botCommands {
		parameters, isMatch := cmd.Match(event.Text)
		if !isMatch {
			continue
		}

		cmd.Execute(NewRequest(event, parameters), response)
		return
	}

	s.defaultHandler(NewRequest(event, &proper.Properties{}), response)
}
