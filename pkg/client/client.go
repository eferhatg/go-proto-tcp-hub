package client

import (
	"bufio"
	"log"
	"net"

	"github.com/eferhatg/uinty-assignment/pkg/protocol"
)

//MaxMsgSize keeps maximum size of message
const MaxMsgSize int = 1048576

//MaxRelayClientCount keeps maximum size of message
const MaxRelayClientCount int = 255

//Client keeps struct
type Client struct {
	UserID   uint64
	Conn     *net.Conn
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Incoming chan *protocol.Message
}

//NewClient initialize new Client
func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:   &conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
}

//Write writes data to Writer
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

//Read reads data from Reader
func (c *Client) Read() ([]byte, error) {
	buf := make([]byte, MaxMsgSize)
	len, err := c.Reader.Read(buf)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return buf[:len], nil
}
