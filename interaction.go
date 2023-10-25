package slacker

import "github.com/slack-go/slack"

// InteractionDefinition structure contains definition of the bot interaction
type InteractionDefinition struct {
	InteractionID string
	Middlewares   []InteractionMiddlewareHandler
	Handler       InteractionHandler
	Type          slack.InteractionType
}

// newInteraction creates a new bot interaction object
func newInteraction(definition *InteractionDefinition) *Interaction {
	return &Interaction{
		definition: definition,
	}
}

// Interaction structure contains the bot's interaction, description and handler
type Interaction struct {
	definition *InteractionDefinition
}

// Definition returns the interaction definition
func (c *Interaction) Definition() *InteractionDefinition {
	return c.definition
}
