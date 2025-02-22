package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, val := range env {
		if !val.NeedRemove {
			err := os.Setenv(k, val.Value)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		err := os.Unsetenv(k)
		if err != nil {
			log.Fatal(err)
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return command.ProcessState.ExitCode()
}
