package skit

import "context"

// Handler is responsible for handling slack message events.
type Handler interface {

	// Handle should return true if the message event was handled by it. If this
	// function returns false, skit will execute the next Handler available and
	// so on.
	Handle(ctx context.Context, sk *Skit, ev *MessageEvent) bool
}

// HandlerFunc implements Handler interface using function type.
type HandlerFunc func(ctx context.Context, sk *Skit, ev *MessageEvent) bool

// Handle dispatches the call to the wrapped function.
func (hf HandlerFunc) Handle(ctx context.Context, sk *Skit, ev *MessageEvent) bool {
	return hf(ctx, sk, ev)
}
