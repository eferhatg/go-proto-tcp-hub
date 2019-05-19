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
	accept  chan *client.Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make([]*client.Client, 0),
		accept:  make(chan *client.Client),
	}
}
func (h *Hub) listen() {

	go func() {

		for {
			select {
			case client := <-h.accept:
				h.acceptClient(client)

			}

		}
	}()
}

func (h *Hub) Start(startport int) error {
	port := strconv.Itoa(startport)

	ll, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Panicln("Error: ", err)
		return err
	}
	h.listen()
	for {
		conn, err := ll.Accept()
		if err != nil {
			log.Panicln("Error: ", err)
		}
		c := client.NewClient(conn)
		h.clients = append(h.clients, c)
		h.accept <- c
	}

}

func (h *Hub) acceptClient(c *client.Client) error {

	go func() {
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
	}()

	return nil
}
