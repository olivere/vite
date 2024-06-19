package vite

import "context"

type contextKey string

var scriptsKey = contextKey("scripts")

// ScriptsFromContext returns the scripts to be injected in the HTML.
func ScriptsFromContext(ctx context.Context) string {
	if md, ok := ctx.Value(scriptsKey).(string); ok {
		return md
	}
	return ""
}

// ScriptsToContext sets the scripts to be injected in the HTML.
func ScriptsToContext(ctx context.Context, scripts string) context.Context {
	return context.WithValue(ctx, scriptsKey, scripts)
}
