package utility

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func GetAllFiles(root, extension string) (files []string, err error) {
	extension = "." + strings.Trim(extension, ".")
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			files = append(files, path)
		}
		return err
	})
	return
}

func GetLines(file string) ([]string, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var r []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		r = append(r, scanner.Text())
	}
	return r, nil
}

func Contains(origin string, sub ...string) bool {
	for _, s := range sub {
		if strings.Contains(origin, s) {
			return true
		}
	}

	return false
}
