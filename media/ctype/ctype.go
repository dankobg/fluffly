package ctype

import (
	"bytes"
	"io"
	"net/http"
)

func DetectContentType(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", r, err
	}
	contentType := http.DetectContentType(buf[:n])
	r = io.MultiReader(bytes.NewReader(buf[:n]), r)
	return contentType, r, nil
}
