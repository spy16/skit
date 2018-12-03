package skit

import (
	"context"
	"reflect"

	"github.com/nlopes/slack"
)

// New initializes an instance of skit with default event handlers.
func New(cfg Config, logger Logger, opts ...Option) (*Skit, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	sl := &Skit{}
	sl.Logger = logger
	sl.cfg = cfg
	sl.connected = false

	for _, opt := range opts {
		opt(sl)
	}
	return sl, nil
}

// Skit represents an instance of skit.
type Skit struct {
	Logger

	// internal states
	self      string
	cfg       Config
	connected bool
	client    *slack.Client

	// event handlers
	onMessage    OnMessage
	onUserTyping OnUserTyping
}

// SendText sends the given message to the channel.
func (sl *Skit) SendText(ctx context.Context, msg string, channel string) error {
	_, _, _, err := sl.client.SendMessageContext(ctx, channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true),
	)
	return err
}

// Listen connects to slack with the given configurations and starts
// the event loop
func (sl *Skit) Listen(ctx context.Context) error {
	sl.client = slack.New(sl.cfg.Token)
	rtm := sl.client.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-rtm.IncomingEvents:
			if err := sl.routeEvent(ev); err != nil {
				return err
			}
		}
	}
}

func (sl *Skit) routeEvent(rtmEv slack.RTMEvent) error {
	switch ev := rtmEv.Data.(type) {
	case *slack.HelloEvent:
		sl.Debugf("HelloEvent received")

	case *slack.ConnectingEvent:
		sl.connected = false
		sl.Infof("connecting to slack: attempt=%d", ev.Attempt)

	case *slack.ConnectedEvent:
		sl.connected = true
		sl.self = ev.Info.User.ID
		sl.Infof("connected to slack: %s", ev.Info.User.ID)

	case *slack.MessageEvent:
		if ev.Msg.User == sl.self {
			return nil
		}
		sl.Debugf("message received: channel=%s", ev.Channel)
		if sl.onMessage != nil && sl.connected {
			sl.onMessage(sl, ev)
		}

	case *slack.UserTypingEvent:
		sl.Debugf("ignoring user typing event")
		if sl.onUserTyping != nil && sl.connected {
			sl.onUserTyping(sl, ev)
		}

	case *slack.RTMError:
		sl.Errorf("rtm error received: %s", ev)
		return ev

	case *slack.LatencyReport:
		sl.Infof("latency received: %s", ev.Value)

	case *slack.UserChangeEvent:
		sl.Debugf("user change event: %s", ev.User.Name)

	default:
		sl.Warnf("unknown event: %s", reflect.TypeOf(ev))
	}

	return nil
}

// Logger implementation is responsible for providing logging functions.
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
