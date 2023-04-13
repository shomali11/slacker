package slacker

// BotInteractionMode instruct the bot on how to handle incoming events that
// originated from a bot.
type BotInteractionMode int

const (
	// BotInteractionModeIgnoreAll instructs our bot to ignore any activity coming
	// from other bots, including our self.
	BotInteractionModeIgnoreAll BotInteractionMode = iota

	// BotInteractionModeIgnoreApp will ignore any events that originate from a
	// bot that is associated with the same App (ie. share the same App ID) as
	// this bot. OAuth scope `user:read` is required for this mode.
	BotInteractionModeIgnoreApp

	// BotInteractionModeIgnoreNone will not ignore any bots, including our self.
	// This can lead to bots "talking" to each other so care must be taken when
	// selecting this option.
	BotInteractionModeIgnoreNone
)
