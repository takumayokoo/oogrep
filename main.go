package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	OoFile{Ext: ".xlsx", TextFilePattern: `xl/sharedStrings\.xml`, Re: xlsxRe},
	OoFile{Ext: ".pptx", TextFilePattern: `ppt/slides/slide[0-9]+\.xml`, Re: pptxRe},
}

func targetFiles() []string {
	var filepaths []string

	if recursive {
		var dirs []string

		if flag.NArg() > 1 {
			dirs = flag.Args()[1:]
		} else {
			d, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			dirs = append(dirs, d)
		}

		for _, d := range dirs {

			err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					fullpath := path
					filepaths = append(filepaths, fullpath)
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
					log.Fatal(err)
				}

				if matchedStrings != nil {
					if fileoonly {
						fmt.Printf("%v\n", fp)
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
