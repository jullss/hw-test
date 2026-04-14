package main

import "os"

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	dir := os.Args[1]

	env, err := ReadDir(dir)
	if err != nil {
		os.Exit(1)
	}

	commandArgs := os.Args[2:]

	code := RunCmd(commandArgs, env)
	os.Exit(code)
}
