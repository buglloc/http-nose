package main

import (
	"flag"
	"fmt"
	"log"
	"encoding/json"
	"strconv"
	"strings"
	"github.com/buglloc/http-nose/tcp"
	"github.com/buglloc/http-nose/httpclient"
)

func formatResponce(body string, contentType string) string {
	head := fmt.Sprintf(
		"HTTP/1.1 200 Ok\r\nServer: Trace Server\r\nContent-Type: %s\r\nContent-Length: %d\r\nConnection: close\r\n\r\n",
		contentType, len(body))
	return head + body
}

func main() {
	portFlag := flag.Int("port", 9000, "Port to bind")
	traceFlag := flag.Bool("trace", false, "Trace mode (analog of HTTP TRACE method)")
	flag.Parse()

	server := tcp.New(fmt.Sprintf(":%d", *portFlag))

	server.OnNewClient(func(c *tcp.Client) {
		log.Printf("New client: %s", c.Conn().RemoteAddr())
	})

	server.OnNewMessage(func(c *tcp.Client, message []byte) {
		log.Printf("New message from: %s", c.Conn().RemoteAddr())

		contentType := "application/json"
		var body string
		if *traceFlag {
			contentType = "text/plain"
			body = strconv.Quote(string(message))
			body = strings.Replace(body, "\\n", "\\n\n", -1)
			body = strings.Trim(body, "\"")
		} else {
			parsed, err := httpclient.ParseRequest(message)
			if err != nil {
				log.Print("Can't parse request:", err)
				c.Close()
			}
			for _, h := range parsed.Headers {
				if strings.ToLower(h.Name) == "host" {
					parsed.Host = h.Value
					break
				}
			}
			encoded, err := json.MarshalIndent(parsed, "", "  ")
			if err != nil {
				log.Print("Failed to encode:", err)
				c.Close()
			}
			body = string(encoded)
		}
		c.Send(formatResponce(body, contentType))
		c.Close()
	})

	server.OnClientConnectionClosed(func(c *tcp.Client, err error) {
		log.Printf("Client disconnected: %s", c.Conn().RemoteAddr())
	})

	server.Listen()
}
