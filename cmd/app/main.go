package main

import "github.com/CyganFx/table-reservation/internal/app"

const (
	configsDir   = "C:/Users/alemh/go/src/github.com/Alemkhan/table-reservation/configs/main.yaml"
	templatesDir = "C:/Users/alemh/go/src/github.com/Alemkhan/table-reservation/ui/html"
)

func main() {
	app.Run(configsDir, templatesDir)
}
