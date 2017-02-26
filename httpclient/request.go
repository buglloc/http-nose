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

func ParseRequest(request []byte) (*Request, error) {
	req := &Request{
		Raw: string(request),
	}

	state := REQUEST_STATE_REQUEST_LINE
	reader := bytes.NewBuffer(request)
	parse := true
	for parse {
		switch state {
		case REQUEST_STATE_REQUEST_LINE:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				state = REQUEST_STATE_END
				break
			}
			if err != nil {
				return nil, err
			}

			line = strings.TrimRight(line, "\r\n")
			splitted := strings.SplitN(line, " ", 3)
			if len(splitted) >= 1 {
				req.Method = splitted[0]
			}
			if len(splitted) >= 2 {
				req.RequestURI = splitted[1]
				uri_splitted := strings.SplitN(splitted[1], "?", 2)
				if len(uri_splitted) >= 1 {
					req.Path = uri_splitted[0]
				}
				if len(uri_splitted) == 2 {
					req.Args = uri_splitted[1]
				}
			}
			if len(splitted) == 3 {
				req.Proto = splitted[2]
			}
			state = REQUEST_STATE_HEADER
		case REQUEST_STATE_HEADER:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				state = REQUEST_STATE_END
				break
			}
			if err != nil {
				return nil, err
			}

			line = strings.TrimRight(line, "\r\n")
			if len(line) == 0 {
				state = REQUEST_STATE_BODY
				break
			}

			splitted := strings.SplitN(line, ":", 2)
			h := Header{Name: splitted[0]}
			if len(splitted) == 2 {
				h.Value = strings.TrimLeft(splitted[1], " ")
			}
			req.Headers = append(req.Headers, h)
		case REQUEST_STATE_BODY:
			var buf bytes.Buffer
			io.Copy(&buf, reader)
			req.Body = buf.String()
			state = REQUEST_STATE_END
		case REQUEST_STATE_END:
			parse = false
		}
	}
	return req, nil
}
