package core

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Application struct {
	Config     Config
	Workspaces []Workspace
}

func NewApplication(cfg Config) (Application, error) {
	workspaces := make([]Workspace, 0, len(cfg.Workspaces))

	for _, cf := range cfg.Workspaces {
		w, err := cf.Workspace()
		if err != nil {
			// TODO Wrap error
			return Application{}, fmt.Errorf("invalid workspace: %w", err)
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
