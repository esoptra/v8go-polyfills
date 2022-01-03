/*
 * Copyright (c) 2021 Xingwang Liao
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

package base64

import (
	"testing"

	"github.com/esoptra/v8go"
)

func TestAtob(t *testing.T) {
	ctx, err := newV8goContext()
	if err != nil {
		t.Error(err)
		return
	}

	val, err := ctx.RunScript("atob()", "atob_undefined.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('')", "atob_empty.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('MTIzNA==')", "atob_1234.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "1234" {
		t.Errorf("assert '1234' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('5rGJ5a2X')", "atob_unicode.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "汉字" {
		t.Errorf("assert '汉字' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("atob('pG8hmgZAFrtaqOIGgc2eYkZo/Xxb0/0ntgky5fZsAg5QpHoh7f6EfufNaTEcGY4oFMcE9ii5TI5S9+0LJqJcHRHpOd8xcWMQ+YWe5DAEag/90uRubFXDuLnK4mHGrEWRdlkO0vp5YDmhUItykyq/GVMzwBmbRKhRWzVxEao9dsXFZrnrTkQE2rdtE81w5kAvhEnYB8q7Yfy0uRN+7U2wLyazy/TqYfx19tDBN66F32rlV8SdTxvewRj4ZMw12/RBLdaiSoj6phMpjOllFgLdUi1RY+roJjNcdaHH988aopZgqnQvTQ8zOczES0tBqt+oRJyrwOhFvpvTZAgg5+ykUg')", "atob_unicode.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "汉字" {
		t.Errorf("assert '汉字' but got '%s'", s)
		return
	}
}

func TestBtoa(t *testing.T) {
	ctx, err := newV8goContext()
	if err != nil {
		t.Error(err)
		return
	}

	val, err := ctx.RunScript("btoa()", "btoa_undefined.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('')", "atob_empty.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "" {
		t.Errorf("assert '' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('1234')", "btoa_1234.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "MTIzNA==" {
		t.Errorf("assert 'MTIzNA==' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa('汉字')", "btoa_unicode.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "5rGJ5a2X" {
		t.Errorf("assert '5rGJ5a2X' but got '%s'", s)
		return
	}

	val, err = ctx.RunScript("btoa({})", "btoa_object.js")
	if err != nil {
		t.Error(err)
		return
	}

	if s := val.String(); s != "W29iamVjdCBPYmplY3Rd" {
		t.Errorf("assert 'W29iamVjdCBPYmplY3Rd' but got '%s'", s)
		return
	}
}

func newV8goContext() (*v8go.Context, error) {
	iso := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(iso)

	if err := InjectTo(iso, global); err != nil {
		return nil, err
	}

	return v8go.NewContext(iso, global), nil
}
