package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func targetFiles() []string {
	var filepaths []string

	if recursive {
		var dirs []string
		relative := false
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if flag.NArg() > 1 {
			dirs = flag.Args()[1:]
		} else {
			relative = true
			dirs = append(dirs, currentDir)
		}

		for _, d := range dirs {

			err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					if relative {
						p := strings.TrimPrefix(path, currentDir)
						p = p[1:]
						filepaths = append(filepaths, p)
					} else {
						filepaths = append(filepaths, path)
					}
				}

				return nil
			})

			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		filepaths = os.Args[1:]
	}

	return filepaths
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&recursive, "r", false, "recursive")
	flag.BoolVar(&fileoonly, "l", false, "print file name")
	flag.Parse()

	pattern := []byte(flag.Args()[0])

	var filepaths = targetFiles()

	for _, fp := range filepaths {
		for _, op := range ooFilePatterns {
			if op.Ext == filepath.Ext(fp) {
				matchedStrings, err := op.Match(fp, pattern, fileoonly)
				if err != nil {
					stderrLog.Println(err)
				}

				fp = strings.Replace(fp, string(fp), "\x1b[36m"+string(fp)+"\x1b[0m", -1)
				if matchedStrings != nil {
					if fileoonly {
						fmt.Printf("%v\n", fp)
					} else {
						for _, s := range matchedStrings {
							s = strings.Replace(s, string(pattern), "\x1b[31m"+string(pattern)+"\x1b[0m", -1)
							fmt.Printf("%v: %v\n", fp, s)
						}
					}
				}
			}
		}
	}
}
