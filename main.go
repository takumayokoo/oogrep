package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	recursive bool
	fileoonly bool
)

var docxRe = regexp.MustCompile(`(?U)<w:t>(.*)<\/w:t>`)
var xlsxRe = regexp.MustCompile(`(?U)<t>(.*)<\/t>`)
var pptxRe = regexp.MustCompile(`(?U)<a:t>(.*)<\/a:t>`)

var ooFilePatterns = []OoFile{
	OoFile{Ext: ".docx", TextFilePattern: `word/document\.xml`, Re: docxRe},
	OoFile{Ext: ".docm", TextFilePattern: `word/document\.xml`, Re: docxRe},
	OoFile{Ext: ".xlsx", TextFilePattern: `xl/sharedStrings\.xml`, Re: xlsxRe},
	OoFile{Ext: ".xlsm", TextFilePattern: `xl/sharedStrings\.xml`, Re: xlsxRe},
	OoFile{Ext: ".pptx", TextFilePattern: `ppt/slides/slide[0-9]+\.xml`, Re: pptxRe},
}

var stderrLog = log.New(os.Stderr, "", 0)

func targetFiles(paths []string) []string {
	var filepaths []string

	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		if fi.Mode().IsDir() {
			dir := path

			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					filepaths = append(filepaths, path)
				}

				return nil
			})

			if err != nil {
				log.Fatal(err)
			}
		} else {
			filepaths = append(filepaths, path)
		}
	}

	return filepaths
}

func isTerminal() bool {
	return terminal.IsTerminal(int(syscall.Stdout))
}

func supportTerminalColor() bool {
	return terminal.IsTerminal(int(syscall.Stdout)) && runtime.GOOS != "windows"
}

func formatFileName(s string) (string) {
	if supportTerminalColor() {
		return "\x1b[36m"+s+"\x1b[0m"
	} else {
		return s
	}
}

func formatMatchText(matchstring, pattern string) string {
	if supportTerminalColor() {
		return strings.Replace(matchstring, string(pattern), "\x1b[31m"+string(pattern)+"\x1b[0m", -1)
	} else {
		return matchstring
	}
}


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&fileoonly, "l", false, "print file name")
	flag.Parse()

	pattern := []byte(flag.Args()[0])

	var paths []string

	if flag.NArg() > 1 {
		paths = flag.Args()[1:]
	} else {
		paths = []string{"./"}
	}

	var filepaths = targetFiles(paths)

	for _, fp := range filepaths {
		for _, op := range ooFilePatterns {
			if op.Ext == filepath.Ext(fp) {
				matchedStrings, err := op.Match(fp, pattern, fileoonly)
				if err != nil {
					stderrLog.Println(err)
				}

				filename := formatFileName(fp)

				if matchedStrings != nil {
					if fileoonly {
						fmt.Println(filename)
					} else {
						if isTerminal() {
							fmt.Println(filename)
							for _, s := range matchedStrings {
								s1 := formatMatchText(s, string(pattern))
								fmt.Printf("\t%v\n", s1)
							}
							fmt.Println()
						} else {
							for _, s := range matchedStrings {
								fmt.Printf("%v: %v\n", fp, s)
							}
						}
					}
				}
			}
		}
	}
}
