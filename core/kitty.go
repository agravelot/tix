package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"
)

type KittySource struct {
	// Password to connect to the kitty process
	// leave it empty if no password is set
	// more info: https://sw.kovidgoyal.net/kitty/conf/#opt-kitty.remote_control_password
	RemotePassword string
	ConfigPath     string
}

func (s KittySource) ListWindows() ([]Window, error) {
	out, err := exec.Command("kitty", "@", "ls").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing kitty: %w", err)
	}

	wins := []Window{}
	err = json.Unmarshal(out, &wins)
	if err != nil {
		return wins, fmt.Errorf("failed to unmashal ls output: %w", err)
	}

	return wins, nil
}

func (s KittySource) StartProject(name string) error {
	pa := path.Join("./", s.ConfigPath, name)
	log.Printf("Starting project %s", pa)

	out, err := exec.Command("kitty", "--session", pa).Output()
	if err != nil {
		log.Println(string(out))
		return fmt.Errorf("failed to start kitty session %s: %w", pa, err)
	}
	return nil
}

func (s KittySource) StopProject(p Workspace) error {
	return fmt.Errorf("not implemented")
}
