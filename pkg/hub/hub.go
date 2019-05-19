package hub

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/eferhatg/uinty-assignment/pkg/client"
	"github.com/eferhatg/uinty-assignment/pkg/protocol"
	"github.com/golang/protobuf/proto"
)

//Hub holds
type Hub struct {
	listener  net.Listener
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
				go h.acceptClient(client)
			}
		}
	}()
}

//Start starts listening tcp port
func (h *Hub) Start(startport int) error {
	port := strconv.Itoa(startport)

	ll, err := net.Listen("tcp", ":"+port)
	h.listener = ll
	if err != nil {
		log.Panicln("Error: ", err)
		return err
	}
	h.listen()
	for {
		conn, err := ll.Accept()
		log.Println("New client")
		if err != nil {
			log.Panicln("Error: ", err)
		}

		c := client.NewClient(conn)
		c.UserID = uint64(len(h.clients) + 1)
		h.clients = append(h.clients, c)
		log.Println(len(h.clients))
		h.accept <- c

	}

}

//acceptClient accepts clients
func (h *Hub) acceptClient(c *client.Client) error {

	for {

		b, err := c.Read()
		log.Println("Yeni Mesaj")
		m := &protocol.Message{}
		proto.Unmarshal(b, m)

		switch m.GetCommand() {
		case protocol.Message_IDENTITY:
			go h.identityResponse(c, m)
		case protocol.Message_LIST:
			go h.listResponse(c, m)
		case protocol.Message_RELAY:
			//	h.relayResponse(current, cmd)

		}
		if m.GetCommand() == protocol.Message_IDENTITY {
			h.identityResponse(c, m)
		}

		if err == io.EOF {
			break
		}

	}

	return nil
}

func (h *Hub) identityResponse(c *client.Client, m *protocol.Message) {

	m.Id = c.UserID
	bt, _ := proto.Marshal(m)

	c.Write(bt)
}

func (h *Hub) listResponse(c *client.Client, m *protocol.Message) {
	log.Print("TEST")
	m.Id = c.UserID
	list := []uint64{}
	for _, cli := range h.clients {
		if c.UserID != cli.UserID {
			list = append(list, cli.UserID)
		}
	}
	log.Print(list)
	m.ConnectedClientIds = list

	bt, _ := proto.Marshal(m)
	c.Write(bt)
}

func (h *Hub) relayResponse(current *client.Client, m *protocol.Message) {
	if len(m.GetBody()) > 1048576 {
		//err := errors.New("Body is bigger than 1024kb")
		return
	}
	if len(m.GetRelayTo()) > 255 {
		//err := errors.New("Reciever count is bigger than 255")
		return
	}

	bt, err := proto.Marshal(m)
	if err != nil {
		return
	}

	for _, cli := range h.clients {
		for _, id := range m.GetRelayTo() {
			if id == cli.UserID && id != m.Id {
				cli.Write(bt)
			}
		}
	}
}
