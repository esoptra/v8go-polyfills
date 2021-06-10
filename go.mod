module github.com/esoptra/v8go-polyfills

go 1.16

require (
	github.com/esoptra/v8go v0.6.1-0.20210610071544-c13ab7e9247a
	golang.org/x/text v0.3.6
)

//replace github.com/esoptra/v8go => ../v8go

retract [v0.1.0, v0.3.0]
