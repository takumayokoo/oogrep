package main

import (
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

func isTerminal() bool {
	return terminal.IsTerminal(int(syscall.Stdout))
}

func supportTerminalColor() bool {
	return false
}
