package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	testFiles := []struct {
		fileName    string
		content     string
		expected    string
		isRemoved   bool
		shouldExist bool
	}{
		{"VALID_VAR", "hello", "hello", false, true},
		{"WITH_SPACES", "world  \t", "world", false, true},
		{"EMPTY_FILE", "", "", true, true},
		{"BINARY_DATA", "foo\x00not_ignore", "foo\nnot_ignore", false, true},
		{"MULTI_LINE", "line1\nline2", "line1", false, true},
		{"INVALID=VAR", "some_data", "", false, false},
	}

	for _, tf := range testFiles {
		err := os.WriteFile(filepath.Join(tmpDir, tf.fileName), []byte(tf.content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	env, err := ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	for _, tf := range testFiles {
		t.Run(tf.fileName, func(t *testing.T) {
			val, exists := env[tf.fileName]

			if exists != tf.shouldExist {
				t.Errorf("exists = %v, want %v", exists, tf.shouldExist)
				return
			}

			if tf.shouldExist {
				if val.Value != tf.expected {
					t.Errorf("Value = %q, want %q", val.Value, tf.expected)
				}

				if val.NeedRemove != tf.isRemoved {
					t.Errorf("NeedRemove = %v, want %v", val.NeedRemove, tf.isRemoved)
				}
			}
		})
	}
}
