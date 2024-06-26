package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

// cmdcompose "github.com/docker/compose/v2/cmd/compose"

func createDockerComposeProject(ctx context.Context, dir string, config []string) (*types.Project, error) {
	cfgs := make([]types.ConfigFile, 0)

	for _, c := range config {
		cfgs = append(cfgs, types.ConfigFile{
			// TODO path
			Filename: dir + "/" + c,
		})
	}

	configDetails := types.ConfigDetails{
		WorkingDir:  dir,
		ConfigFiles: cfgs,
		Environment: nil,
	}

	// TODO project name cleaner
	split := strings.Split(dir, "/")
	projectName := split[len(split)-1]

	p, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(projectName, true)
	})
	if err != nil {
		return p, fmt.Errorf("error load project '%s': %w", projectName, err)
	}

	// TODO use included service label injection
	addServiceLabels(p)

	return p, nil
}

// createDockerService creates a docker service which can be
// used to interact with docker-compose.
func createDockerService() (api.Service, error) {
	var srv api.Service
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return srv, fmt.Errorf("error creating docker cli: %w", err)
	}

	err = dockerCli.Initialize(&flags.ClientOptions{Context: getDockerContext(), LogLevel: "error"})
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}

/*
addServiceLabels adds the labels docker compose expects to exist on services.
This is required for future compose operations to work, such as finding
containers that are part of a service.
*/
func addServiceLabels(project *types.Project) {
	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False", // default, will be overridden by `run` command
		}
		project.Services[i] = s
	}
}

func getDockerContext() string {
	env := os.Getenv("DOCKER_CONTEXT")

	if env == "" {
		env = "default"
	}
	// TODO Get current context from docker

	return env
}
