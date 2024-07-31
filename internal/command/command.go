package command

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/renevo/zombieutils/internal/command/zombie"
	"github.com/spf13/cobra"
)

func Execute(args []string) error {
	verboseLogging := false
	nocolorLogging := false
	jsonLogging := false

	rootCommand := &cobra.Command{
		Use:   "zombieutils",
		Short: "7 Days To Die Utilities",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// logger setup
			var logLeveler slog.LevelVar
			var logHandler slog.Handler
			logOutput := os.Stdout

			switch {
			case jsonLogging:
				logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &logLeveler})
			case nocolorLogging:
				logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &logLeveler})
			default:
				logHandler = tint.NewHandler(colorable.NewColorable(logOutput), &tint.Options{
					Level:   &logLeveler,
					NoColor: !isatty.IsTerminal(logOutput.Fd()),
				})
			}

			if verboseLogging {
				logLeveler.Set(slog.LevelDebug)
			}

			slog.SetDefault(slog.New(logHandler))

			return nil
		},
	}

	rootCommand.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "verbose output")
	rootCommand.PersistentFlags().BoolVarP(&jsonLogging, "json", "j", false, "output logging as json")
	rootCommand.PersistentFlags().BoolVar(&nocolorLogging, "no-color", false, "disable colorized output")

	// add commands here:
	rootCommand.AddCommand(
		zombie.New(),
	)

	// execute
	rootCommand.SetArgs(args)
	return rootCommand.Execute()
}
