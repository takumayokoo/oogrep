package main

import (
	"fmt"
	"github.com/takumayokoo/oogrep/term"
	"strings"
)

func formatMatchText(matchstring, pattern string) string {
	if term.IsTerminal() && term.SupportTerminalColor() {
		return strings.Replace(matchstring, string(pattern), "\x1b[31m"+string(pattern)+"\x1b[0m", -1)
	} else {
		return matchstring
	}
}

func formatFileName(s string) string {
	if term.IsTerminal() && term.SupportTerminalColor() {
		return "\x1b[36m" + s + "\x1b[0m"
	} else {
		return s
	}
}

func PrintFileNameLine(filename string) {
	fmt.Println(formatFileName(filename))
}

func PrintMatchLine(filename string, matchedStrings []string, pattern string) {
	if term.IsTerminal() {
		for _, s := range matchedStrings {
			s1 := formatMatchText(s, string(pattern))
			fmt.Printf("\t%v\n", s1)
		}
		fmt.Println()
	} else {
		for _, s := range matchedStrings {
			fmt.Printf("%v: %v\n", filename, s)
		}
	}
}
