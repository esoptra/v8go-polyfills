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

package textDecoder

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"

	"github.com/esoptra/v8go"
)

type Decoder struct {
}

func NewDecode(opt ...Option) *Decoder {
	c := &Decoder{}

	for _, o := range opt {
		o.apply(c)
	}

	return c
}

//Encoding struct resembles internal.Encoding in "golang.org/x/text/encoding/internal"
//this is to avoid import of internal package
type Encoding struct {
	encoding.Encoding
	Name string
	MIB  uint16
}

func getDecoder(label string) (*encoding.Decoder, error) {

	if label == "" {
		//utf-8 by default
		return nil, nil
	}
	sanitize := func(input string) string {
		x := strings.ToUpper(strings.ReplaceAll(input, "-", " "))
		//fmt.Println("sanitized", x)
		return x
	}
	sanitizedlabel := sanitize(label)
	//fmt.Println("sanitized label", sanitizedlabel)
	for _, eachEncoding := range charmap.All {

		name := ""
		cast, ok := eachEncoding.(*charmap.Charmap)
		if !ok {
			castInternal, ok := eachEncoding.(*Encoding)
			if !ok {
				continue
			} else {
				name = castInternal.Name
			}
		} else {
			name = cast.String()
		}
		if sanitize(name) == sanitizedlabel {
			return eachEncoding.NewDecoder(), nil
		}
	}

	return nil, fmt.Errorf("Not supported label %q", label)
}

func (c *Decoder) TextDecoderFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		iso := ctx.Isolate()
		label := ""
		_ = label
		if len(info.Args()) > 0 {
			label = info.Args()[0].String()
			if len(info.Args()) > 1 {
				fmt.Printf("Options %q Not yet supported\n", info.Args()[1].String())
			}
		}

		decodeFnTmp := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) <= 0 {
				strErr, _ := v8go.NewValue(iso, "Expected an arguments\n")

				return iso.ThrowException(strErr)
			}
			s := args[0].Uint8Array()
			result := ""

			dec, err := getDecoder(label)
			if err != nil {
				strErr, _ := v8go.NewValue(iso, err.Error())
				return iso.ThrowException(strErr)
			}
			if dec != nil {
				bUTF := make([]byte, len(s)*3)
				n, _, err := dec.Transform(bUTF, s, true)
				if err != nil {
					strErr, _ := v8go.NewValue(iso, fmt.Sprintf("Error transforming: %#v", err))
					return iso.ThrowException(strErr)
				}
				result = string(bUTF[:n])
			} else {
				result = string(s)
			}

			//fmt.Println(result)
			v, err := v8go.NewValue(iso, result)
			if err != nil {
				strErr, _ := v8go.NewValue(iso, fmt.Sprintf("error creating new val: %#v", err))

				return iso.ThrowException(strErr)
			}
			return v
		})

		resTmp := v8go.NewObjectTemplate(iso)

		if err := resTmp.Set("decode", decodeFnTmp, v8go.ReadOnly); err != nil {
			strErr, _ := v8go.NewValue(iso, fmt.Sprintf("error setting encode function template: %#v", err))
			return iso.ThrowException(strErr)
		}

		resObj, err := resTmp.NewInstance(ctx)
		if err != nil {
			strErr, _ := v8go.NewValue(iso, fmt.Sprintf("error new instance from ctx: %#v", err))
			return iso.ThrowException(strErr)
		}
		return resObj.Value
	}
}
