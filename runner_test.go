package polyfills

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/esoptra/v8go"
	"github.com/esoptra/v8go-polyfills/console"
	"github.com/esoptra/v8go-polyfills/fetch"
	"github.com/esoptra/v8go-polyfills/uuid"

	"github.com/esoptra/v8go-polyfills/crypto"
)

func TestRunPromise(t *testing.T) {
	addr := "localhost:10001"
	go fetch.StartHttpServer(addr)
	time.Sleep(time.Second * 5)

	runTestRunner(t, `epsilon = async (event) => {
		const headers = new Headers()
        headers.append("Authorization", "Basic dGVzdDp0ZXN0");
        headers.append("X-Client-Name", "twintag.js");
        headers.append("X-Client-Version", "12.3.4");

		let options = {
			method: 'GET',
        	headers: headers,
        	body: 'hello twintag',
			redirect: 'manual'
		}
		let path = "http://127.0.0.1:10001/auth"
		console.log('calling fetch', path)
		const resp = await fetch(path, options);
		return new Response(resp)
		}
		let res = epsilon();
		`)

	runTestRunner(t, `epsilon = async (event) => {
		let hed = new Headers()
		hed.append("X-Client-Name", "twintag.js")
		console.log(hed.get('X-Client-Name'))
		const url = 'http://127.0.0.1:10001/'
		let resp = await fetch(url)
		let respText = await resp.text()
		return new Response(respText)
	}
	let res = epsilon();
	`)

	runTestRunner(t, `epsilon = async (event) => {
		const headers = new Headers()
        headers.append("Authorization", "Basic dGVzdDp0ZXN0");
        headers.append("X-Client-Name", "twintag.js");
        headers.append("X-Client-Version", "12.3.4");

		let options = {
			method: 'GET',
        	headers: headers,
        	body: 'hello twintag',
			redirect: 'manual'
		}
		let path = "http://127.0.0.1:10001/auth"
		console.log('calling fetch', path)
		const resp = await fetch(new Request(path, options));
		return new Response(resp)
		}
		let res = epsilon();
		`)

	runTestRunner(t, `epsilon = async (event) => {
		let options = {
			method: 'GET',
        	headers: null,
			redirect: ''
		}
		let path = "http://127.0.0.1:10001/auth"
		console.log('calling fetch', path)
		const resp = await fetch(new Request(path, options));
		return new Response(resp)
		}
		let res = epsilon();
		`)
}

func runTestRunner(t *testing.T, script string) {

	ctxCn, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()

	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()
	global, _ := v8go.NewObjectTemplate(iso)

	fetcher := fetch.NewFetcher()
	if err := fetch.InjectWithFetcherTo(iso, global, fetcher); err != nil {
		t.Error(err)
		return
	}

	runner, err := NewRunner(iso, global)
	if err != nil {
		t.Error(err)
		return
	}

	ctx, err := v8go.NewContext(iso, global)
	if err != nil {
		t.Error(err)
		return
	}
	if err := crypto.InjectTo(ctx); err != nil {
		panic(err)
	}

	if err := console.InjectTo(ctx); err != nil {
		panic(err)
	}

	if err := fetch.InjectHTTPProperties(ctx); err != nil {
		panic(err)
	}

	val, err := runner.RunPromise(ctxCn, ctx, script)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := val.AsObject()
	if err != nil {
		t.Error(err)
		return
	}
	status, err := res.Get("status")
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println("status : ", status.String())
	body, err := res.Get("body")
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println("body : ", body.String())
	if uuid.IsUUID(body.String()) {
		val, ok := fetcher.ResponseMap.Load(body.String())
		if ok {
			result := val.(io.ReadCloser)
			defer result.Close()
			bodyBytes, err := ioutil.ReadAll(result)
			if err != nil {
				t.Error(fmt.Errorf("Error while getting status of epsilon execution : %#v", err))
				return
			}
			fmt.Printf("status %s, bodyBytes %s", status.String(), string(bodyBytes))
			return
		}
	}
	fmt.Printf("status %s, body %s", status.String(), body.String())
}
