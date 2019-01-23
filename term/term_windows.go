package term

import (
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

var supportColor = false

func init() {
	if !IsTerminal() {
		return
	}

	h, err := GetStdHandle(syscall.Handle(uintptr(STD_OUTPUT_HANDLE)))
	if err != nil {
		panic(err)
	}

	mode, err := GetConsoleMode(h)
	if err != nil {
		panic(err)
	}

	mode = mode | ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = SetConsoleMode(h, mode)
	if err != nil {
		panic(err)
	}

	mode2, err := GetConsoleMode(h)
	if err != nil {
		panic(err)
	}

	supportColor = mode == mode2
}

func IsTerminal() bool {
	return terminal.IsTerminal(int(syscall.Stdout))
}

func SupportTerminalColor() bool {
	// Windows support console ANSI Escape Sequences only on Windows 10.
	return supportColor
}
