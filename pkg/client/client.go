package client

import (
	"bufio"
	"log"
	"net"
)

const MAX_MSG_SIZE int = 1048576

type Client struct {
	Conn   *net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
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

	b, _, err := c.Reader.ReadLine()
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return b, nil
}