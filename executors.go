package slacker

func executeCommand(ctx CommandContext, handler CommandHandler, middlewares ...CommandMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}

func executeInteraction(ctx InteractionContext, handler InteractionHandler, middlewares ...InteractionMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}

func executeJob(ctx JobContext, handler JobHandler, middlewares ...JobMiddlewareHandler) func() {
	return func() {
		if handler == nil {
			return
		}

		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}

		handler(ctx)
	}
}
