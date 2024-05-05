package app

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
)

// Config represents the configuration of tix
type Config struct {
	Workspaces []Workspace `toml:"workspace"`
}

// Workspace represents a workspace
type Workspace struct {
	Name      string
	Directory string
	// TODO Define default values
	Timeout          int
	SetupCommands    []string
	TeardownCommands []string
}

func (w Workspace) Setup() {
	log.Println("Setting up workspace : ", w.Name)
	// TODO configure timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	w.runCommand(ctx, w.SetupCommands...)
}

func (w Workspace) Teardown() {
	log.Println("Tearing down workspace : ", w.Name)
	// TODO configure timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	w.runCommand(ctx, w.TeardownCommands...)
}

// TODO Implement timeout
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
			// TODO redirect to stdout
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
