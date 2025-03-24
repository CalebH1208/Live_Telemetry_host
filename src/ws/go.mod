module ws

go 1.24.1

replace car => ../car

require (
	car v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.3
)
