package yace

import (
	"testing"

	"github.com/kjansson/yac-p/v2/internal/test_utils"
)

func TestConfigLoad(t *testing.T) {
	_, err := NewYaceClient(
		test_utils.GetTestConfigLoader(),
		YaceOpts{},
	)

	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
}
