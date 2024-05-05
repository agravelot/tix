package main

import (
	"log"
	"os"

	"github.com/agravelot/tix/app"
	"github.com/agravelot/tix/ui"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	// TODO make it configurable
	f, err := os.ReadFile("config/workspaces.toml")
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	var cfg app.Config
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}

	err = ui.New(cfg)
	if err != nil {
		log.Fatalf("error on ui: %v", err)
	}
}
