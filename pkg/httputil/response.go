package httputil

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
)

type ResponseCapture struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	body        *bytes.Buffer
	hijacker    http.Hijacker
}

func NewResponseCapture(w http.ResponseWriter) *ResponseCapture {
	hijacker, _ := w.(http.Hijacker)
	return &ResponseCapture{
		ResponseWriter: w,
		wroteHeader:    false,
		body:           new(bytes.Buffer),
		hijacker:       hijacker,
	}
}

func (its ResponseCapture) Header() http.Header {
	return its.ResponseWriter.Header()
}

func (its ResponseCapture) Write(data []byte) (int, error) {
	if !its.wroteHeader {
		its.WriteHeader(http.StatusOK)
	}
	its.body.Write(data)
	return its.ResponseWriter.Write(data)
}

func (its *ResponseCapture) WriteHeader(statusCode int) {
	its.status = statusCode
	its.wroteHeader = true
	its.ResponseWriter.WriteHeader(statusCode)
}

func (its ResponseCapture) Bytes() []byte {
	return its.body.Bytes()
}

func (its ResponseCapture) StatusCode() int {
	return its.status
}

func (its *ResponseCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if its.hijacker == nil {
		return nil, nil, errors.New("http.Hijacker not implemented by underlying http.ResponseWriter")
	}
	return its.hijacker.Hijack()
}
