package yace

import (
	"testing"

	"github.com/kjansson/yac-p/pkg/tests"
)

func TestConfigLoad(t *testing.T) {
	c := &YaceClient{}

	err := c.Init(tests.GetTestConfigLoader()) // Initialize all components
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
}
