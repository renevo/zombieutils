package zombie

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
	"github.com/renevo/zombieutils/pkg/zombie"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	configFile := "./zombie.hcl"

	srv := zombie.Default()
	serverConfig := struct {
		Server []*zombie.Server `hcl:"server,block"`
	}{
		Server: []*zombie.Server{srv},
	}

	loadConfig := func() error {
		if len(configFile) > 0 {
			if err := hclsimple.DecodeFile(configFile, nil, &serverConfig); err != nil {
				return errors.Wrap(err, "failed to parse config file")
			}

			if len(serverConfig.Server) != 1 {
				return errors.New("you must specify exactly one server block in the configuration file")
			}
		}

		return nil
	}

	zombieCommand := &cobra.Command{
		Use:   "zombie",
		Short: "7 Days to Die Server",
		Long:  "Commands to install and run a 7 days to die server",
	}

	zombieCommand.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validates configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			logrus.Info("Validated server configuration")
			return nil
		},
	})

	zombieCommand.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install 7 days to die",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			return srv.Install(context.Background())
		},
	})

	zombieCommand.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run a 7 days to die server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := srv.Install(ctx); err != nil {
				return err
			}

			sigCh := make(chan os.Signal, 2)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			go func() {
				sig := <-sigCh
				logrus.Infof("Stopping server... %v", sig)
				cancel()
			}()

			err := srv.Run(ctx)

			if err != nil {
				logrus.Infof("Stopped Server: %v", err)
			} else {
				logrus.Infof("Stopped Server")
			}

			return err
		},
	})

	zombieCommand.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "specify an optional configuration file")

	return zombieCommand
}
