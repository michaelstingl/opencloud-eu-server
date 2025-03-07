package command

import (
	"context"
	"os"

	"github.com/opencloud-eu/opencloud/pkg/clihelper"
	occfg "github.com/opencloud-eu/opencloud/pkg/config"
	"github.com/opencloud-eu/opencloud/services/auth-bearer/pkg/config"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		// start this service
		Server(cfg),

		// interaction with this service

		// infos about this service
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the opencloud auth-bearer command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "auth-bearer",
		Usage:    "Provide bearer authentication for OpenCloud",
		Commands: GetCommands(cfg),
	})

	return app.RunContext(cfg.Context, os.Args)
}

// SutureService allows for the accounts command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new auth-bearer.SutureService
func NewSutureService(cfg *occfg.Config) suture.Service {
	cfg.AuthBearer.Commons = cfg.Commons
	return SutureService{
		cfg: cfg.AuthBearer,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
