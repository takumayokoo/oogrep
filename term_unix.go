// +build !windows

package main

import (
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

func isTerminal() bool {
	return terminal.IsTerminal(syscall.Stdout)
}

func supportTerminalColor() bool {
	return true
}
