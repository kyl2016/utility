package utility

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func GetLinesFromFile(file string) ([]string, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			lines = append(lines, text)
		}
	}
	return lines, nil
}

func GetFilesFromDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, path.Join(dir, entry.Name()))
	}
	return files, nil
}

func SaveJsonObj(obj interface{}, toFile string) error {
	data, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err
	}
	return writeFile(data, toFile)
}

func SaveFromReader(reader io.Reader, toFile string) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return writeFile(data, toFile)
}

func writeFile(data []byte, toFile string) error {
	_, err := os.Stat(path.Dir(toFile))
	if err != nil {
		err = os.MkdirAll(path.Dir(toFile), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(toFile, data, os.ModePerm)
}

func Exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func FormatFileName(file string) string {
	file = strings.ReplaceAll(file, "/", "_")
	file = strings.ReplaceAll(file, "\\", "_")
	return file
}