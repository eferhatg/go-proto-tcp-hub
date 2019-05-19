package client

import (
	"bufio"
	"io"
	"log"
	"net"
)

const MAX_MSG_SIZE int = 1048576

type Client struct {
	UserID uint64
	Conn   *net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func NewClient(conn net.Conn, reader io.Reader, writer io.Writer) *Client {
	return &Client{
		Conn:   &conn,
		Reader: bufio.NewReader(reader),
		Writer: bufio.NewWriter(writer),
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
