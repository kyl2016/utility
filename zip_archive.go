package utility

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type ZipArchiveReader interface {
	ReadNextFile() (*ZipArchiveInnerFile, error)
	Extract(targetDir string, overwrite bool) ([]string, error)
}

type ZipArchiveInnerFile struct {
	Name     string
	Content  []byte
	FileInfo os.FileInfo
}

func NewZipArchiveReader(r io.ReaderAt, size int64) (ZipArchiveReader, error) {
	zipFP, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	inst := &zipArchiveReader{reader: zipFP}
	return inst, nil
}

func NewZipArchiveReaderFromFile(zipPath string) (ZipArchiveReader, error) {
	zipFP, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	inst := &zipArchiveReader{reader: &zipFP.Reader}
	runtime.SetFinalizer(inst, func(z *zipArchiveReader) {
		_ = zipFP.Close()
	})
	return inst, nil
}

type zipArchiveReader struct {
	reader *zip.Reader
	index  int
}

func (z *zipArchiveReader) Rewind() {
	z.index = 0
}

/// Read next file from reader, skip all directory.
func (z *zipArchiveReader) ReadNextFile() (*ZipArchiveInnerFile, error) {
	var file *zip.File
	var info os.FileInfo
	for {
		file = z.nextFile()
		if file == nil {
			return nil, nil
		}
		info = file.FileInfo()
		if !info.IsDir() {
			break
		}
	}
	var err error
	inst := &ZipArchiveInnerFile{Name: file.Name, FileInfo: info}
	if !info.IsDir() {
		var fp io.ReadCloser
		fp, err = file.Open()
		if err == nil {
			var buf bytes.Buffer
			_, err = io.Copy(&buf, fp)
			_ = fp.Close()
			if err == io.EOF {
				err = nil
			}
			if err == nil {
				inst.Content = buf.Bytes()
			}
		} else {
			err = fmt.Errorf("read file %v failed: %w", file.Name, err)
		}
	}
	return inst, err
}

func (z *zipArchiveReader) nextFile() *zip.File {
	if z.reader != nil && z.index < len(z.reader.File) {
		i := z.index
		z.index += 1
		return z.reader.File[i]
	}
	return nil
}

func (z *zipArchiveReader) Extract(targetDir string, overwrite bool) ([]string, error) {
	fileStatus, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return nil, err
		}
	} else if !fileStatus.IsDir() {
		if overwrite {
			err = os.Remove(targetDir)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unzip target dir exist but no overwrite allowed")
		}
	}
	extractedFiles := make([]string, len(z.reader.File))
	z.Rewind()
	for {
		file := z.nextFile()
		if file == nil {
			break
		}
		fpSrc, err := file.Open()
		if err != nil {
			return nil, err
		}
		shouldCopy := true
		pathTar := filepath.Join(targetDir, file.Name)
		if status, _ := os.Stat(pathTar); status != nil {
			if overwrite {
				err = os.Remove(pathTar)
				shouldCopy = err == nil
			} else {
				shouldCopy = false
			}
		}
		if shouldCopy {
			if file.FileInfo().IsDir() {
				err = os.MkdirAll(pathTar, file.Mode())
			} else {
				dirTar := filepath.Dir(pathTar)
				file, _ := os.Stat(dirTar)
				if file == nil || !file.IsDir() {
					_ = os.MkdirAll(dirTar, 0755)
				}
				var fpTar *os.File
				fpTar, err = os.OpenFile(pathTar, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				if err == nil {
					_, err = io.Copy(fpTar, fpSrc)
					if err == io.EOF {
						err = nil
					}
				}
				if fpTar != nil {
					_ = fpTar.Close()
				}
				extractedFiles = append(extractedFiles, pathTar)
			}
		}
		_ = fpSrc.Close()
		if err != nil {
			return extractedFiles, err
		}
	}
	return extractedFiles, nil
}

type ZipArchiveWriter interface {
	WriteFile(file string, pathInZIP string, deflate bool) (int64, error)
	WriteBytes(bs []byte, pathInZIP string, deflate bool) (int, error)
	Close()
}

func NewZipArchiveWriterToFile(path string) (ZipArchiveWriter, error) {
	zipFP, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	inst := &zipArchiveWriter{writer: zip.NewWriter(zipFP)}
	runtime.SetFinalizer(inst, func(z *zipArchiveWriter) {
		_ = zipFP.Close()
	})
	return inst, nil
}

func NewZipArchiveWriter(w io.Writer) ZipArchiveWriter {
	return &zipArchiveWriter{writer: zip.NewWriter(w)}
}

type zipArchiveWriter struct {
	writer *zip.Writer
}

func (z *zipArchiveWriter) WriteFile(file string, pathInZIP string, deflate bool) (int64, error) {
	fpSrc, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer fpSrc.Close()
	info, err := fpSrc.Stat()
	if err != nil {
		return 0, err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return 0, err
	}
	if len(pathInZIP) > 0 {
		header.Name = pathInZIP
	}
	if deflate {
		header.Method = zip.Deflate
	} else {
		header.Method = zip.Store
	}
	writer, err := z.writer.CreateHeader(header)
	if err != nil {
		return 0, err
	}
	return io.Copy(writer, fpSrc)
}

// WriteBytes as zip-inner file.
// Returns: wrote bytes count if no error.
// Panic: if writer wasn't initialize.
func (z *zipArchiveWriter) WriteBytes(bs []byte, pathInZIP string, deflate bool) (int, error) {
	header := &zip.FileHeader{Name: pathInZIP}
	header.UncompressedSize64 = uint64(len(bs))
	if deflate {
		header.Method = zip.Deflate
	} else {
		header.Method = zip.Store
	}
	header.SetMode(0755)
	header.SetModTime(time.Now())
	writer, err := z.writer.CreateHeader(header)
	if err != nil {
		return 0, err
	}
	return writer.Write(bs)
}

func (z *zipArchiveWriter) Close() {
	if z.writer != nil {
		_ = z.writer.Close()
	}
}
