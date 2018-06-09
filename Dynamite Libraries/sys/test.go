package main

import (
	"os"
	"os/exec"
)

func System(cmd_string, args string) {
	// s := strings.Split(args, " ")
	// cmd := exec.Command(cmd_string, s...)
	cmd := exec.Command(cmd_string, args)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	// System("systeminfo > test.txt")
}
