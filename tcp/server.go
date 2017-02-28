package tcp

import (
	"net"
	"log"
	"io"
	"time"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
}

// TCP server
type server struct {
	clients                  []*Client
	address                  string // Address to open connection: localhost:9999
	timeout                  time.Duration
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, rd io.Reader)
}

// Read client data from channel
func (c *Client) listen() {
	c.Server.onNewMessage(c, c.conn)
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, rd io.Reader)) {
	s.onNewMessage = callback
}

// Start network server
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server. ", err)
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		err := conn.SetReadDeadline(time.Now().Add(s.timeout))
		if err != nil {
			log.Fatal("SetReadDeadline failed:", err)
		}

		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// Creates new tcp server instance
func New(address string, timeout int) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
		timeout: time.Duration(5 * time.Second),
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, rd io.Reader) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}
