package slacker

// BotMode instruct the bot on how to handle incoming events that
// originated from a bot.
type BotMode int

const (
	// BotModeIgnoreAll instructs our bot to ignore any activity coming
	// from other bots, including our self.
	BotModeIgnoreAll BotMode = iota

	// BotModeIgnoreApp will ignore any events that originate from a
	// bot that is associated with the same App (ie. share the same App ID) as
	// this bot. OAuth scope `user:read` is required for this mode.
	BotModeIgnoreApp

	// BotModeIgnoreNone will not ignore any bots, including our self.
	// This can lead to bots "talking" to each other so care must be taken when
	// selecting this option.
	BotModeIgnoreNone
)
