//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func Tests() error {
	fmt.Println("running tests...")
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
