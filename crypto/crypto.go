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
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/esoptra/v8go"
	"github.com/esoptra/v8go-polyfills/uuid"
)

type Crypto struct {
	KeyMap sync.Map
}

func NewCrypto(opt ...Option) *Crypto {
	c := &Crypto{}

	for _, o := range opt {
		o.apply(c)
	}

	return c
}

//cryptoVerifyFunctionCallback implements https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/verify
//const result = crypto.subtle.verify(algorithm, key, signature, data);
//result is a Promise with a Boolean: true if the signature is valid, false otherwise.
func (c *Crypto) cryptoVerifyFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		iso, _ := ctx.Isolate()
		resolver, err := v8go.NewPromiseResolver(ctx)
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error creating newPromiseResolver with ctx: %#v", err))
		}
		go func() {
			passed := false
			defer func() {
				v, err := v8go.NewValue(iso, passed)
				if err != nil {
					resolver.Reject(newErrorValue(iso, "Error creating newvalue: %#v\n", err))
					return
				}
				resolver.Resolve(v)
			}()
			args := info.Args()
			if len(args) != 4 {
				resolver.Reject(newErrorValue(iso, "Expected algorithm, key, signature, data (4) arguments\n"))
				return
			}

			algorithm, algoName, err := getAlgorithm(args[0])
			if err != nil {
				resolver.Reject(newErrorValue(iso, "Error parsing Algorithm arg: %#v\n", err))
				return
			}
			_ = algorithm //Need this in future
			keyBytes, err := args[1].MarshalJSON()
			if err != nil {
				resolver.Reject(newErrorValue(iso, "Error marshalling key arg: %#v\n", err))
				return
			}

			if !args[2].IsUint8Array() {
				resolver.Reject(newErrorValue(iso, "Expecting signature in []byte(ArrayBuffer) format\n"))
				return
			}

			key := &CryptoKey{}
			err = json.Unmarshal(keyBytes, key)
			if err != nil {
				resolver.Reject(newErrorValue(iso, "Error unmarshalling keybytes arg: %#v\n", err))
				return
			}

			if algoName == "RSA-OAEP" && key.Type == "public" {
				fmt.Println("extract pub key", key.Kid)
				//this expecting public rsa key
				pubKey, ok := c.KeyMap.Load(key.Kid)
				if !ok {
					resolver.Reject(newErrorValue(iso, "Invalid Key : %#v\n", key))
					return
				}
				sign := args[2].Uint8Array()
				fmt.Println("sign", string(sign))
				payload := args[3].Uint8Array()
				fmt.Println("payload", string(payload))

				err = jwt.SigningMethodRS256.Verify(string(payload), string(sign), pubKey)
				passed = (err == nil)

			}

		}()

		return resolver.GetPromise().Value
	}
}

