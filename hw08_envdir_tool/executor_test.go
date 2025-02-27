package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestRunCmd(t *testing.T) {
	t.Run("test example env", func(t *testing.T) {
		args := []string{"echo", "arg1=1", "arg2=2"}
		env := Environment{}
		rsc := RunCmd(args, env)
		require.Equal(t, 0, rsc)
	})
}
