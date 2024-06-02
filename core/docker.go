package core

import (
	"context"
	"fmt"
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
		return p, fmt.Errorf("error load: %w", err)
	}
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

	// TODO Configurable ?
	dockerContext := "default"

	myOpts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(myOpts)
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}
