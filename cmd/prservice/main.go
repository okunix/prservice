package main

import (
	"flag"

	"github.com/okunix/prservice/internal/app"
	"github.com/okunix/prservice/internal/app/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "/etc/prservice/prservice.yaml", "yaml config file path")
}

func main() {
	flag.Parse()
	if _, err := config.Read(configPath); err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
