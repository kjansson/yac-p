package controller

import (
	"testing"

	"github.com/kjansson/yac-p/internal/test_utils"
	"github.com/kjansson/yac-p/pkg/types"
)

func TestConfigLoad(t *testing.T) {

	_, err := NewController(types.Config{
		ConfigFileLoader: test_utils.GetTestConfigLoader(),
		RemoteWriteURL:   "http://localhost:9090/api/v1/write",
	})
	if err != nil {
		t.Fatalf("Failed to create controller: %v", err)
	}
}
