package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
)

type OoFile struct {
	Ext             string
	TextFilePattern string
	Re              *regexp.Regexp
}

func readFilesFromZip(f *zip.ReadCloser, textFilePattern string) ([]byte, error) {
	var inBuf []byte
	found := false
	buf := bytes.NewBuffer(inBuf)
	fnRegexp := regexp.MustCompile(textFilePattern)

	for _, v := range f.File {
		if fnRegexp.MatchString(v.Name) {
			found = true

			doc, err := v.Open()
			if err != nil {
				return nil, err
			}

			defer doc.Close()

			bs, err := ioutil.ReadAll(doc)
			if err != nil {
				return nil, err
			}

			_, err = buf.Write(bs)
			if err != nil {
				return nil, err
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("%v not found.", textFilePattern)
	}

	return buf.Bytes(), nil
}

func (f *OoFile) Match(filepath string, pattern []byte, onlyfirst bool) ([]string, error) {
	z, err := zip.OpenReader(filepath)
	if err == zip.ErrFormat {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%v: %v", filepath, err)
	}
	defer z.Close()

	bs, err := readFilesFromZip(z, f.TextFilePattern)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", filepath, err)
	}

	var ret []string
	for _, b := range f.Re.FindAllSubmatch(bs, -1) {
		contents := b[1]
		if bytes.Contains(contents, pattern) {
			con := string(contents)
			//			log.Println(con)
			if onlyfirst {
				return []string{con}, nil
			} else {
				ret = append(ret, con)
			}
		}
	}

	return ret, nil
}
