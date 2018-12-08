package lua

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// NewWrapper initializes a new wrapper with empty Lua state
func NewWrapper(paths ...string) *Wrapper {
	wr := &Wrapper{}

	if len(paths) > 0 {
		pathVal := strings.Join(paths, ";")
		os.Setenv("LUA_PATH", pathVal)
	}

	wr.state = lua.NewState()
	return wr
}

// Wrapper is a thin wrapper around gopher-lua defined
// lua state
type Wrapper struct {
	state *lua.LState
}

// State returns the actual lua state object
func (wr *Wrapper) State() *lua.LState {
	return wr.state
}

// Bind creates a global variable with given value and name
// thus exposing the value to the lua script
func (wr *Wrapper) Bind(name string, v interface{}) {
	wr.state.SetGlobal(name, luar.New(wr.state, v))
}

// BindAll exposes all the values in the map to the lua scripts
// by iteratively calling Bind with key and value
func (wr *Wrapper) BindAll(vals map[string]interface{}) {
	for key, val := range vals {
		wr.Bind(key, val)
	}
}

// Execute the given lua script string
func (wr *Wrapper) Execute(src string) error {
	return wr.state.DoString(src)
}

// ExecuteFile reads and executes the lua file
func (wr *Wrapper) ExecuteFile(fileName string) error {
	return wr.state.DoFile(fileName)
}

// GetFunction returns a function object defined in Lua if found.
func (wr *Wrapper) GetFunction(name string) (*lua.LFunction, error) {
	fn := wr.state.GetGlobal(name)

	lfn, ok := fn.(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", name)
	}

	return lfn, nil
}

// Call a lua function by its name. Args are automatically converted to
// appropriate types using the Luar library
func (wr *Wrapper) Call(name string, args ...interface{}) (lua.LValue, error) {
	lfn, err := wr.GetFunction(name)
	if err != nil {
		return nil, err
	}
	return wr.CallFunc(lfn, args...)
}

// CallFunc calls the given lua function object with given args.
func (wr *Wrapper) CallFunc(lfn *lua.LFunction, args ...interface{}) (lua.LValue, error) {
	wr.state.Push(lfn)
	for _, arg := range args {
		wr.state.Push(luar.New(wr.state, arg))
	}
	err := wr.state.PCall(len(args), 1, nil)
	if err != nil {
		return nil, err
	}

	top := wr.state.GetTop()
	retVal := wr.state.Get(top)
	return retVal, nil
}

// Reset the lua state. Same as creating new instance
// of Wrapper
func (wr *Wrapper) Reset() {
	wr.state = lua.NewState()
}
