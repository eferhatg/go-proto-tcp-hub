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

//Hub holds
type Hub struct {
	clients   []*client.Client
	accept    chan *client.Client
	terminate chan bool
}

//NewHub initialize Hub
func NewHub() *Hub {
	return &Hub{
		clients:   make([]*client.Client, 0),
		accept:    make(chan *client.Client),
		terminate: make(chan bool),
	}
}

//listen starts channel listening
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

//Start starts listening tcp port
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
		if <-h.terminate {
			break
		}
	}
	return nil

}

//acceptClient accepts clients
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
