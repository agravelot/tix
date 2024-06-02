package core

import "errors"

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
	Timeout       int
	DockerCompose struct{ Configs []string }
}

func (c ConfigWorkspace) Workspace() (Workspace, error) {
	// TODO Validate
	if c.Name == "" {
		return Workspace{}, errors.New("name is required")
	}

	return Workspace{
		Name:             c.Name,
		Directory:        c.Directory,
		SetupCommands:    c.SetupCommands,
		TeardownCommands: c.TeardownCommands,
		Timeout:          c.Timeout,
		DockerCompose:    c.DockerCompose,
	}, nil
}
