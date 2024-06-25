package core

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/agravelot/tix/color"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/pkg/api"
)

type Task interface {
	Run() error
}

type ShellTask struct{}

func (t ShellTask) Run() error {
	return nil
}

type DockerComposeTask struct{}

func (t DockerComposeTask) Run() error {
	return nil
}

type Applet struct {
	// TODO enum
	Icon      string
	IsRunning bool
}

// Workspace represents a workspace
type Workspace struct {
	Name             string
	Directory        string
	Shell            string
	SetupCommands    []string
	TeardownCommands []string
	DockerCompose    struct {
		Configs []string
		project *types.Project
	}
	Timeout int
	Applets []Applet
}

func (w Workspace) IsRunning() bool {
	for _, a := range w.Applets {
		if a.IsRunning {
			return true
		}
	}

	return false
}

// TODO inject applet source into workspace
func (w *Workspace) RefreshApplets(ctx context.Context, srv api.Service) error {
	// TODO move it and avoid perf issue
	r := KittySource{
		// TODO make it configurable
		ConfigPath:     "./config",
		RemotePassword: "my passphrase",
	}

	projects, err := r.ListWindows()
	if err != nil {
		return fmt.Errorf("unable listing projects: %w", err)
	}

	applets := []Applet{
		{
			Icon:      "kitty",
			IsRunning: r.HasTabOpen(projects, w.Name),
		},
	}

	if len(w.DockerCompose.Configs) > 0 && w.DockerCompose.project != nil {
		sum, err := srv.Ps(ctx, w.DockerCompose.project.Name, api.PsOptions{All: true})
		if err != nil {
			return fmt.Errorf("unable to list docker compose: %w", err)
		}

		for _, s := range sum {
			applets = append(applets, Applet{
				Icon:      "docker",
				IsRunning: s.State == "running",
			})
		}
	}

	w.Applets = applets

	return nil
}

// TODO Add context
func (w Workspace) Setup(ctx context.Context, srv api.Service) error {
	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	wg := sync.WaitGroup{}

	// TODO Refactor outside of workspace
	if len(w.DockerCompose.Configs) > 0 && w.DockerCompose.project != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createOpts := api.CreateOptions{
				Recreate:             api.RecreateDiverged,
				RecreateDependencies: api.RecreateDiverged,
				Inherit:              true,
				Services:             w.DockerCompose.project.ServiceNames(),
				Build: &api.BuildOptions{
					Services: w.DockerCompose.project.ServiceNames(),
				},
			}

			startOpts := api.StartOptions{
				Project: w.DockerCompose.project,
			}

			err := srv.Up(ctx, w.DockerCompose.project, api.UpOptions{
				Create: createOpts,
				Start:  startOpts,
			})
			if err != nil {
				log.Println("error running docker compose up: ", err)
				// return fmt.Errorf("unable running docker compose up: %w", err)
			}
		}()
	}

	w.runCommand(ctx, w.SetupCommands...)
	wg.Wait()

	return nil
}

func (w Workspace) Teardown(ctx context.Context, srv api.Service) error {
	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	wg := sync.WaitGroup{}

	if len(w.DockerCompose.Configs) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Setting down docker compose: ", w.Name)

			err := srv.Stop(ctx, w.DockerCompose.project.Name, api.StopOptions{})
			if err != nil {
				log.Println("error running docker compose down: ", err)
				// return fmt.Errorf("unable running docker compose down: %w", err)
			}
		}()
	}

	w.runCommand(ctx, w.TeardownCommands...)
	wg.Wait()

	return nil
}

func NewWorkspace(ctx context.Context, c ConfigWorkspace) (Workspace, error) {
	// TODO Validate
	if c.Name == "" {
		return Workspace{}, errors.New("name is required")
	}

	w := Workspace{
		Name:             c.Name,
		Directory:        c.Directory,
		SetupCommands:    c.SetupCommands,
		TeardownCommands: c.TeardownCommands,
		Timeout:          c.Timeout,
		DockerCompose: struct {
			Configs []string
			project *types.Project
		}{Configs: c.DockerCompose.Configs},
	}

	if w.Timeout == 0 {
		w.Timeout = 5
	}

	if len(w.DockerCompose.Configs) != 0 {
		p, err := createDockerComposeProject(ctx, w.Directory, w.DockerCompose.Configs)
		if err != nil {
			return Workspace{}, fmt.Errorf("error init docker project: %w", err)
		}

		w.DockerCompose.project = p
	}

	return w, nil
}

// runCommand runs multiple commands concurrently
// Context is used to cancel the commands in case of timeout
func (w Workspace) runCommand(ctx context.Context, cmd ...string) {
	c := make(chan struct{})
	wg := sync.WaitGroup{}

	for i, cmd := range cmd {
		log.Println("Running command : ", cmd)
		wg.Add(1)

		col := color.ColorByIndex(i)

		// TODO find a way to make it on blocking
		// TODO ensure ni memory leak

		go func(i int, cm string) {
			defer wg.Done()

			err := func() error {
				log.Printf("%d: Running command : %s", i, cm)
				cmd := exec.Command(os.Getenv("SHELL"), "-c", cm)

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					return fmt.Errorf("unable to get stdout: %w", err)
				}
				stderr, err := cmd.StderrPipe()
				if err != nil {
					return fmt.Errorf("unable to get stdout: %w", err)
				}

				cmd.Dir = w.Directory
				scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
				err = cmd.Start()
				if err != nil {
					return fmt.Errorf("unable to run command: %w", err)
				}

				for scanner.Scan() {
					log.Printf(color.Colorize(col, fmt.Sprintf("%d: > %s", i, scanner.Text())))
				}

				if scanner.Err() != nil {
					err := cmd.Process.Kill()
					if err != nil {
						return fmt.Errorf("unable to kill process: %w", err)
					}
					err = cmd.Wait()
					if err != nil {
						return fmt.Errorf("unable to wait command (%s): %w", cmd, err)
					}
					return fmt.Errorf("unable to read stdout: %w", err)
				}

				err = cmd.Wait()
				if err != nil {
					return fmt.Errorf("error while running command (%s): %w", cmd, err)
				}

				log.Printf(color.Colorize(col, fmt.Sprintf("%d: Command finished", i)))

				return nil
			}()
			if err != nil {
				log.Printf("error running command : %v", err)
			}
		}(i, cmd)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		log.Println("commands finished")
	case <-ctx.Done():
		log.Println("timeout reached" + ctx.Err().Error())
	}
}
