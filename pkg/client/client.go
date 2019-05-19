package client

import (
	"bufio"
	"log"
	"net"

	"github.com/eferhatg/uinty-assignment/pkg/protocol"
)

const MAX_MSG_SIZE int = 1048576

type Client struct {
	UserID   uint64
	Conn     *net.Conn
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Incoming chan *protocol.Message
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:   &conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
}

func (c *Client) Write(data []byte) error {

	_, err := c.Writer.Write(data)
	if err != nil {
		return err
	}

	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Read() ([]byte, error) {
	buf := make([]byte, 1024)
	len, err := c.Reader.Read(buf)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return buf[:len], nil
}
