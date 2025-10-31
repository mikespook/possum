module examples/helloworld

go 1.25

replace github.com/mikespook/possum => ../..

require github.com/mikespook/possum v0.0.0-20251030011011-a93c811be178

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
