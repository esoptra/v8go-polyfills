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

	"github.com/esoptra/v8go"
)

func InjectWith(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...Option) error {
	e := NewDecode(opt...)
	decodeFnTmp, err := v8go.NewFunctionTemplate(iso, e.TextDecoderFunctionCallback())
	if err != nil {
		return fmt.Errorf("v8go-polyfills/textDecoder NewFunctionTemplate: %w", err)
	}
	if err := global.Set("TextDecoder", decodeFnTmp); err != nil {
		return fmt.Errorf("v8go-polyfills/textDecoder global.set: %w", err)
	}

	return nil
}
