package skit

import (
	"context"
	"reflect"

	"github.com/nlopes/slack"
)

// New initializes an instance of skit with default event handlers.
func New(cfg Config, logger Logger) (*Skit, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	sk := &Skit{}
	sk.Logger = logger
	sk.cfg = cfg
	sk.connected = false

	return sk, nil
}

// Skit represents an instance of skit.
type Skit struct {
	Logger

	// internal states
	self      string
	cfg       Config
	connected bool
	client    *slack.Client

	handlers []Handler
}

// Register a handler to handle message events. Handlers will be executed
// in the order they are registered in.
func (sk *Skit) Register(handler Handler) {
	sk.handlers = append(sk.handlers, handler)
}

// SendText sends the given message to the channel.
func (sk *Skit) SendText(ctx context.Context, msg string, channel string) error {
	_, _, _, err := sk.client.SendMessageContext(ctx, channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true),
	)
	return err
}

// Listen connects to slack with the given configurations and starts
// the event loop
func (sk *Skit) Listen(ctx context.Context) error {
	sk.client = slack.New(sk.cfg.Token)
	rtm := sk.client.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-rtm.IncomingEvents:
			if err := sk.routeEvent(ev); err != nil {
				return err
			}
		}
	}
}

func (sk *Skit) routeEvent(rtmEv slack.RTMEvent) error {
	switch ev := rtmEv.Data.(type) {
	case *slack.HelloEvent:
		sk.Debugf("HelloEvent received")

	case *slack.ConnectingEvent:
		sk.connected = false
		sk.Infof("connecting to slack: attempt=%d", ev.Attempt)

	case *slack.ConnectedEvent:
		sk.connected = true
		sk.self = ev.Info.User.ID
		sk.Infof("connected to slack: %s", ev.Info.User.ID)

	case *slack.MessageEvent:
		if ev.Msg.User == sk.self {
			return nil
		}
		sk.Debugf("message received: channel=%s", ev.Channel)
		sk.handleMessageEvent(ev)

	case *slack.UserTypingEvent:
		sk.Debugf("ignoring user typing event")

	case *slack.RTMError:
		sk.Errorf("rtm error received: %s", ev)
		return ev

	case *slack.LatencyReport:
		sk.Infof("latency received: %s", ev.Value)

	case *slack.UserChangeEvent:
		sk.Debugf("user change event: %s", ev.User.Name)

	default:
		sk.Warnf("unknown event: %s", reflect.TypeOf(ev))
	}

	return nil
}

func (sk *Skit) handleMessageEvent(ev *slack.MessageEvent) {
	for _, handler := range sk.handlers {
		if handler.Handle(sk, ev) {
			return
		}
	}

	msg := "I don't know what to say. :neutral_face:"
	sk.SendText(context.Background(), msg, ev.Channel)
}

// Logger implementation is responsible for providing logging functions.
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
