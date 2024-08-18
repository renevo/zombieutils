package discord_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/renevo/zombieutils/internal/discord"
	"github.com/renevo/zombieutils/pkg/zombie"
	"github.com/subosito/gotenv"
)

func TestStatus(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	_ = gotenv.Load("../../.env")
	disco, err := discord.New()

	if err != nil {
		t.Skipf("Failed to create discord session: %v", err)
		return
	}

	api := zombie.NewAPI()

	for i := 0; i < 10; i++ {
		stats, err := api.ServerStats()
		if err != nil {
			t.Errorf("Failed to get server stats: %v", err)
		}

		t.Logf("Server stats: %s", stats)

		if err := disco.UpdateStatus(stats); err != nil {
			t.Errorf("Failed to update discord status: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	_ = disco.Close()
}
