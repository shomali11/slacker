package slacker

import "github.com/slack-go/slack"

func executeCommand(ctx *CommandContext, handler CommandHandler, middlewares ...CommandMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}

func executeInteraction(ctx *InteractionContext, handler InteractionHandler, middlewares ...InteractionMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}

func executeSuggestion(ctx *InteractionContext, handler SuggestionHandler, middlewares ...SuggestionMiddlewareHandler) slack.OptionsResponse {
	if handler == nil {
		return slack.OptionsResponse{}
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler(ctx)
}

func executeJob(ctx *JobContext, handler JobHandler, middlewares ...JobMiddlewareHandler) func() {
	if handler == nil {
		return func() {}
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return func() {
		handler(ctx)
	}
}
