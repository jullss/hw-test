package main

import (
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	fromPath := "testdata/from_path.txt"
	toPath := "testdata/to_path.txt"

	str := "hello world"
	content := []byte(str)

	err := os.WriteFile(fromPath, content, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(fromPath)
	defer os.Remove(toPath)

	t.Run("full copy", func(t *testing.T) {
		err := Copy(fromPath, toPath, 0, 0)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		result, _ := os.ReadFile(toPath)
		if string(result) != str {
			t.Errorf("Expected 'hello world', got '%s'", string(result))
		}
	})

	t.Run("offset and limit", func(t *testing.T) {
		os.Remove(toPath)

		err := Copy(fromPath, toPath, 6, 5)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		result, _ := os.ReadFile(toPath)
		if string(result) != "world" {
			t.Errorf("Expected 'world', got '%s'", string(result))
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		err := Copy(fromPath, toPath, 100, 0)

		if !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Errorf("Expected ErrOffsetExceedsFileSize, got %v", err)
		}
	})
}
