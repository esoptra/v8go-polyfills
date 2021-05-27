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
	"testing"

	"rogchap.com/v8go"

	"github.com/esoptra/v8go-polyfills/console"
)

// func TestInject1(t *testing.T) {
// 	t.Parallel()
// 	iso, _ := v8go.NewIsolate()
// 	ctx, _ := v8go.NewContext(iso)

// 	if err := InjectTo(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	if err := console.InjectTo(ctx); err != nil {
// 		t.Error(err)
// 	}

// 	if _, err := ctx.RunScript(`const encoder = TextEncoder()
// 	const view = encoder.encode('€')
// 	console.log(view); // Uint8Array(3) [226, 130, 172]`, "encoder.js"); err != nil {
// 		t.Error(err)
// 	}
// }

func TestInject(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)
	global, _ := v8go.NewObjectTemplate(iso)

	if err := InjectWithGlobalTo(iso, global); err != nil {
		t.Error(err)
	}

	ctx, err := v8go.NewContext(iso, global)
	if err != nil {
		t.Error(err)
		return
	}
	if err := console.InjectTo(ctx); err != nil {
		t.Error(err)
	}

	val, err := ctx.RunScript(`const encoder = new TextEncoder()
	console.log(typeof encoder.encode);
	const view = encoder.encode('€')
	console.log("=>", view); 
	console.log(typeof view); //expecting the type as Object (uint8Array)
	view `, "encoder.js")
	if err != nil {
		t.Error(err)
	}

	idx, ok := val.ArrayIndex()
	if ok {
		fmt.Println("returned val is array", idx)
	} else {
		fmt.Println("returned val is not array", idx)
	}
}
