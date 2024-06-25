package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/agravelot/tix/core"
	"github.com/agravelot/tix/ui"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	// TODO make it configurable
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".config", "tix", "config.toml")
	log.Printf("config path: %s", configPath)

	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	var cfg core.Config
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}

	a, err := core.NewApplication(cfg)
	if err != nil {
		log.Fatalf("error on application: %v", err)
	}

	err = ui.New(a)
	if err != nil {
		log.Fatalf("error on ui: %v", err)
	}
}
