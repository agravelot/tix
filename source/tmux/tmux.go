package tmux

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/agravelot/tix/project"
)

type TmuxSource struct{}

func (s TmuxSource) ListProjects() ([]project.Project, error) {
	listOut, err := exec.Command("tmuxinator", "list").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing tmuxinator projects: %w", err)
	}

	listSessionOut, err := exec.Command("tmux", "list-sessions").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing tmux sessions: %w", err)
	}

	projects := []project.Project{}

	for _, p := range strings.Fields(string(listOut))[2:] {
		// TODO Find a better way to check if a project is running
		running := strings.Contains(string(listSessionOut)+":", p)
		projects = append(projects, project.Project{Name: p, Selected: running, Opened: running})
	}

	return projects, nil
}

func (s TmuxSource) StartProject(p project.Project) error {
	_, err := exec.Command("tmuxinator", "start", p.Name).Output()
	if err != nil {
		return fmt.Errorf("failed to start project %s: %w", p.Name, err)
	}
	return nil
}
