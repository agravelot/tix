package core

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/agravelot/tix/color"
	"github.com/docker/compose/v2/pkg/api"
)

// Workspace represents a workspace
type Workspace struct {
	Name             string
	Directory        string
	Shell            string
	SetupCommands    []string
	TeardownCommands []string
	DockerCompose    struct{ Configs []string }
	Timeout          int
}

func (w Workspace) Setup() error {
	log.Println("Setting up workspace : ", w.Name)
	ctx := context.Background()

	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	log.Println("Setting up comp : ", w.DockerCompose)

	if len(w.DockerCompose.Configs) > 0 {
		log.Println("Setting up docker compose: ", w.Name)

		p, err := createDockerProject(ctx, w.Directory, w.DockerCompose.Configs)
		if err != nil {
			return fmt.Errorf("error create docker project: %w", err)
		}

		srv, err := createDockerService()
		if err != nil {
			return fmt.Errorf("error create docker service: %w", err)
		}

		fmt.Println("Docker service up...")
		err = srv.Up(ctx, p, api.UpOptions{})
		if err != nil {
			return fmt.Errorf("unable running docker compose up: %w", err)
		}
	}

	w.runCommand(ctx, w.SetupCommands...)

	return nil
}

func (w Workspace) Teardown() {
	log.Println("Tearing down workspace : ", w.Name)
	ctx := context.Background()

	if w.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(w.Timeout)*time.Second)
		defer cancel()
	}

	w.runCommand(ctx, w.TeardownCommands...)
}

type WorkspaceBuilder struct {
	Name      string
	Directory string
}

// func NewWorkspace(name string, params ) Workspace {
// 	return Workspace{
// 		Name: name,
// 	}
// }

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
