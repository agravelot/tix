package core

import "github.com/agravelot/tix/workspace"

// Config represents the configuration of tix
type Config struct {
	Workspaces []ConfigWorkspace `toml:"workspace"`
}

type ConfigWorkspace struct {
	Name      string
	Directory string
	Shell     string
	// TODO Define default values
	Timeout          int
	SetupCommands    []string
	TeardownCommands []string
}

type Application struct {
	Config     Config
	Workspaces []workspace.Workspace
}

func NewApplication(cfg Config) Application {
	// TODO create workspaces
	return Application{
		Config: cfg,
	}
}
