package hub

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/eferhatg/uinty-assignment/pkg/client"
	"github.com/eferhatg/uinty-assignment/pkg/protocol"
	"github.com/golang/protobuf/proto"
)

//Hub holds
type Hub struct {
	listener  net.Listener
	clients   map[uint64]*client.Client
	accept    chan *client.Client
	terminate chan bool
	mutex     *sync.Mutex
}

//NewHub initialize Hub
func NewHub() *Hub {
	return &Hub{
		clients:   make(map[uint64]*client.Client, 0),
		accept:    make(chan *client.Client),
		terminate: make(chan bool),
		mutex:     &sync.Mutex{},
	}
}

//listen starts channel listening
func (h *Hub) listen() {
	go func() {
		for {
			select {
			case client := <-h.accept:
				go h.handleClient(client)
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
		if err != nil {
			log.Panicln("Error: ", err)
		}

		c := client.NewClient(conn)
		c.UserID = uint64(len(h.clients) + 1)
		h.mutex.Lock()
		h.clients[c.UserID] = c
		h.mutex.Unlock()
		log.Printf("New client connected. Client id: %s. %s clients connected now.", strconv.FormatUint(c.UserID, 10), strconv.Itoa(len(h.clients)))
		h.accept <- c
	}

}

//handleClient listens and handles client connections
func (h *Hub) handleClient(c *client.Client) error {

	for {

		b, err := c.Read()

		m := &protocol.Message{}
		proto.Unmarshal(b, m)

		switch m.GetCommand() {
		case protocol.Message_IDENTITY:
			go h.identityResponse(c, m)
		case protocol.Message_LIST:
			go h.listResponse(c, m)
		case protocol.Message_RELAY:
			h.relayResponse(c, m)
		}

		if err == io.EOF {
			h.mutex.Lock()
			delete(h.clients, c.UserID)
			h.mutex.Unlock()
			log.Printf("A client disconnected. Dsconnected client id: %s. Total %s clients connected now.", strconv.FormatUint(c.UserID, 10), strconv.Itoa(len(h.clients)))

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

	m.Id = c.UserID
	list := []uint64{}
	for k, c := range h.clients {
		if c.UserID != k {
			list = append(list, k)
		}
	}
	log.Print(list)
	m.ConnectedClientIds = list

	bt, _ := proto.Marshal(m)
	c.Write(bt)
}

func (h *Hub) relayResponse(c *client.Client, m *protocol.Message) {
	if len(m.GetBody()) > 1048576 {
		m.BodyType = protocol.Message_ERROR
		m.Body = []byte("Body is bigger than 1024kb")
		bt, _ := proto.Marshal(m)
		c.Write(bt)
		return
	}
	if len(m.GetRelayTo()) > 255 {
		m.BodyType = protocol.Message_ERROR
		m.Body = []byte("Reciever count is bigger than 255")
		bt, _ := proto.Marshal(m)
		c.Write(bt)
		return
	}

	bt, err := proto.Marshal(m)
	if err != nil {
		return
	}

	for _, cli := range h.clients {
		for _, id := range m.GetRelayTo() {
			if id == cli.UserID && id != c.UserID {

				cli.Write(bt)
			}
		}
	}
}
