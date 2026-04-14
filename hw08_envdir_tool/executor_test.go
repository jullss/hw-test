package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		env := Environment{}
		cmd := []string{}

		code := RunCmd(cmd, env)
		if code != 1 {
			t.Errorf("Expected 1, got %d", code)
		}
	})

	t.Run("successful execution", func(t *testing.T) {
		os.Setenv("TEST_UNSET", "should_be_remove")
		defer os.Unsetenv("TEST_UNSET")

		env := Environment{
			"TEST_NEW":   {Value: "new_value", NeedRemove: false},
			"TEST_UNSET": {Value: "", NeedRemove: true},
		}

		cmd := []string{"env"}

		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("Expected 0, got %d", code)
		}
	})
}
