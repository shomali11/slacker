package slacker

// LoggingMiddleware middleware that logs requests
func LoggingMiddleware() MiddlewareHandler {
	return func(next CommandHandler) CommandHandler {
		return func(botCtx CommandContext) {
			infof(
				"%s executed \"%s\" with parameters %v in channel %s\n",
				botCtx.Event().UserID,
				botCtx.Usage(),
				botCtx.Request().Properties(),
				botCtx.Event().Channel.ID,
			)
			next(botCtx)
		}
	}
}
