/*
 * Copyright (c) 2021 Twintag
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package textEncoder

import (
	"errors"
	"fmt"

	"rogchap.com/v8go"
)

/**
Inject basic textEncoder encode and decode support.
*/
func InjectTo(ctx *v8go.Context, opt ...Option) error {
	if ctx == nil {
		return errors.New("v8go-polyfills/textEncoder: ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("v8go-polyfills/textEncoder Isolate: %w", err)
	}

	c := NewEncode(opt...)
	encodeFn, err := v8go.NewFunctionTemplate(iso, c.TextEncoderFunctionCallback())
	if err != nil {
		return fmt.Errorf("v8go-polyfills/textEncoder NewFunctionTemplate: %w", err)
	}

	global, _ := v8go.NewObjectTemplate(iso)
	if err := global.Set("TextEncoder", encodeFn); err != nil {
		return fmt.Errorf("v8go-polyfills/textEncoder global.Set: %w", err)
	}

	return nil
}

func InjectWithGlobalTo(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...Option) error {

	e := NewEncode(opt...)
	fetchFn, err := v8go.NewFunctionTemplate(iso, e.TextEncoderFunctionCallback())
	if err != nil {
		return fmt.Errorf("v8go-polyfills/textEncoder NewFunctionTemplate: %w", err)
	}

	// encodeFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	// 	ctx := info.Context()
	// 	args := info.Args()
	// 	if len(args) <= 0 {
	// 		// TODO: this should return an error, but v8go not supported now
	// 		val, _ := v8go.NewValue(iso, "")
	// 		return val
	// 	}
	// 	resolver, _ := v8go.NewPromiseResolver(ctx)
	// 	// go func() {
	// 	input := args[0].String()
	// 	v, _ := v8go.NewValue(iso, []byte(input))
	// 	resolver.Resolve(v)
	// 	// }()

	// 	return resolver.GetPromise().Value
	// })
	// if err != nil {
	// 	return fmt.Errorf("v8go-polyfills/textEncoder encodeFnTmp def: %w", err)
	// }
	// con, err := v8go.NewObjectTemplate(iso)
	// if err != nil {
	// 	return fmt.Errorf("v8go-polyfills/textEncoder new object temaplate: %w", err)
	// }

	// if err := con.Set("encode", encodeFnTmp, v8go.ReadOnly); err != nil {
	// 	return fmt.Errorf("v8go-polyfills/textEncoder encodeFnTmp registration: %w", err)
	// }

	// ctx, _ := v8go.NewContext(iso)
	// conObj, err := con.NewInstance(ctx)
	// if err != nil {
	// 	return fmt.Errorf("v8go-polyfills/textEncoder newinstance: %w", err)
	// }

	if err := global.Set("TextEncoder", fetchFn); err != nil {
		return fmt.Errorf("v8go-polyfills/textEncoder global.set: %w", err)
	}

	return nil
}
