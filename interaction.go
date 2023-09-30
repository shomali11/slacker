package slacker

// InteractionDefinition structure contains definition of the bot interaction
type InteractionDefinition struct {
	BlockID     string
	CallbackID  string
	Middlewares []InteractionMiddlewareHandler
	Handler     InteractionHandler
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
