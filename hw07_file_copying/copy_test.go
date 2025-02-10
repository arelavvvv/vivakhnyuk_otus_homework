package main

import (
	"io"
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("offset 0, limit 0", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 0))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset0_limit0.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 10", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 10))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset0_limit10.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 1000))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset0_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 0, limit 10000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 0, 10000))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset0_limit10000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 100, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 100, 1000))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset100_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})

	t.Run("offset 6000, limit 1000", func(t *testing.T) {
		require.Equal(t, nil, Copy("testdata/input.txt", "out.txt", 6000, 1000))
		require.Equal(t, true, FEqual("out.txt", "testdata/out_offset6000_limit1000.txt"))
		require.Equal(t, nil, os.Remove("out.txt"))
	})
}

func FEqual(file1Path, file2Path string) bool {
	file1, err := os.Open(file1Path)
	if err != nil {
		return false
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		return false
	}
	defer file2.Close()

	buffer1 := make([]byte, 32*1024)
	buffer2 := make([]byte, 32*1024)

	for {
		n1, err1 := file1.Read(buffer1)
		if err1 != nil && err1 != io.EOF {
			return false
		}

		n2, err2 := file2.Read(buffer2)
		if err2 != nil && err2 != io.EOF {
			return false
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true
		}

		if err1 == io.EOF || err2 == io.EOF {
			return false
		}

		if n1 != n2 {
			return false
		}

		for i := 0; i < n1; i++ {
			if buffer1[i] != buffer2[i] {
				return false
			}
		}
	}
}
