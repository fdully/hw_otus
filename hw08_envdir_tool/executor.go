package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const IncorrectCommand = -1

// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd, env []string) (returnCode int) {
	// Place your code here
	if len(cmd) == 0 {
		return IncorrectCommand
	}
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Env = env

	// run command
	if err := c.Run(); err != nil {
		fmt.Printf("execute command error: %s\n", err.Error())
	}
	return c.ProcessState.ExitCode()
}

func MakeCommandEnv(osEnvironment []string, customEnvironment map[string]string) []string {
	var e = osEnvironment

	for key, val := range customEnvironment {
		// find empty var and remove from env
		if val == "" {
			for i, v := range e {
				k := strings.Split(v, "=")[0]
				if key == k {
					e = append(e[:i], e[i+1:]...)
				}
			}
			continue
		}
		e = append(e, fmt.Sprintf("%s=%s", key, val))
	}

	return e
}
