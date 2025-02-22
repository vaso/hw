package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestReadDir(t *testing.T) {
	t.Run("test example env", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)

		expectedEnv := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: false,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}
		require.Equal(t, expectedEnv, env)
	})
}
