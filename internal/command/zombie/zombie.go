package zombie

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/renevo/zombieutils/internal/discord"
	"github.com/renevo/zombieutils/pkg/zombie"
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
				return fmt.Errorf("failed to parse config file: %w", err)
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

			slog.Info("Validated server configuration")
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

			gmsgs := make(chan string, 100)
			messenger, err := discord.New()
			if err != nil {
				return fmt.Errorf("failed to create discord messager: %w", err)
			}
			defer messenger.Close()

			api := zombie.NewAPI()
			go func(ctx context.Context) {
				t := time.NewTicker(time.Minute)
				defer t.Stop()

				for {
					select {
					case <-ctx.Done():
						return

					case <-t.C:
						stats, err := api.ServerStats()
						if err != nil {
							slog.Error("Failed to get server stats", "err", err)
							continue
						}

						if err := messenger.UpdateStatus(stats); err != nil {
							slog.Error("Failed to update discord status", "err", err)
						}
					}
				}
			}(ctx)

			messagerWG := sync.WaitGroup{}
			messagerWG.Add(1)
			go func() {
				messenger.Publish(gmsgs)
				messagerWG.Done()
			}()

			sigCh := make(chan os.Signal, 2)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			go func() {
				sig := <-sigCh
				slog.Info("Stopping server", "signal", sig)
				cancel()
			}()

			if err = srv.Run(ctx, gmsgs); err != nil {
				slog.Info("Stopped Server", "err", err)
			} else {
				slog.Info("Stopped Server")
			}

			// close and drain messages to discord
			close(gmsgs)
			messagerWG.Wait()

			return err
		},
	})

	zombieCommand.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "specify an optional configuration file")

	return zombieCommand
}
