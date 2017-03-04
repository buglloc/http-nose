package httpclient

import (
	"strings"
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

func (h *Header) EqualName(name string) bool {
	return NormalizeHeaderName(h.Name) == NormalizeHeaderName(name)
}

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

func (r *Request) HeadersSlice(name string) []Header {
	testName := NormalizeHeaderName(name)
	result := make([]Header, 0, 1)
	for _, h := range r.Headers {
		if NormalizeHeaderName(h.Name) == testName {
			result = append(result, h)
		}
	}
	return result
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
	res = append(res, []byte(r.BuildRequestUri())...)
	if len(r.Proto) > 0 {
		res = append(res, delim...)
		res = append(res, []byte(r.Proto)...)
	}
	return res
}

func (r *Request) BuildRequestUri() []byte {
	if r.RequestURI != "" {
		return []byte(r.RequestURI)
	}
	res := make([]byte, 0)
	res = append(res, []byte(r.Path)...)
	if r.Args != "" {
		res = append(res, '?')
		res = append(res, []byte(r.Args)...)
	}
	return res
}

func (r *Request) BuildHeaders(delim []byte) [][]byte {
	if len(delim) == 0 {
		delim = []byte(":")
	}

	headers := r.Headers
	if r.Host != "" {
		headers = append([]Header{{"Host", r.Host}}, r.Headers...)
	}
	res := make([][]byte, len(headers))
	for i, h := range headers {
		res[i] = append(res[i], []byte(h.Name)...)
		res[i] = append(res[i], delim...)
		res[i] = append(res[i], []byte(h.Value)...)
	}
	return res
}

func NormalizeHeaderName(name string) string {
	return strings.Trim(strings.ToLower(name), " \r\n")
}