//cryptoImportKeyFunctionCallback implements https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/importKey
//const result = crypto.subtle.importKey(format, keyData, algorithm, extractable, keyUsages);
//result is a Promise that fulfills with the imported key as a CryptoKey object.
func (c *Crypto) cryptoImportKeyFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		iso, _ := ctx.Isolate()
		resolver, err := v8go.NewPromiseResolver(ctx)
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error creating newPromiseResolver with ctx: %#v", err))
		}
		go func() {
			args := info.Args()

			format := args[0].String()
			var result interface{}
			if format == "jwk" {
				keyData, err := args[1].AsObject() //object type
				if err != nil {
					resolver.Reject(newErrorValue(iso, "error getting keyData %#v", err))
					return
				}
				algorithm, algoName, err := getAlgorithm(args[2])
				if err != nil {
					resolver.Reject(newErrorValue(iso, "Error parsing Algorithm arg: %#v\n", err))
					return
				}

				extractable := args[3] //boolean
				if !extractable.IsBoolean() {
					resolver.Reject(newErrorValue(iso, "Expected extractable argument as boolean type\n"))
					return
				}

				keyUsages := args[4] //array
				if !keyUsages.IsArray() {
					resolver.Reject(newErrorValue(iso, "Expected keyUsages argument as array type\n"))
					return
				}

				keyDataBytes, err := keyData.MarshalJSON()
				if err != nil {
					resolver.Reject(newErrorValue(iso, "error marshalling keyData %#v", err))
					return
				}
				//fmt.Println(string(keyDataBytes))

				isKeySet := keyData.Has("keys")

				var key interface{}
				if algoName == "RSA-OAEP" {
					fmt.Println("iskeyset", isKeySet)
					//this expecting public rsa key
					if isKeySet {
						keys, err := parseKeySet(keyDataBytes)
						if err != nil {
							resolver.Reject(newErrorValue(iso, "Could not parse DER encoded key (encryption key): %#v", err))
							return
						}
						//select the first key from the set
						key = keys[0]
						fmt.Println("public key size==>", keys[0].Size())
					} else {
						key, err = parseKey(keyDataBytes)
						if err != nil {
							resolver.Reject(newErrorValue(iso, "Could not parse DER encoded key (encryption key): %#v", err))
							return
						}
					}

				}

				fmt.Println(key)
				miniPub := uuid.NewUuid()
				c.KeyMap.Store(miniPub, key)

				result = &CryptoKey{
					Type:        "public",
					Kid:         miniPub,
					Extractable: extractable.Boolean(),
					Algorithm:   algorithm,
					Usages:      keyUsages,
				}
			} else {
				resolver.Reject(newErrorValue(iso, "format %q not supported", format))
				return
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				resolver.Reject(newErrorValue(iso, "error marshalling jsonKey: %#v", err))
				return
			}
			v, err := v8go.JSONParse(info.Context(), string(resultBytes))
			if err != nil {
				resolver.Reject(newErrorValue(iso, "error jsonParse on result: %#v", err))
				return
			}

			resolver.Resolve(v)
		}()

		return resolver.GetPromise().Value
	}
}

type RSAAlgo struct {
	Name           string           `json:"name"`           //"RSA-OAEP",
	ModulusLength  int              `json:"modulusLength"`  // 4096,
	PublicExponent map[string]uint8 `json:"publicExponent"` // new Uint8Array([1, 0, 1]),
	Hash           string           `json:"hash"`           // "SHA-256"
}

//for symmetric algo https://developer.mozilla.org/en-US/docs/Web/API/CryptoKey
type CryptoKey struct {
	Type        string      `json:"type"`
	Kid         string      `json:"kid"` //additional property to refer the actual key
	Extractable bool        `json:"extractable"`
	Algorithm   interface{} `json:"algorithm"`
	Usages      interface{} `json:"usages"`
}

//for public-key algorithms
type CryptoKeyPair struct {
	PrivateKey CryptoKey `json:"privateKey"`
	PublicKey  CryptoKey `json:"publicKey"`
}

