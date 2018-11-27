package skit

import "github.com/nlopes/slack"

// OnMessage is called when a message is received.
type OnMessage func(sk *Skit, ev *slack.MessageEvent)

// OnUserTyping is called when a user starts typing direct message.
type OnUserTyping func(sk *Skit, ev *slack.UserTypingEvent)
