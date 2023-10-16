package utility

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

const (
	HTTPHeaderContentEncoding = "Content-Encoding"
	HTTPHeaderContentLength   = "Content-Length"
	HTTPHeaderContentType     = "Content-Type"

	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain"
	ContentTypeBIN  = "application/octet-stream"

	ContentTypeFormData = "multipart/form-data"
)

func ReadDecodedResponseBody(resp *http.Response) ([]byte, error) {
	body, err := ReadAllFromReadCloser(resp.Body)
	if err != nil {
		return nil, err
	}
	return ReadDecodedBodyWithHeader(body, resp.Header.Get(HTTPHeaderContentEncoding))
}

func ReadDecodedBodyWithHeader(body []byte, contentEncoding string) ([]byte, error) {
	switch contentEncoding {
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return ReadBytes(reader)
	case "deflate":
		reader := flate.NewReader(bytes.NewReader(body))
		defer reader.Close()
		return ReadBytes(reader)
	default:
		return body, nil
	}
}

func IsRequestFromLocalhost(req *http.Request) bool {
	return strings.Contains(req.Host, "127.0.0.1") || strings.Contains(req.Host, "localhost")
}

func SetBodyToRequest(req *http.Request, body []byte) {
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	req.Header.Set(HTTPHeaderContentLength, AnyToString(req.ContentLength))
}

func SetBodyToResponse(resp *http.Response, body []byte) {
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	resp.Header.Set(HTTPHeaderContentLength, AnyToString(resp.ContentLength))
}

// GetLocalhostIP 取得本机在局域网的IP
func GetLanIP() (net.IP, error) {
	ips, err := GetNetIPs()
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if !ip.IsLoopback() && !strings.Contains(ip.String(), ":") {
			return ip, nil
		}
	}
	return nil, errors.New("not found")
}

func GetNetIPs() (ips []net.IP, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, addrErr := i.Addrs()
		if addrErr != nil {
			err = addrErr
			return
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ips = append(ips, v.IP)
			case *net.IPAddr:
				ips = append(ips, v.IP)
			}
		}
	}
	return
}
