module github.com/lnikon/glcs-cmd

go 1.17

require (
	github.com/go-kit/kit v0.12.0
	github.com/lnikon/glcs v0.0.0-00010101000000-000000000000
)

replace github.com/lnikon/glcs => ./glcs

require (
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
)
