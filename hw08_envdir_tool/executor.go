package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	cmdName := filepath.Clean(cmd[0])
	command := exec.Command(cmdName, cmd[1:]...)

	envVars := make([]string, 0, len(env))

	for _, item := range os.Environ() {
		key := strings.SplitN(item, "=", 2)[0]
		if v, ok := env[key]; ok && v.NeedRemove {
			continue
		}
		envVars = append(envVars, item)
	}

	for k, v := range env {
		if !v.NeedRemove {
			envVars = append(envVars, k+"="+v.Value)
		}
	}

	command.Env = envVars
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
