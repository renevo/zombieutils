package steam

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/renevo/zombieutils/pkg/logutil"
)

type Game struct {
	ID   int
	Beta string
}

func (g Game) Install(steamcmd, installPath string) error {
	var args []string

	if len(g.Beta) > 0 {
		args = []string{
			"+force_install_dir", installPath,
			"+login", "anonymous",
			"+app_update", fmt.Sprintf("%d", g.ID),
			"-beta", g.Beta,
			"validate",
			"+quit",
		}
	} else {
		args = []string{
			"+force_install_dir", installPath,
			"+login", "anonymous",
			"+app_update", fmt.Sprintf("%d", g.ID),
			"validate",
			"+quit",
		}
	}

	cmd := exec.Command(steamcmd, args...)
	cmd.Dir = installPath

	cmd.Stdin = os.Stdin
	cmd.Stdout = logutil.Writer{}
	cmd.Stderr = logutil.Writer{IsErr: true}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install steam game %d: %w", g.ID, err)
	}

	return nil
}
