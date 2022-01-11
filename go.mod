module github.com/esoptra/v8go-polyfills

go 1.16

require (
	github.com/esoptra/v8go v0.6.1-0.20220110104324-494e51808509
	github.com/lestrrat-go/jwx v1.2.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/text v0.3.6
)

// replace github.com/esoptra/v8go => ../v8go

retract [v0.1.0, v0.3.0]