//cryptoGenerateKeyFunctionCallback implements https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/generateKey
//const result = crypto.subtle.generateKey(algorithm, extractable, keyUsages);
//result is a Promise that fulfills with a CryptoKey (for symmetric algorithms) or a CryptoKeyPair (for public-key algorithms).
func (c *Crypto) cryptoGenerateKeyFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		iso, _ := ctx.Isolate()
		resolver, err := v8go.NewPromiseResolver(ctx)
		if err != nil {
			return iso.ThrowException(fmt.Sprintf("error creating newPromiseResolver with ctx: %#v", err))
		}
		go func() {

			// if err != nil {
			// 	return iso.ThrowException(fmt.Sprintf("error getting Isolate from ctx: %#v", err))
			// }
			args := info.Args()
			if len(args) != 3 {
				resolver.Reject(newErrorValue(iso, "Expected algorithm, extractable, keyUsages (3) arguments\n"))
				return
			}

			algorithm, algoName, err := getAlgorithm(args[0])
			if err != nil {
				resolver.Reject(newErrorValue(iso, "Error parsing Algorithm arg: %#v\n", err))
				return
			}

			extractable := args[1] //boolean
			if !extractable.IsBoolean() {
				resolver.Reject(newErrorValue(iso, "Expected extractable argument as boolean type\n"))
				return
			}

			keyUsages := args[2] //array
			if !keyUsages.IsArray() {
				resolver.Reject(newErrorValue(iso, "Expected keyUsages argument as array type\n"))
				return
			}

			var result interface{}

			if algoName == "RSA-OAEP" {
				primeBits := 2048

				rsaAlgo := algorithm.(*RSAAlgo)
				if rsaAlgo.ModulusLength != 0 {
					primeBits = rsaAlgo.ModulusLength
				}
				// The GenerateKey method takes in a reader that returns random bits, and
				// the number of bits
				privateKey, err := rsa.GenerateKey(rand.Reader, primeBits) //2048 by default
				if err != nil {
					resolver.Reject(newErrorValue(iso, "error generating RSA key: %#v", err))
					return
				}

				fmt.Printf("%v\n", algorithm)
				// The public key is a part of the *rsa.PrivateKey struct
				publicKey := privateKey.PublicKey
				fmt.Printf("%#v", publicKey)

				// jsonKey := jose.JSONWebKey{
				// 	Key:       privateKey,
				// 	Algorithm: algoName,
				// }

				// jsonKeyBytes, err := jsonKey.MarshalJSON()
				// if err != nil {
				// 	resolver.Reject(newErrorValue(iso, "error marshalling jsonKey: %#v", err))
				// 	return
				// }
				//store a pointer reference with the fetcher
				miniPriv := uuid.NewUuid()
				c.KeyMap.Store(miniPriv, privateKey)
				miniPub := uuid.NewUuid()
				c.KeyMap.Store(miniPub, &privateKey.PublicKey)

				result = &CryptoKeyPair{
					PrivateKey: CryptoKey{
						Type:        "private",
						Kid:         miniPriv,
						Extractable: extractable.Boolean(),
						Algorithm:   algorithm,
						Usages:      keyUsages.Object(),
					},
					PublicKey: CryptoKey{
						Type:        "public",
						Kid:         miniPub,
						Extractable: extractable.Boolean(),
						Algorithm:   algorithm,
						Usages:      keyUsages.Object(),
					},
				}
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				resolver.Reject(newErrorValue(iso, "error marshalling jsonKey: %#v", err))
				return
			}
			v, err := v8go.JSONParse(info.Context(), string(resultBytes))
			if err != nil {
				resolver.Reject(newErrorValue(iso, "error jsonParse on result: %#v", err))
				return
			}

			// v, err := v8go.NewValue(iso, result)
			// if err != nil {
			// 	resolver.Reject(newErrorValue(iso, "error new value for generateKey: %#v", err))
			// 	return
			// }

			resolver.Resolve(v)
		}()

		return resolver.GetPromise().Value
	}
}

func newErrorValue(iso *v8go.Isolate, format string, a ...interface{}) *v8go.Value {
	e, _ := v8go.NewValue(iso, fmt.Sprintf(format, a...))
	return e
}

func getAlgorithm(v *v8go.Value) (interface{}, string, error) {
	algorithm, err := v.AsObject() //object type
	if err != nil {
		return nil, "", fmt.Errorf("Expected algorithm argument as Object type: %#v\n", err)
	}
	if !algorithm.Has("name") {
		return nil, "", fmt.Errorf("Missing algorithm's name property\n")
	}
	algoName, err := algorithm.Get("name")
	if err != nil {
		return nil, "", fmt.Errorf("Missing algorithm's name property:%#v\n", err)
	}
	if algoName.String() == "" {
		return nil, "", fmt.Errorf("Missing algorithm's name property value\n")
	}
	res, err := algorithm.MarshalJSON()
	if err != nil {
		return nil, "", fmt.Errorf("Error Marshalling algorithm:%#v\n", err)
	}

	var result interface{}
	if algoName.String() == "RSA-OAEP" {
		rsa := &RSAAlgo{}
		fmt.Println(string(res))
		err = json.Unmarshal(res, rsa)
		if err != nil {
			return nil, "", fmt.Errorf("Error UnMarshalling algorithm:%#v\n", err)
		}
		result = rsa
	}

	return result, algoName.String(), nil
}
