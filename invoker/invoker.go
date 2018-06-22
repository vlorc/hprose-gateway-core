package invoker

import (
	"context"
	"github.com/vlorc/hprose-gateway-types"
	"reflect"
)

type Invoker []types.Plugin

func NewInvoker(p ...types.Plugin) Invoker {
	return Invoker(p)
}

func (i Invoker) Invoke(ctx context.Context, name string, args []reflect.Value) ([]reflect.Value, error) {
	pos := len(i) - 1
	var next types.InvokeHandler
	next = func(c context.Context, m string, a []reflect.Value) (result []reflect.Value, err error) {
		if pos--; pos >= 0 {
			result, err = i[pos].Handler(next, c, m, a)
		}
		return
	}

	return i[pos].Handler(next, ctx, name, args)
}

func (i Invoker) Len() int {
	return len(i)
}
func (i Invoker) Less(a, b int) bool {
	return i[a].Level() < i[b].Level()
}
func (i Invoker) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}
