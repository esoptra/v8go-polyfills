package polyfills

import (
	"context"
	"fmt"

	"github.com/esoptra/v8go"
)

type Runner struct {
	resCh  chan *v8go.Value
	errCh  chan *v8go.Value
	global *v8go.ObjectTemplate
}

//NewRunner creates a runner obj with result and error of 'chan *v8go.Value' type
func NewRunner(iso *v8go.Isolate, global *v8go.ObjectTemplate) (*Runner, error) {
	resCh := make(chan *v8go.Value, 1)
	errCh := make(chan *v8go.Value, 1)

	errFun := func() v8go.FunctionCallback {
		return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			fmt.Println("errToGo")
			args := info.Args()
			errCh <- args[0]
			return nil
		}
	}

	resFun := func() v8go.FunctionCallback {
		return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			fmt.Println("resToGo")
			args := info.Args()
			resCh <- args[0]
			return nil
		}
	}
	errFn, err := v8go.NewFunctionTemplate(iso, errFun())
	if err != nil {
		return nil, fmt.Errorf("v8go-runner/errFn: %w", err)
	}
	resFn, err := v8go.NewFunctionTemplate(iso, resFun())
	if err != nil {
		return nil, fmt.Errorf("v8go-runner/resFn: %w", err)
	}
	if err := global.Set("errToGo", errFn, v8go.ReadOnly); err != nil {
		return nil, fmt.Errorf("v8go-runner/errToGo: %w", err)
	}
	if err := global.Set("resToGo", resFn, v8go.ReadOnly); err != nil {
		return nil, fmt.Errorf("v8go-runner/resToGo: %w", err)
	}

	return &Runner{
		resCh:  resCh,
		errCh:  errCh,
		global: global,
	}, nil
}

// RunPromise runs a function that resolves a promise and waits for the
// promise to resolve, reject or for the context to timeout.
// Make sure the script includes 'let res = epsilon(data);'
func (r *Runner) RunPromise(ctx context.Context, v8ctx *v8go.Context, script string) (*v8go.Value, error) {
	code := script + `
	Promise.resolve(res)`

	jsCode := code + ".then(resToGo).catch(errToGo)"
	_, err := v8ctx.RunScript(jsCode, "script.js")
	if err != nil {
		fmt.Printf("RunScript error: %v\nCode: %s", err, jsCode)
		return nil, fmt.Errorf("Eval error: %v\nCode: %s", err, jsCode)
	}
	//fmt.Println("end RunScript")

	select {
	case res := <-r.resCh:
		return res, nil
	case errVal := <-r.errCh:
		return nil, fmt.Errorf("%v", errVal)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout in RunPromise: %v", ctx.Err())
	}
}
