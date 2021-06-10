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
	"fmt"

	"github.com/esoptra/v8go"
)

type Encoder struct {
}

func NewEncode(opt ...Option) *Encoder {
	c := &Encoder{}

	for _, o := range opt {
		o.apply(c)
	}

	return c
}

func (c *Encoder) TextEncoderFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		iso, _ := ctx.Isolate()

		encodeFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) <= 0 {
				return iso.ThrowException(fmt.Sprintf("Expected an arguments\n"))
			}
			s := args[0].String()
			v, err := v8go.NewValue(iso, []byte(s))
			if err != nil {
				return iso.ThrowException(fmt.Sprintf("error creating new val: %#v", err))
			}
			return v
		})
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error creating encode() template: %#v", err))
		}

		/*
			encodeIntoFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
				args := info.Args()
				if len(args) <= 0 {
					return iso.ThrowException(fmt.Sprintf("Expected an arguments\n"))
				}
				s := args[0].String()
				result := args[1].Uint8Array()
				i := 0
				for ; i < len(s); i++ {
					fmt.Printf("%d ", s[i])
					result[i] = s[i]
				}
				v, err := v8go.JSONParse(info.Context(), fmt.Sprintf(`{"read": %d, "written": %d }`, i, i))
				if err != nil {
					return iso.ThrowException(fmt.Sprintf("Error parsing json val\n"))
				}

				return v
			})

			if err != nil {
				return iso.ThrowException(fmt.Sprintf("error creating encodeInto() template: %#v", err))
			}
		*/
		resTmp, err := v8go.NewObjectTemplate(iso)
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error creating object template: %#v", err))
		}

		if err := resTmp.Set("encode", encodeFnTmp, v8go.ReadOnly); err != nil {
			return iso.ThrowException(fmt.Sprintf("error setting encode function template: %#v", err))
		}
		/*
			if err := resTmp.Set("encodeInto", encodeIntoFnTmp, v8go.ReadOnly); err != nil {
				return iso.ThrowException(fmt.Sprintf("error setting encodeInto function template: %#v", err))
			}
		*/
		resObj, err := resTmp.NewInstance(ctx)
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error new instance from ctx: %#v", err))
		}
		return resObj.Value
	}
}
