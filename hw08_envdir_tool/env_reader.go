package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

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

		name := file.Name()
		filePath := filepath.Join(dir, name)

		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		line, err := reader.ReadString('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		if fileInfo.Size() == 0 {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		lineBytes := []byte(line)
		lineBytes = bytes.ReplaceAll(lineBytes, []byte{0x00}, []byte{'\n'})
		line = string(lineBytes)

		line = strings.TrimRight(line, " \t\r\n")

		env[name] = EnvValue{Value: line, NeedRemove: false}
	}

	return env, nil
}
