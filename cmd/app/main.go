package main

import "github.com/CyganFx/table-reservation/internal/app"

const (
	configsDir   = "./configs/main.yaml"
	templatesDir = "./ui/html/"
)

func main() {
	app.Run(configsDir, templatesDir)
}
