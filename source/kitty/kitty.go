package kitty

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/agravelot/tix/project"
	"github.com/pkg/errors"
)

type KittySource struct {
	// Password to connect to the kitty process
	// leave it empty if no password is set
	// more info: https://sw.kovidgoyal.net/kitty/conf/#opt-kitty.remote_control_password
	RemotePassword string
	ConfigPath     string
}

func (s KittySource) ListProjects() ([]project.Project, error) {
	// listOut, err := exec.Command("tmuxinator", "list").Output()
	// if err != nil {
	// 	return nil, fmt.Errorf("error listing tmuxinator projects: %w", err)
	// }

	out, err := exec.Command("kitty", "@", "ls").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing kitty sessions: %w", err)
	}

	projects := []project.Project{}

	list, err := os.ReadDir(s.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list projects: ")
	}

	wins := []Window{}
	err = json.Unmarshal(out, &wins)
	if err != nil {
		return projects, fmt.Errorf("failed to unmashal ls output: %w", err)
	}

	for _, v := range list {
		p := project.Project{Name: v.Name()}
		for _, w := range wins {
			for _, t := range w.Tabs {
				if t.Title == p.Name {
					p.Opened = true
					// TOODO break more levels
					break
				}
			}
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (s KittySource) StartProject(p project.Project) error {
	pa := path.Join("./", s.ConfigPath, p.Name)
	log.Printf("Starting project %s", pa)

	out, err := exec.Command("kitty", "--session", pa).Output()
	if err != nil {
		log.Println(string(out))
		return fmt.Errorf("failed to start kitty session %s: %w", pa, err)
	}
	return nil
}
