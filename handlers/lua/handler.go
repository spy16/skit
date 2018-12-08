package lua

import (
	"context"
	"fmt"

	"github.com/spy16/skit"
	lua "github.com/yuin/gopher-lua"
)

// New initializes the Lua based skit handler.
func New(src, handlerName string, luaPaths []string) (*Lua, error) {
	wr := NewWrapper(luaPaths...)
	if err := wr.Execute(src); err != nil {
		return nil, err
	}

	lfn, err := wr.GetFunction(handlerName)
	if err != nil {
		return nil, err
	}

	lh := &Lua{}
	lh.wr = wr
	lh.handlerFunc = lfn
	return lh, nil
}

// Lua implements skit.Handler using a lua scripting layer.
type Lua struct {
	handler     string
	handlerFunc *lua.LFunction
	wr          *Wrapper
}

// Handle dispatches the event object to lua handler function and returns the result.
func (lh *Lua) Handle(ctx context.Context, sk *skit.Skit, ev *skit.MessageEvent) bool {
	val, err := lh.wr.CallFunc(lh.handlerFunc, ctx, sk, ev)
	if err != nil {
		sk.Errorf("failed to call lua function: %v", err)
		sk.SendText(ctx, fmt.Sprintf("Something went terribly wrong: %s", err), ev.Channel)
		return true
	}

	if handled, ok := val.(lua.LBool); ok {
		return handled.String() == "true"
	}

	sk.Warnf("lua handler function returned non-bool value: %v", val)
	sk.SendText(ctx, ":scream: My handlers are misbehaving. :pensive:", ev.Channel)
	return true
}
