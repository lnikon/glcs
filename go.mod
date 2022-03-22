module github.com/lnikon/glcs

go 1.17

require (
	github.com/go-kit/kit v0.12.0
	github.com/go-kit/log v0.2.0
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.4
)

replace github.com/lnikon/glcs/computation => ./computation

replace github.com/lnikon/glcs/server => ./server

replace github.com/lnikon/glcs/utilities => ./utilities

require github.com/go-logfmt/logfmt v0.5.1 // indirect
