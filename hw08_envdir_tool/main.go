package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage /path command arg1 arg2 ...")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
