package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/config"
)

func initializeTestCache(t *testing.T) {
	t.Helper()
	err := Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
}
