package slacker

// LoggingMiddleware middleware that logs requests
func LoggingMiddleware() CommandMiddlewareHandler {
	return func(next CommandHandler) CommandHandler {
		return func(ctx CommandContext) {
			infof(
				"%s executed \"%s\" with parameters %v in channel %s\n",
				ctx.Event().UserID,
				ctx.Definition().Usage,
				ctx.Request().Properties(),
				ctx.Event().Channel.ID,
			)
			next(ctx)
		}
	}
}
