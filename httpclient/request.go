package httpclient

import (
	"io"
	"strings"
	"bytes"
)

type Request struct {
	Method     string
	RequestURI string
	Path       string
	Args       string
	Proto      string
	Headers    []Header
	Host       string
	Body       string
	Raw        string
}

type Header struct {
	Name  string
	Value string
}

const (
	REQUEST_STATE_REQUEST_LINE = 0
	REQUEST_STATE_HEADER       = 1
	REQUEST_STATE_BODY         = 2
	REQUEST_STATE_END          = 3
)

func (r *Request) Clone() *Request {
	return &Request{
		Method:     r.Method,
		RequestURI: r.RequestURI,
		Path:       r.Path,
		Args:       r.Args,
		Proto:      r.Proto,
		Headers:    r.Headers,
		Body:       r.Body,
		Host:       r.Host,
	}
}

func (r *Request) AddHeader(name string, value string) {
	r.Headers = append(r.Headers, Header{Name: name, Value: value})
}

func (r *Request) RemoveHeader(name string) {
	testName := strings.ToLower(name)
	for i, h := range r.Headers {
		if strings.ToLower(h.Name) == testName {
			r.Headers = append(r.Headers[:i], r.Headers[i+1:]...)
		}
	}
}

func (r *Request) Build(lineDelim []byte, headerDelim []byte) []byte {
	if len(lineDelim) == 0 {
		lineDelim = []byte("\r\n")
	}
	if len(headerDelim) == 0 {
		headerDelim = []byte(":")
	}

	res := make([]byte, 0)
	res = append(res, r.BuildRequestLine(nil)...)
	res = append(res, lineDelim...)
	for _, h := range r.BuildHeaders(headerDelim) {
		res = append(res, h...)
		res = append(res, lineDelim...)
	}
	res = append(res, lineDelim...)
	res = append(res, r.Body...)
	return res
}

func (r *Request) BuildRequestLine(delim []byte) []byte {
	if len(delim) == 0 {
		delim = []byte(" ")
	}

	res := make([]byte, 0)
	res = append(res, []byte(r.Method)...)
	res = append(res, delim...)
	res = append(res, []byte(r.RequestURI)...)
	if len(r.Proto) > 0 {
		res = append(res, delim...)
		res = append(res, []byte(r.Proto)...)
	}
	return res
}

func (r *Request) BuildHeaders(delim []byte) [][]byte {
	if len(delim) == 0 {
		delim = []byte(":")
	}

	res := make([][]byte, len(r.Headers))
	for i, h := range r.Headers {
		res[i] = append(res[i], []byte(h.Name)...)
		res[i] = append(res[i], delim...)
		res[i] = append(res[i], []byte(h.Value)...)
	}
	return res
}