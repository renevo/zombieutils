package zombie

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/zombieutils/pkg/logutil"
	"github.com/sirupsen/logrus"
)

func (s *Server) Run(ctx context.Context) error {
	installDirectory, _ := filepath.Abs(filepath.FromSlash(s.Path))
	configFile, _ := filepath.Abs(filepath.FromSlash(s.Config))
	args := []string{
		"-logfile", "/dev/stdout",
		"-quit",
		"-batchmode",
		"-nographics",
		"-dedicated",
		fmt.Sprintf("-configfile=%s", configFile),
	}

	cmd := exec.Command(filepath.Join(installDirectory, "7DaysToDieServer.x86_64"), args...)
	cmd.Dir = installDirectory
	cmd.Env = append(cmd.Env, fmt.Sprintf("LD_LIBRARY_PATH=.:%s/7DaysToDieServer_Data/Plugins/x86_64", installDirectory))
	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s/GameData", installDirectory))

	cmd.Stdin = os.Stdin
	cmd.Stdout = logutil.Writer(logrus.Info)
	cmd.Stderr = logutil.Writer(logrus.Error)

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "faild to start server %q", cmd.Path)
	}

	<-ctx.Done()

	_ = cmd.Process.Signal(os.Interrupt)

	return cmd.Wait()
}
