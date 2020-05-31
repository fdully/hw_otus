package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("normal run, with zero exit code", func(t *testing.T) {
		exitCode := RunCmd([]string{"echo"}, nil)
		require.Equal(t, exitCode, 0)
	})

	t.Run("run without command", func(t *testing.T) {
		exitCode := RunCmd([]string{}, nil)
		require.Equal(t, exitCode, -1)
	})

	t.Run("check command environment", func(t *testing.T) {

		var customEnvironment = Environment{"BAR": "bar", "FOO": "   foo\nwith new line", "HELLO": "\"hello\"", "UNSET": ""}

		if err := os.Setenv("UNSET", "Four"); err != nil {
			t.Fatalf("can't set env var %s\n", err.Error())
		}

		commandEnvironment := MakeCommandEnv(os.Environ(), customEnvironment)

		for _, v := range commandEnvironment {
			if strings.Contains(v, "UNSET=") {
				t.Fatalf("UNSET must be removed. Got: %s", v)
			}
		}

	})
}
