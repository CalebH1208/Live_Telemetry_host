module main

go 1.24.1

replace car => ./car

require (
	car v0.0.0-00010101000000-000000000000
	serial v0.0.0-00010101000000-000000000000
	ws v0.0.0-00010101000000-000000000000
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	go.bug.st/serial v1.6.3 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace serial => ./serial

replace ws => ./ws
