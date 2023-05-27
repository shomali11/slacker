package slacker

// LoggingMiddleware middleware that logs requests
func LoggingMiddleware() MiddlewareHandler {
	return func(next CommandHandler) CommandHandler {
		return func(ctx CommandContext) {
			infof(
				"%s executed \"%s\" with parameters %v in channel %s\n",
				ctx.Event().UserID,
				ctx.Usage(),
				ctx.Request().Properties(),
				ctx.Event().Channel.ID,
			)
			next(ctx)
		}
	}
}
