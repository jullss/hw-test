package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		str, needRemove, err := readLine(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", filePath, err)
			continue
		}

		env[file.Name()] = EnvValue{Value: str, NeedRemove: needRemove}
	}
	return env, nil
}

func readLine(filePath string) (string, bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false, err
	}

	if len(data) == 0 {
		return "", true, nil
	}

	line := string(data)

	if i := strings.Index(line, "\n"); i != -1 {
		line = line[:i]
	}

	line = strings.ReplaceAll(line, "\x00", "\n")
	line = strings.TrimRight(line, " \t")

	return line, false, nil
}
