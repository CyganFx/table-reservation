package main

import "github.com/CyganFx/table-reservation/internal/app"

const HttpPort = ":5000"

// logic in ./internal/app/app.go
func main() {
	app.Run(HttpPort)
}
