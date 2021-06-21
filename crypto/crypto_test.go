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

package crypto

import (
	"fmt"
	"testing"

	"github.com/esoptra/v8go"
	"github.com/esoptra/v8go-polyfills/console"
	"github.com/esoptra/v8go-polyfills/fetch"
	"github.com/esoptra/v8go-polyfills/textEncoder"
)

func TestCrypto(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()

	con, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		t.Error(err)
		return
	}
	if err := fetch.InjectTo(iso, con); err != nil {
		t.Error(err)
		return
	}
	ctx, err := textEncoder.InjectWith(iso, con)
	if err != nil {
		t.Error(err)
	}

	// ctx, err := v8go.NewContext(iso, con)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	if err := InjectWith(iso, ctx); err != nil {
		t.Error(err)
		return
	}

	if err := console.InjectTo(ctx); err != nil {
		t.Error(err)
	}

	val, err := ctx.RunScript(`const fetchKeys = async () => {
			  const data = await fetch('https://login.microsoftonline.com/24b080cd-5874-44ab-9862-8d7e0e0781ab/v2.0/.well-known/openid-configuration')
			  const jwksResp = await data.json();
			  //console.log(jwksResp.jwks_uri)
			  const keysData = await fetch(jwksResp.jwks_uri);
			  return await keysData.json();
			};

			Raj();
			Ra();
			epsilon = async (event) => {
	
 	let data = await fetchKeys()

	const algo = {
		name: "RSA-OAEP",
		modulusLength: 4096,
		publicExponent: new Uint8Array([1, 0, 1]),
		hash: "SHA-256"
	  };
	let importedKey = await crypto.subtle.importKey('jwk', data, algo, true, ["encrypt", "decrypt"]);
	console.log(importedKey.kid);

	const token = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Im5PbzNaRHJPRFhFSzFqS1doWHNsSFJfS1hFZyJ9.eyJhdWQiOiIwOGQ0NWY3Zi0xNmM5LTQ1ZGUtYmFkZC05NDc3ZGRjZTVlMzYiLCJpc3MiOiJodHRwczovL2xvZ2luLm1pY3Jvc29mdG9ubGluZS5jb20vZWMxMDAyZDctMDM0OC00MGFlLWFlNGUtOTBjMDA1MDZlYWNkL3YyLjAiLCJpYXQiOjE2MjM3NDMzNDcsIm5iZiI6MTYyMzc0MzM0NywiZXhwIjoxNjIzNzQ3MjQ3LCJhaW8iOiJBV1FBbS84VEFBQUFtbEJHSkpMT25BbllSMFRTVFUvdS9NQzhQNDdTSWJncWltY0xva1B0OFpLbGpvdW5IdmJ3M0d4UHJqTnlWNDkrbW5JTUhUNW84R2tXSHprelN4Qk5QSE0vWmVDWFdSc0FITmxtNW1oTURaWEMvWkJ5N3UwT0xVbTMrU0ZBSWR0NiIsImlkcCI6Imh0dHBzOi8vc3RzLndpbmRvd3MubmV0LzU2OTc1MzBiLWRmYzMtNDhlZi1iMWFjLTk4ZTRmZTY4MTI1Mi8iLCJuYW1lIjoiYmFja2tlbSIsIm9pZCI6IjIyNjZhMGYxLWFiZWYtNGRmNC1hY2UwLTNhZDk4NTcxOWRjMCIsInByZWZlcnJlZF91c2VybmFtZSI6Ik1pY2hpZWwuZGViYWNra2VyQHR3aW50YWcuY29tIiwicmgiOiIwLkFRd0Exd0lRN0VnRHJrQ3VUcERBQlFicXpYOWYxQWpKRnQ1RnV0MlVkOTNPWGpZTUFLVS4iLCJzdWIiOiJRR1BSamZrWUhaLTl2MFlrU2lQdm1YX3BhQTAzYzRDbGZrcUlkQWpoMDFvIiwidGlkIjoiZWMxMDAyZDctMDM0OC00MGFlLWFlNGUtOTBjMDA1MDZlYWNkIiwidXRpIjoiOHRaWlB1WHdQVUtPRUVLOGhraWxBQSIsInZlciI6IjIuMCJ9.j58zhFkqOPtcxB-gA1LdLYJYQw_oVZ2vDiZXD6M9nZNWbgAmFFkvN7CuhQFYR5rM9XaGrO-Rn4X6X389aFk-sZKQUOtVqmW4VT8_yT2iSGVspL5BcwWYeR0vEjO_5UNoavSunXz_qOFzzQqUYZ2-ex3KG9x7cL1Tc1kVv2JmAtUB-yK5t5yZU1BzNteIDCC4QEUa_vBxZrTwVEkRW_fT26TonWZTikYvi80COSFlMRiDD-gK2QFHrjcyPvhETTYDzXYhHoJDolcey59ERu9301SE9flTMigVpJlL5SreMIWhy1-vWt5lbCPOA246o3hEa_HAmAVgIdC1t1tSsj61hw'
	const splitToken = token.split('.')

	const encoder = new TextEncoder()
	const signature = encoder.encode(splitToken[2])
	const payload = encoder.encode(splitToken[0]+'.'+splitToken[1])

	let isvalid = await crypto.subtle.verify(algo, importedKey, signature, payload)
	
	let returnVal = 'failed'
	if (isvalid){
		returnVal = 'success'
	}
	returnVal;
	};
	let res = epsilon();
	Promise.resolve(res)`, "crypto.js")
	if err != nil {
		t.Error(err)
	}

	proms, err := val.AsPromise()
	if err != nil {
		t.Error(err)
		return
	}

	for proms.State() == v8go.Pending {
		continue
	}

	res := proms.Result().String()

	fmt.Println("returned val is", res)
}

func TestInject(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)

	// global, err := v8go.NewObjectTemplate(iso)
	// if err != nil {
	// 	t.Error(err)
	// }

	if err := InjectTo(ctx); err != nil {
		t.Error(err)
	}

	// ctx1, err := v8go.NewContext(iso, global)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// if err := console.InjectTo(ctx1); err != nil {
	// 	t.Error(err)
	// }

	val, err := ctx.RunScript(`let view = await crypto.verify('test')
	console.log("=>", view); 
	view `, "encoder.js")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("returned val is", val.String())

}
