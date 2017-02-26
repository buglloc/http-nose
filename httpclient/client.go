package httpclient

import (
	"net"
	"time"
	"log"
	"github.com/buglloc/http-nose/tcp"
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
	message, err := tcp.ReadAll(conn)
	if err != nil {
		log.Print("Failed to read response: ", err)
		return nil, err
	}

	resp, err := ParseResponse(message)
	if err != nil {
		log.Print("Failed to parse response: ", err)
		return nil, err
	}

	return resp, nil
}
