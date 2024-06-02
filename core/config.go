package core

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
