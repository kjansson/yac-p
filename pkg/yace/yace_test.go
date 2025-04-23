package yace

import (
	"testing"

	"github.com/kjansson/yac-p/internal/test_utils"
	"github.com/kjansson/yac-p/pkg/types"
)

func TestConfigLoad(t *testing.T) {
	c := &YaceClient{}

	conf := types.Config{
		ConfigFileLoader: test_utils.GetTestConfigLoader(),
	}

	err := c.Init(conf) // Initialize all components
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
}
