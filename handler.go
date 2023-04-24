package slacker

// MiddlewareHandler represents the middleware handler function
type MiddlewareHandler func(CommandHandler) CommandHandler

// CommandHandler represents the command handler function
type CommandHandler func(CommandContext)

// InteractiveHandler represents the interactive handler function
type InteractiveHandler func(InteractiveContext)

// JobHandler represents the job handler function
type JobHandler func(jobCtx JobContext) error
