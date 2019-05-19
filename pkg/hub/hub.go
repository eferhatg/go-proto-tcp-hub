package hub

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/eferhatg/uinty-assignment/pkg/client"
)

type Hub struct {
	clients []*client.Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make([]*client.Client, 0),
	}
}

func (h *Hub) Listen(startport int) error {
	port := strconv.Itoa(startport)

	ll, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Panicln("Error: ", err)
		return err
	}

	conn, err := ll.Accept()
	if err != nil {
		log.Panicln("Error: ", err)
		return err
	}
	c := client.NewClient(conn)
	h.clients = append(h.clients, c)

	for {

		b, err := c.Read()
		message := string(b)
		fmt.Print("Got:", message+"\n")

		newmessage := strings.ToUpper(message)

		c.Write([]byte(newmessage + "\n"))

		if err == io.EOF {
			break
		}

	}
	return nil
}
