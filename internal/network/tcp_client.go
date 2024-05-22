package network

import (
	"fmt"
	"net"
	"time"
)

type TCPClient struct {
	conn           net.Conn
	dialTimeout    time.Duration
	idleTimeout    time.Duration
	maxMessageSize int
}

func NewTCPClient(address string, maxMessageSize int, idleTimeout time.Duration) (*TCPClient, error) {
	if idleTimeout <= 0 {
		return nil, fmt.Errorf("idle timeout must be greater than zero")
	}

	if maxMessageSize <= 0 {
		return nil, fmt.Errorf("max message size must be greater than zero")
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %w", err)
	}

	return &TCPClient{
		conn:           conn,
		idleTimeout:    idleTimeout,
		maxMessageSize: maxMessageSize,
	}, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if err := c.conn.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := c.conn.Write(request); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	response := make([]byte, c.maxMessageSize)
	count, err := c.conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return response[:count], nil
}
