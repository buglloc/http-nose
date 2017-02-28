package main

import (
	"flag"
	"fmt"
	"log"
	"encoding/json"
	"strconv"
	"strings"
	"io"
	"github.com/buglloc/http-nose/tcp"
	"github.com/buglloc/http-nose/cmd/trace-server/parser"
)

func formatResponce(body string, contentType string) string {
	head := fmt.Sprintf(
		"HTTP/1.1 200 Ok\r\nServer: Trace Server\r\nContent-Type: %s\r\nContent-Length: %d\r\nConnection: close\r\n\r\n",
		contentType, len(body))
	return head + body
}

func parseMessage(rd io.Reader) (*parser.Request, error) {
	request := parser.NewRequest()
	err := request.Parse(rd)
	if err != nil {
		return nil, err
	}

	if request.BodyReader != nil {
		if request.BodyReader.Buffered() != 0 {
			rest, err := tcp.ReadAll(request.BodyReader)
			if err == nil {
				request.Body = string(rest)
				request.Raw += request.Body
			}

		}
		request.BodyReader = nil
	}

	return request, nil
}

func main() {
	portFlag := flag.Int("port", 9000, "Port to bind")
	traceFlag := flag.Bool("trace", false, "Trace mode (analog of HTTP TRACE method)")
	flag.Parse()

	server := tcp.New(fmt.Sprintf(":%d", *portFlag), 2)

	server.OnNewClient(func(c *tcp.Client) {
		log.Printf("New client: %s", c.Conn().RemoteAddr())
	})

	server.OnNewMessage(func(c *tcp.Client, rd io.Reader) {
		log.Printf("New message from: %s", c.Conn().RemoteAddr())

		req, err := parseMessage(rd)
		if err != nil {
			log.Println("Failed to parse request: ", err)
			c.Close()
			return
		}

		contentType := "application/json"
		var body string
		if *traceFlag {
			contentType = "text/plain"
			body = strconv.Quote(req.Raw)
			body = strings.Replace(body, "\\n", "\\n\n", -1)
			body = strings.Trim(body, "\"")
		} else {
			for _, h := range req.Headers {
				if strings.ToLower(h.Name) == "host" {
					req.Host = h.Value
					break
				}
			}
			encoded, err := json.MarshalIndent(req, "", "  ")
			if err != nil {
				log.Println("Failed to encode: ", err)
				c.Close()
			}
			body = string(encoded)
		}
		err = c.Send(formatResponce(body, contentType))
		if err != nil {
			log.Println("Failed to send response: ", err)
		}
		c.Close()
	})

	server.OnClientConnectionClosed(func(c *tcp.Client, err error) {
		log.Printf("Client disconnected: %s", c.Conn().RemoteAddr())
	})

	server.Listen()
}
