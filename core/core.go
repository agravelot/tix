package core

import (
	"errors"
	"log"
	"os"

	"github.com/agravelot/tix/workspace"
)

// Config represents the configuration of tix
type Config struct {
	Shell      string
	Workspaces []ConfigWorkspace `toml:"workspace"`
}

type ConfigWorkspace struct {
	Name             string
	Directory        string
	SetupCommands    []string
	TeardownCommands []string
	// Default to 5 seconds
	Timeout int
}

type Application struct {
	Config     Config
	Workspaces []workspace.Workspace
}

func NewApplication(cfg Config) (Application, error) {
	workspaces := make([]workspace.Workspace, len(cfg.Workspaces))
	for _, cf := range cfg.Workspaces {
		w := workspace.Workspace{
			Name: cf.Name,
		}

		workspaces = append(workspaces, w)
	}

	if cfg.Shell == "" {
		log.Println("No shell defined, using default from SHELL eenvironment variable")
		cfg.Shell = os.Getenv("SHELL")
		if cfg.Shell == "" {
			return Application{}, errors.New("no shell defined")
		}
	}

	return Application{
		Config:     cfg,
		Workspaces: workspaces,
	}, nil
}
