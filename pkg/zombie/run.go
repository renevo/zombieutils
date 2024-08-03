package zombie

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/renevo/zombieutils/pkg/logutil"
)

func (s *Server) Run(ctx context.Context, gmsgs chan<- string) error {
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
	cmd.SysProcAttr = getSysProcAttr()

	stdin := bytes.Buffer{}

	cmd.Stdin = &stdin
	cmd.Stdout = logutil.Writer{GlobalMessagesChan: gmsgs}
	cmd.Stderr = logutil.Writer{GlobalMessagesChan: gmsgs, IsErr: true}

	go func() {
		<-ctx.Done()
		_ = cmd.Process.Signal(syscall.SIGINT)
	}()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running server: %s: %w", cmd.Path, err)
	}

	return nil
}
