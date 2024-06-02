package core

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

func createDockerProject(ctx context.Context, dir string, config []string) (*types.Project, error) {
	cfgs := make([]types.ConfigFile, 0)

	for _, c := range config {
		cfgs = append(cfgs, types.ConfigFile{
			// TODO path
			Filename: dir + "/" + c,
		})
	}

	configDetails := types.ConfigDetails{
		WorkingDir: dir,
		// ConfigFiles: types.ToConfigFiles(config),
		ConfigFiles: cfgs,
		Environment: nil,
	}

	// TODO add project name
	projectName := "api-go"

	p, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(projectName, true)
	})
	if err != nil {
		return p, fmt.Errorf("error load: %w", err)
	}
	addServiceLabels(p)
	return p, nil
}

// createDockerService creates a docker service which can be
// used to interact with docker-compose.
func createDockerService() (api.Service, error) {
	var srv api.Service
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return srv, err
	}

	dockerContext := "default"

	// Magic line to fix error:
	// Failed to initialize: unable to resolve docker endpoint: no context store initialized
	myOpts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(myOpts)
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}

func myExec(ctx context.Context, srv api.Service, p *types.Project, cmd string) {
	result, err := srv.Exec(ctx, p.Name, api.RunOptions{
		Service:     "testservice",
		Command:     []string{"/bin/bash", "-c", cmd},
		WorkingDir:  "/bin",
		Tty:         true,
		Environment: []string{"TONE=test1"},
	})
	log.Println("Command result:", result, " and err:", err)
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
