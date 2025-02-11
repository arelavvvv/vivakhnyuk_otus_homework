package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		require.Nil(t, err)

		cmd := []string{"go", "version"}

		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)
	})
}