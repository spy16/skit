package skit

import "github.com/nlopes/slack"

// Handler is responsible for handling slack message events.
type Handler interface {

	// Handle should return true if the message event was handled by it. If this
	// function returns false, skit will execute the next Handler available and
	// so on.
	Handle(sk *Skit, ev *slack.MessageEvent) bool
}

// HandlerFunc implements Handler interface using function type.
type HandlerFunc func(sk *Skit, ev *slack.MessageEvent) bool

// Handle dispatches the call to the wrapped function.
func (hf HandlerFunc) Handle(sk *Skit, ev *slack.MessageEvent) bool {
	return hf(sk, ev)
}
