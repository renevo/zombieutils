package zombie

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

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

	stdin := bytes.Buffer{}

	cmd.Stdin = &stdin
	cmd.Stdout = logutil.Writer(logrus.Info)
	cmd.Stderr = logutil.Writer(logrus.Error)

	go func() {
		<-ctx.Done()
		_ = cmd.Process.Signal(syscall.SIGINT)
	}()

	return errors.Wrapf(cmd.Run(), "failed running server: %s", cmd.Path)
}
