package skit

// Option implementation can apply configurations to skit instance.
type Option func(sk *Skit) error

// WithMessageHandler registers h as message handler.
func WithMessageHandler(h OnMessage) Option {
	return func(sk *Skit) error {
		sk.onMessage = h
		return nil
	}
}

// WithUserTypingHandler registers h as a handler for User typing event.
func WithUserTypingHandler(h OnUserTyping) Option {
	return func(sk *Skit) error {
		sk.onUserTyping = h
		return nil
	}
}
