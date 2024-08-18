package zombie_test

import (
	"testing"

	"github.com/renevo/zombieutils/pkg/zombie"
)

func TestAPI(t *testing.T) {
	api := zombie.NewAPI()
	stats, err := api.ServerStats()
	if err != nil {
		t.Fatalf("failed to get server stats: %v", err)
	}

	t.Logf("Server stats: %+v", stats)
}
