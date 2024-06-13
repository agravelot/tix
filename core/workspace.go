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
	"github.com/docker/compose/v2/pkg/api"
	"github.com/kr/pretty"
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
	DockerCompose    struct{ Configs []string }
	Timeout          int
	Applets          []Applet
}

func (w Workspace) IsRunning() bool {
	for _, a := range w.Applets {
		if a.IsRunning {
			return true
		}
	}

	return false
}

func (w *Workspace) RefreshApplets(ctx context.Context, srv api.Service) error {
	if len(w.DockerCompose.Configs) == 0 {
		return nil
	}

	p, err := createDockerComposeProject(ctx, w.Directory, w.DockerCompose.Configs)
	if err != nil {
		return fmt.Errorf("error create docker project: %w", err)
	}

	sum, err := srv.Ps(ctx, p.Name, api.PsOptions{All: true})
	if err != nil {
		return err
	}

	var applets []Applet

	for _, s := range sum {
		applets = append(applets, Applet{
			Icon:      "docker",
			IsRunning: s.State == "running",
		})
	}

	w.Applets = applets

	pretty.Println(w.Applets)

	return nil
}

func (w Workspace) Setup(srv api.Service) error {
	log.Println("Setting up workspace : ", w.Name)
	ctx := context.Background()

	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	// TODO Refactor outside of workspace
	if len(w.DockerCompose.Configs) > 0 {
		log.Println("Setting up docker compose: ", w.Name)

		// TODO one per workspace
		p, err := createDockerComposeProject(ctx, w.Directory, w.DockerCompose.Configs)
		if err != nil {
			return fmt.Errorf("error create docker project: %w", err)
		}

		createOpts := api.CreateOptions{
			Recreate:             api.RecreateDiverged,
			RecreateDependencies: api.RecreateDiverged,
			Inherit:              true,
			Services:             p.ServiceNames(),
			Build: &api.BuildOptions{
				Services: p.ServiceNames(),
			},
		}

		startOpts := api.StartOptions{
			Project: p,
		}

		err = srv.Up(ctx, p, api.UpOptions{
			Create: createOpts,
			Start:  startOpts,
		})
		if err != nil {
			return fmt.Errorf("unable running docker compose up: %w", err)
		}
	}

	w.runCommand(ctx, w.SetupCommands...)

	return nil
}

func (w Workspace) Teardown(srv api.Service) error {
	log.Println("Tearing down workspace : ", w.Name)
	ctx := context.Background()

	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	if len(w.DockerCompose.Configs) > 0 {
		log.Println("Setting down docker compose: ", w.Name)

		// TODO one per workspace
		p, err := createDockerComposeProject(ctx, w.Directory, w.DockerCompose.Configs)
		if err != nil {
			return fmt.Errorf("error down docker project: %w", err)
		}

		// err = srv.Down(ctx, p.Name, api.DownOptions{})
		err = srv.Stop(ctx, p.Name, api.StopOptions{})
		if err != nil {
			return fmt.Errorf("unable running docker compose down: %w", err)
		}
	}

	w.runCommand(ctx, w.TeardownCommands...)

	return nil
}

type WorkspaceBuilder struct {
	Name      string
	Directory string
}

func NewWorkspace(c ConfigWorkspace) (Workspace, error) {
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

// runCommand runs multiple commands concurrently
// Context is used to cancel the commands in case of timeout
func (w Workspace) runCommand(ctx context.Context, cmd ...string) {
	c := make(chan struct{})
	wg := sync.WaitGroup{}

	for i, cmd := range cmd {
		log.Println("Running command : ", cmd)
		wg.Add(1)

		col := color.ColorByIndex(i)

		go func(i int, cm string) {
			log.Printf("%d: Running command : %s", i, cm)
			cmd := exec.Command(os.Getenv("SHELL"), "-c", cm)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatal(fmt.Errorf("unable to get stdout: %w", err))
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Fatal(fmt.Errorf("unable to get stdout: %w", err))
			}

			cmd.Dir = w.Directory
			scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
			err = cmd.Start()
			if err != nil {
				log.Fatal(fmt.Errorf("unable to run command: %w", err))
			}

			for scanner.Scan() {
				log.Printf(color.Colorize(col, fmt.Sprintf("%d: > %s", i, scanner.Text())))
			}

			if scanner.Err() != nil {
				err := cmd.Process.Kill()
				if err != nil {
					log.Fatal(fmt.Errorf("unable to kill process: %w", err))
				}
				err = cmd.Wait()
				if err != nil {
					log.Fatal(fmt.Errorf("unable to wait command (%s): %w", cmd, err))
				}
				log.Fatal(fmt.Errorf("unable to read stdout: %w", err))
			}

			err = cmd.Wait()
			if err != nil {
				log.Fatal(fmt.Errorf("unable to wait command (%s): %w", cmd, err))
			}

			log.Printf(color.Colorize(col, fmt.Sprintf("%d: Command finished", i)))

			wg.Done()
		}(i, cmd)
	}

	go func() {
		wg.Wait()
		defer close(c)
	}()

	select {
	case <-c:
		log.Println("Commands finished")
	case <-ctx.Done():
		log.Println("Timeout reached" + ctx.Err().Error())
	}
}
