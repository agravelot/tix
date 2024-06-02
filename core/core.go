package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/docker/compose/v2/pkg/api"
)

type Application struct {
	// TODO public
	Docker     api.Service
	Config     Config
	Workspaces []Workspace
}

func NewApplication(cfg Config) (Application, error) {
	// TODO configurable
	srv, err := createDockerService()
	if err != nil {
		return Application{}, fmt.Errorf("error create docker service: %w", err)
	}

	workspaces := make([]Workspace, 0, len(cfg.Workspaces))

	for _, cf := range cfg.Workspaces {
		w, err := NewWorkspace(cf)
		if err != nil {
			// TODO Wrap error
			return Application{}, fmt.Errorf("invalid workspace: %w", err)
		}

		// TODO non blocking
		w.RefreshApplets(context.Background(), srv)

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
		Docker:     srv,
	}, nil
}
