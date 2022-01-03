package steam

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/renevo/zombieutils/pkg/logutil"
	"github.com/sirupsen/logrus"
)

type Game struct {
	ID   int
	Beta string
}

func (g Game) Install(steamcmd, installPath string) error {
	installDirectory, _ := filepath.Abs(filepath.FromSlash(installPath))

	if err := os.MkdirAll(installDirectory, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create install directory: %q", installDirectory)
	}

	var args []string

	if len(g.Beta) > 0 {
		args = []string{
			"+force_install_dir", installDirectory,
			"+login", "anonymous",
			"+app_update", fmt.Sprintf("%d", g.ID),
			"-beta", g.Beta,
			"validate",
			"+quit",
		}
	} else {
		args = []string{
			"+force_install_dir", installDirectory,
			"+login", "anonymous",
			"+app_update", fmt.Sprintf("%d", g.ID),
			"validate",
			"+quit",
		}
	}

	cmd := exec.Command(steamcmd, args...)
	cmd.Dir = installDirectory

	cmd.Stdin = os.Stdin
	cmd.Stdout = logutil.Writer(logrus.Info)
	cmd.Stderr = logutil.Writer(logrus.Error)

	return errors.Wrapf(cmd.Run(), "failed to install steam game %d", g.ID)
}
