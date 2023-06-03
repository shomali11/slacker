package slacker

// InteractionDefinition structure contains definition of the bot interaction
type InteractionDefinition struct {
	BlockID     string
	Description string
	Middlewares []InteractionMiddlewareHandler
	Handler     InteractionHandler

	// HideHelp will hide this interaction definition from appearing in the `help` results.
	HideHelp bool
}

// newInteraction creates a new bot interaction object
func newInteraction(blockID string, definition *InteractionDefinition) Interaction {
	definition.BlockID = blockID

	return &interaction{
		definition: definition,
	}
}

// Interaction interface
type Interaction interface {
	Definition() *InteractionDefinition
}

// interaction structure contains the bot's interaction, description and handler
type interaction struct {
	definition *InteractionDefinition
}

// Definition returns the interaction definition
func (c *interaction) Definition() *InteractionDefinition {
	return c.definition
}
