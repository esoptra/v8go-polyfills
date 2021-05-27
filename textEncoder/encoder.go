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
	"log"

	"rogchap.com/v8go"
)

// type TextEncoder interface {
// 	EncodeFunctionCallback() v8go.FunctionCallback
// 	EncodeIntoFunctionCallback() v8go.FunctionCallback
// }

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
	fmt.Println("at the root level TextEncoderFunctionCallback")
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		fmt.Println("at the call func lwvl")
		// resolver, _ := v8go.NewPromiseResolver(ctx)

		iso, _ := ctx.Isolate()

		encodeFnTmp, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			fmt.Println("1")
			val, _ := v8go.NewValue(iso, "")
			//_ = info.Context()

			args := info.Args()
			if len(args) <= 0 {
				fmt.Println("2")

				// TODO: this should return an error, but v8go not supported now

				return val
			}
			fmt.Println("3")
			//
			// go func() {
			s := args[0].String()
			fmt.Println("s=>", s)
			retVal := make([]uint8, 0)
			for x := range s {
				retVal = append(retVal, uint8(x))
			}
			// str := ""
			// for i := 0; i < len(s); i++ {
			// 	fmt.Printf("%d ", s[i])
			// 	if str == "" {
			// 		str = "["
			// 	} else {
			// 		str = str + ", "
			// 	}
			// 	str = str + fmt.Sprintf("%d", s[i])
			// }
			// str = str + "]"
			v, err := v8go.NewValue(iso, retVal)
			if err != nil {
				fmt.Println("error creating new val ", err)
				return val
			}

			// obj, err := v.AsObject()
			// if err != nil {
			// 	fmt.Println("error conveting to object ", err)
			// 	return val
			// }

			// resolver, err := v8go.NewPromiseResolver(ctx)
			// if err != nil {
			// 	fmt.Println("error creating NewPromiseResolver ", err)
			// 	return val
			// }
			// resolver.Resolve(v)
			// // }()
			// fmt.Println("5")
			// return resolver.GetPromise().Value

			return v
		})
		val, _ := v8go.NewValue(iso, "")

		if err != nil {
			log.Printf("error creating Encode() template: %#v", err)
			return val
		}
		resTmp, err := v8go.NewObjectTemplate(iso)
		if err != nil {
			log.Printf("error creating object template: %#v", err)
			return val
		}

		fmt.Println("set encode")
		if err := resTmp.Set("encode", encodeFnTmp, v8go.ReadOnly); err != nil {
			log.Printf("error setting encode function template: %#v", err)
			return val
		}

		// ctx, err = v8go.NewContext(iso)
		// if err != nil {
		// 	log.Printf("error new context from iso: %#v", err)
		// 	return val
		// }
		resObj, err := resTmp.NewInstance(ctx)
		if err != nil {
			log.Printf("error new instance from ctx: %#v", err)
			return val
		}

		//resolver.Resolve(resObj)
		fmt.Println("resolved")
		//	return resolver.GetPromise().Value
		return resObj.Value
	}
}
