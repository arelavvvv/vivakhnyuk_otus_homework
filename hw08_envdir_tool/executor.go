package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	command.Env = os.Environ()

	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return
			}
		}
		command.Env = append(command.Env, k+"="+v.Value)
	}

	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		log.Println("Error running command:", err)
		return 1
	}

	return 0
}
