module github.com/esoptra/v8go-polyfills

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/esoptra/v8go v0.6.1-0.20210618103653-571c0d9060e5
	github.com/lestrrat-go/jwx v1.2.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/text v0.3.6
)

//replace github.com/esoptra/v8go => ../v8go

retract [v0.1.0, v0.3.0]
