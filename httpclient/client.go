package httpclient

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

const timeout = 10

type Client struct {
	Target string
}

func (c *Client) MakeRequest(r *Request) (*Response, error) {
	return c.MakeRawRequest(r.Build(nil, nil))
}

func (c *Client) MakeRawRequest(r []byte) (*Response, error) {
	conn, err := net.Dial("tcp", c.Target)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
	defer conn.Close()
	_, err = conn.Write(r)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(conn)
	httpResp, err := http.ReadResponse(br, nil)
	if err != nil {
		//log.Print("Failed to read response: ", err)
		//log.Print("Request: ", string(r))
		return nil, err
	}

	response, err := NewFromHttpResponse(httpResp)
	if err != nil {
		log.Print("Failed to parse response: ", err)
		log.Print("Request: ", string(r))
		return nil, err
	}

	return response, nil
}
