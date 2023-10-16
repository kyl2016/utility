package utility

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Attribute prefix for file in app directory
	FilePrefix = "@file:"
)

// filename is like: @file:filename.txt
func ExtractFilename(value string) string {
	if strings.HasPrefix(value, FilePrefix) {
		return strings.Trim(value[len(FilePrefix):], " \r\n")
	}
	return ""
}

func ListFilesAtDirectory(dir string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			files = append(files, path)
		}
		return err
	})
	return
}

func MoveFile(src, dest string, overrite bool) error {
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if !overrite {
			return os.ErrExist
		}
		if err := os.RemoveAll(dest); err != nil {
			return err
		}
	}
	// check parent-dir of dest
	destParentDir := filepath.Dir(dest)
	if _, err := os.Stat(destParentDir); os.IsNotExist(err) {
		if err = os.MkdirAll(destParentDir, os.ModePerm); err != nil {
			return err
		}
	}
	return os.Rename(src, dest)
}

func ReadAllFromReadCloser(r io.ReadCloser) ([]byte, error) {
	var bs []byte
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	bs = buf.Bytes()
	return bs, err
}

func ReadBytes(fp io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, fp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadFileFromPath(path string) ([]byte, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return ReadBytes(fp)
}

func WriteBytesToFile(bs []byte, path string) (int, error) {
	return WriteBytesToFileWithPermission(bs, path, 0755)
}

func WriteBytesToFileWithPermission(bs []byte, path string, perm os.FileMode) (int, error) {
	var oldFileSize int64
	if file, err := os.Stat(path); !os.IsNotExist(err) {
		oldFileSize = file.Size()
	}
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return 0, err
	}
	count, err := fp.Write(bs)
	if oldFileSize > int64(count) {
		// truncate if old file exists and length < new bytes
		fp.Truncate(int64(count))
	}
	fp.Close()
	return count, err
}
