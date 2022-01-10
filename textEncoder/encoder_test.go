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

	"github.com/esoptra/v8go"
	"github.com/esoptra/v8go-polyfills/console"
)

func TestInject(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(iso)

	err := InjectTo(iso, global)
	if err != nil {
		t.Error(err)
	}

	ctx := v8go.NewContext(iso, global)
	if err := console.InjectTo(ctx); err != nil {
		t.Error(err)
	}

	val, err := ctx.RunScript(`const encoder = new TextEncoder()
	console.log(typeof encoder.encode);
	const view = encoder.encode('eyJhdWQiOiJlNjI0YTdjMi0zZTUzLTQ2NTktOGY5Yi1kN2MxOWZjZjAxZjciLCJpc3MiOiJodHRwczovL2xvZ2luLm1pY3Jvc29mdG9ubGluZS5jb20vMjRiMDgwY2QtNTg3NC00NGFiLTk4NjItOGQ3ZTBlMDc4MWFiL3YyLjAiLCJpYXQiOjE2MzkxMzU3NDYsIm5iZiI6MTYzOTEzNTc0NiwiZXhwIjoxNjM5MTM5NjQ2LCJuYW1lIjoiQXNoaXNoIFNoYXJtYSAoRGV2T24pIiwib2lkIjoiOTZmODM2N2QtY2M2NC00NjMwLWI0MGQtYTUwNTVjMjAwOGVkIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYXNoaXNoLnNoYXJtYUBkZXZvbi5ubCIsInJoIjoiMC5BUUlBellDd0pIUllxMFNZWW8xLURnZUJxOEtuSk9aVFBsbEdqNXZYd1pfUEFmY0NBTzguIiwic3ViIjoiLVNDRE5lR2IwVVc1TzZ5NkoxMERyNWhFZWxIR0lSdU5uNnd3NTZuMHRyMCIsInRpZCI6IjI0YjA4MGNkLTU4NzQtNDRhYi05ODYyLThkN2UwZTA3ODFhYiIsInV0aSI6InJfZERnTWdncGtPN01pQnhPNndTQUEiLCJ2ZXIiOiIyLjAifQ')
	console.log("=>", view); 
	console.log(typeof view); //expecting the type as Object (uint8Array)

	const utf8 = new ArrayBuffer(7);
	let encodedResults = encoder.encodeInto('H€llo', utf8);
	console.log("=>", utf8, encodedResults.read, encodedResults.written); 
	console.log(typeof utf8); //expecting the type as Object (uint8Array)

	utf8 `, "encoder.js")
	if err != nil {
		t.Error(err)
	}

	ok := val.IsArrayBuffer()
	if ok {
		bytes := val.ArrayBuffer().GetBytes()
		fmt.Println("returned val is array", bytes)
	} else {
		fmt.Println("returned val is not array", val.Object().Value)
	}

}
