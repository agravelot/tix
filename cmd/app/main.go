package main

import (
	"log"
	"os"

	"github.com/agravelot/tix/core"
	"github.com/agravelot/tix/ui"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	// TODO make it configurable
	f, err := os.ReadFile("config/workspaces.toml")
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	var cfg core.Config
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}

	a := core.NewApplication(cfg)

	err = ui.New(a)
	if err != nil {
		log.Fatalf("error on ui: %v", err)
	}
}
