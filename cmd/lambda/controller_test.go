package main

import (
	"testing"

	"github.com/kjansson/yac-p/v3/internal/test_utils"
)

func TestConfigLoad(t *testing.T) {

	_, err := NewController(Config{
		ConfigFileLoader: test_utils.GetTestConfigLoader(),
		RemoteWriteURL:   "http://localhost:9090/api/v1/write",
	})

	if err != nil {
		t.Fatalf("Failed to create controller: %v", err)
	}
}
