package utility

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var fileContents = map[string]string{
	"a.txt": "123",
	"b.txt": "中文",
}

func TestZIPWithFile(t *testing.T) {
	tmpDir := os.TempDir()
	for filename, content := range fileContents {
		path := filepath.Join(tmpDir, filename)
		file, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		_, _ = io.WriteString(file, content)
		_ = file.Close()
	}
	zipPath := filepath.Join(tmpDir, "zip.zip")
	zipFiles := map[string]string{
		filepath.Join(tmpDir, "a.txt"): "",
		filepath.Join(tmpDir, "b.txt"): "a/b/index.txt",
	}
	{
		zipWriter, err := NewZipArchiveWriterToFile(zipPath)
		if err != nil {
			t.Error(err)
			return
		}
		for fileSrc, fileInZip := range zipFiles {
			_, err = zipWriter.WriteFile(fileSrc, fileInZip, true)
			if err != nil {
				t.Error(err)
				return
			}
		}
		zipWriter.Close()
	}
	for file := range zipFiles {
		_ = os.Remove(file)
	}
	zipReader, err := NewZipArchiveReaderFromFile(zipPath)
	if err != nil {
		t.Error(err)
	} else {
		extractedFiles, err := zipReader.Extract(tmpDir, true)
		if err != nil {
			t.Error(err)
		}
		for _, file := range extractedFiles {
			_ = os.Remove(file)
		}
	}
	_ = os.Remove(zipPath)
}

func TestReZIP(t *testing.T) {
	// TODO: make it work
	dir, _ := os.Getwd()
	zipFileSRC := dir +
		"FF383829-D168-4CED-B4A3-DC4FF7A34E2D.zip"
	// "149A8201-A74C-49B4-B52F-9DA40E70ED6C.zip"
	zipFileTAR := dir + "t.zip"
	zipReader, err := NewZipArchiveReaderFromFile(zipFileSRC)
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	zipWriter := NewZipArchiveWriter(&buf)
	var zFile *ZipArchiveInnerFile
	files := make(map[string]int)
	for {
		zFile, err = zipReader.ReadNextFile()
		if err != nil {
			t.Error(err)
			return
		}
		if zFile == nil {
			break
		}
		wroteCount, err := zipWriter.WriteBytes(zFile.Content, zFile.Name, true)
		if err != nil {
			t.Error(err)
			return
		}
		files[zFile.Name] = wroteCount
	}
	zipWriter.Close()
	bs := buf.Bytes()
	_, err = WriteBytesToFile(bs, zipFileTAR)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestZIPInMemory(t *testing.T) {
	var bs []byte
	var err error
	{
		var buf bytes.Buffer
		zipWriter := NewZipArchiveWriter(&buf)
		for filename, content := range fileContents {
			fmt.Println("writing", filename, content)
			_, err = zipWriter.WriteBytes([]byte(content), filename, true)
			if err != nil {
				t.Error(err)
				break
			}
		}
		zipWriter.Close()
		bs = buf.Bytes()
		if err != nil {
			return
		}
	}
	fmt.Println("bs.len=", len(bs))
	{
		//path := filepath.Join(os.TempDir(), "zip.zip")
		//fp, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		//fp.Write(bs)
		//fp.Close()
		//fmt.Println(path)
	}
	reader := bytes.NewReader(bs)
	zipReader, err := NewZipArchiveReader(reader, int64(len(bs)))
	if err != nil {
		t.Error(err)
		return
	}
	var zipFile *ZipArchiveInnerFile
	for {
		zipFile, err = zipReader.ReadNextFile()
		if err != nil {
			t.Error(err)
			break
		}
		if zipFile == nil {
			break
		}
		fmt.Println(zipFile.Name, len(zipFile.Content), string(zipFile.Content))
	}
}
