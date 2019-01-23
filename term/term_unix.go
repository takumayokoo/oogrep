// +build !windows

package term

import (
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func IsTerminal() bool {
	return terminal.IsTerminal(syscall.Stdout)
}

func SupportTerminalColor() bool {
	return true
}
