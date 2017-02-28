package parser

import (
	"io"
	"errors"
	"bufio"
	"strings"
	"github.com/buglloc/http-nose/httpclient"
	//"fmt"
)

type Request struct {
	httpclient.Request
	BodyReader        *bufio.Reader `json:"-"`
	mayHaveBody       bool
	contentLength     int64
	contentLengthRead bool
}

func NewRequest() *Request {
	return &Request{}
}

func (pr *Request) Parse(rd io.Reader) error {
	reader := bufio.NewReader(rd)

	line, err := reader.ReadString('\n')
	pr.Raw += line
	if err != nil {
		return err
	}

	err = pr.parseRequestLine(line)
	if err != nil {
		return err
	}

	err = pr.parseHeaders(reader)
	if err != nil {
		return err
	}

	if pr.mayHaveBody {
		pr.BodyReader = reader
	}
	return nil
}

func (pr *Request) parseRequestLine(line string) error {
	line = strings.TrimRight(line, "\r\n")
	splitted := strings.SplitN(line, " ", 3)
	if len(splitted) == 1 {
		return errors.New("No path in request line")
	}
	pr.Method = splitted[0]
	pr.RequestURI = splitted[1]
	pr.parseRequestURI(pr.RequestURI)
	if len(splitted) == 3 {
		pr.Proto = splitted[2]
	}
	return nil
}

func (pr *Request) parseRequestURI(uri string) error {
	uri_splitted := strings.SplitN(uri, "?", 2)
	pr.Path = uri_splitted[0]
	if len(uri_splitted) == 2 {
		pr.Args = uri_splitted[1]
	}

	return nil
}

func (pr *Request) parseHeaders(rd *bufio.Reader) error {
	for {
		line, err := rd.ReadString('\n')

		pr.Raw += line
		if err != nil {
			return err
		} else if err == io.EOF {
			pr.mayHaveBody = false
			break
		}

		// End of headers
		line = strings.TrimRight(line, "\r\n")

		if len(line) == 0 {
			pr.mayHaveBody = true
			break
		}

		if strings.IndexAny(line, " \t") == 0 {
			// Multiline header
			last := len(pr.Headers)
			if last == 0 {
				return errors.New("Bad header: " + line)
			}
			pr.Headers[last-1].Value += line
			continue
		}

		splitted := strings.SplitN(line, ":", 2)
		h := httpclient.Header{Name: splitted[0]}
		if len(splitted) == 2 {
			h.Value = strings.Trim(splitted[1], " ")
		}
		pr.Headers = append(pr.Headers, h)
	}
	return nil
}