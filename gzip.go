package utility

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GZipCompress(in []byte) ([]byte, error) {
	var buffer bytes.Buffer
	var err error
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err == nil {
		err = writer.Close()
	} else {
		_ = writer.Close()
	}
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GZipDecompress(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
