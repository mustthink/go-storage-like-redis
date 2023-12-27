package main

import (
	"flag"

	"github.com/mustthink/go-storage-like-redis/config"
	"github.com/mustthink/go-storage-like-redis/internal"
)

func main() {
	configPath := flag.String("config", config.DefaultConfig, "path to config")
	flag.Parse()

	app := internal.NewApplication(*configPath)
	app.Run()
}
