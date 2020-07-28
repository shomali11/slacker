package slacker

import (
	"context"
	"errors"

	"github.com/slack-go/slack"
)

// EventHandler handles an incoming RTM event.
type EventHandler func(ctx context.Context, s *Slacker, msg slack.RTMEvent) error

// DefaultEventHandler it the default event handler.
func DefaultEventHandler(ctx context.Context, s *Slacker, msg slack.RTMEvent) error {
	switch event := msg.Data.(type) {
	case *slack.ConnectedEvent:
		if s.initHandler == nil {
			return nil
		}
		go s.initHandler()

	case *slack.MessageEvent:
		if s.isFromBot(event) {
			return nil
		}

		if !s.isBotMentioned(event) && !s.isDirectMessage(event) {
			return nil
		}
		go s.handleMessage(ctx, event)

	case *slack.RTMError:
		if s.errorHandler == nil {
			return nil
		}
		go s.errorHandler(event.Error())

	case *slack.InvalidAuthEvent:
		return errors.New(invalidToken)

	default:
		if s.fallbackEventHandler == nil {
			return nil
		}
		go s.fallbackEventHandler(event)
	}

	return nil
}
