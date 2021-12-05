package slacker

import (
	"time"

	allot "github.com/sdslabs/allot/pkg"
)

// NewCommandEvent creates a new command event
func NewCommandEvent(command string, parameters []allot.Parameter, event *MessageEvent) *CommandEvent {
	return &CommandEvent{
		Timestamp:  time.Now(),
		Command:    command,
		Parameters: parameters,
		Event:      event,
	}
}

// CommandEvent is an event to capture executed commands
type CommandEvent struct {
	Timestamp  time.Time
	Command    string
	Parameters []allot.Parameter
	Event      *MessageEvent
}
