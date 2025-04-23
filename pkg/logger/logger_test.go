package logger

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/kjansson/yac-p/pkg/types"
)

// {"time":"2025-04-23T16:53:55.002724+02:00","level":"INFO","msg":"test message","key1":"value1"}
type TestLogEntry struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Key1  string `json:"key1"`
}

func TestLogLevel(t *testing.T) {

	l := &SlogLogger{}

	tmpFile, err := os.CreateTemp(".", "logtest")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}()

	defer func() {
		err := tmpFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}
	}()

	err = l.Init(types.Config{
		Debug:          false,
		LogDestination: tmpFile,
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	l.Log("debug", "test message", "key1", "value1")

	entry := make([]byte, 256)
	n, err := tmpFile.ReadAt(entry, 0)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	if n > 0 {
		t.Fatalf("Expected no logs due to lower log level attempt, got '%s'", string(entry))
	}
}

func TestLogFormat(t *testing.T) {

	l := &SlogLogger{}

	tmpFile, err := os.CreateTemp(".", "logtest")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}()

	defer func() {
		err := tmpFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}
	}()

	err = l.Init(types.Config{
		Debug:          false,
		LogDestination: tmpFile,
		LogFormat:      "json",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	l.Log("info", "test message", "key1", "value1")

	entry := make([]byte, 256)
	n, err := tmpFile.ReadAt(entry, 0)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	// Check if the log entry is in JSON format
	entryJSON := TestLogEntry{}
	err = json.Unmarshal(entry[:n], &entryJSON)
	if err != nil {
		t.Fatalf("Expected JSON format, got '%s'", string(entry[:n]))
	}
}
