// +build js,wasm

package wasm

import (
	"syscall/js"
	// "fmt"
)

const (
	goBridge = "__gobridge__"
)

var (
	bridge js.Value
)

// GoFunc is go function being exported to js
type GoFunc func(this js.Value, args []js.Value) (interface{}, error)

// JSFunc is function exported
type JSFunc func(this js.Value, args []js.Value) interface{}

func init() {
	bridge = js.Global().Get(goBridge)
}

// RegisterFuncPromise register function as Promise on bridge
func RegisterFuncPromise(name string, fn GoFunc) {
	bridge.Set(name, js.FuncOf(wrapFuncPromise(fn)))
}

// RegisterFuncCallback register function with callback on bridge
func RegisterFuncCallback(name string, fn GoFunc) {
	bridge.Set(name, js.FuncOf(wrapFuncCallback(fn)))
}

// RegisterFunc register function on bridge
func RegisterFunc(name string, fn GoFunc) {
	bridge.Set(name, js.FuncOf(wrapFunc(fn)))
}

func wrapFuncPromise(fn GoFunc) JSFunc {
	return func(this js.Value, args []js.Value) interface{} {
		// fmt.Println("before wrap", len(args), args)
		// defer fmt.Println("after wrap")
		// argLen := len(args)
		resolve := args[0]
		reject := args[1]
		ret, err := fn(this, args[2:])
		if err != nil {
			return reject.Invoke(err.Error())
		}
		return resolve.Invoke(ret)

	}
}

func wrapFuncCallback(fn GoFunc) JSFunc {
	return func(this js.Value, args []js.Value) interface{} {
		argLen := len(args)
		cb := args[argLen-1]
		ret, err := fn(this, args[:argLen-1])
		if err != nil {
			cb.Invoke(err.Error(), js.Undefined())
		} else {
			cb.Invoke(js.Undefined(), ret)
		}
		return ret
	}
}

type fnResult struct {
	err    string
	result interface{}
}

func wrapFunc(fn GoFunc) JSFunc {
	return func(this js.Value, args []js.Value) interface{} {
		ret, err := fn(this, args)
		return fnResult{
			err:    err.Error(),
			result: ret,
		}
	}
}
