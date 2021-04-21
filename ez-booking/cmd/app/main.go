package main

import "github.com/CyganFx/table-reservation/ez-booking/internal/app"

const HttpPort = ":5001"

// logic in ./internal/app/app.go
func main() {
	app.Run(HttpPort)
}
