package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/eferhatg/go-proto-tcp-hub/pkg/client"
	"github.com/eferhatg/go-proto-tcp-hub/pkg/protocol"
	"github.com/golang/protobuf/proto"
)

func main() {

	conn, err := net.Dial("tcp", ":1087")
	if conn == nil {
		log.Fatal("No connection")
		return
	}
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	c := client.NewClient(conn)
	go handle(c)
	for {

		//Reading choice from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Write text: ")
		text, _ := reader.ReadString('\n')

		//Writing to socket predefined messages depending on the choice
		m := protocol.Message{}
		switch strings.Trim(text, "\n") {
		case "identity":
			sendIdentity(c)
		case "list":
			sendList(c)
		case "relay":
			sendRelay(c)
		default:
			fmt.Print("Command not found\n")
		}

		if proto.Size(&m) > 0 {
			bt, _ := proto.Marshal(&m)
			c.Write(bt)
		}

	}
}

func sendIdentity(c *client.Client) {
	m := protocol.Message{
		Command: protocol.Message_IDENTITY,
	}
	if proto.Size(&m) > 0 {
		bt, _ := proto.Marshal(&m)
		c.Write(bt)
	}
}

func sendList(c *client.Client) {
	m := protocol.Message{
		Command: protocol.Message_LIST,
	}
	if proto.Size(&m) > 0 {
		bt, _ := proto.Marshal(&m)
		c.Write(bt)
	}
}

func sendRelay(c *client.Client) {
	to := []uint64{1, 2, 3, 4, 5, 6}
	bodyType := protocol.Message_PLAIN_TEXT
	body := []byte("This can be anything up to 1mb")
	m := protocol.Message{
		Command:  protocol.Message_RELAY,
		Id:       c.UserID,
		RelayTo:  to,
		BodyType: bodyType,
		Body:     body,
	}
	if proto.Size(&m) > 0 {
		bt, _ := proto.Marshal(&m)
		c.Write(bt)
	}
}

func handle(c *client.Client) {
	for {
		b, err := c.Read()
		if err != nil {
			log.Printf(err.Error())
		}
		nm := &protocol.Message{}
		proto.Unmarshal(b, nm)
		handleMsg(nm)
		if err == io.EOF {
			log.Println("HUB SHUT DOWN")
			break
		}
	}
	os.Exit(1)
}

func handleMsg(m *protocol.Message) {
	switch m.GetCommand() {
	case protocol.Message_IDENTITY:
		//Got msg.Id
		id := m.GetId()
		fmt.Printf("\nMy ID is %s", strconv.FormatUint(id, 10)+"\n")

	case protocol.Message_LIST:
		//Got msg.ConnectedClientIds
		ids := m.GetConnectedClientIds()

		if len(ids) == 0 {
			fmt.Printf("\nYou are the only client connected\n")
		} else {
			fmt.Printf("\nConnected client ids are [%s]\n", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]"))
		}

	case protocol.Message_RELAY:
		/*
			Got msg.Body and msg.BodyType
			You can convert/unmarshal body according to bodytype
		*/
		fmt.Printf("\nYou have a message from %s\n", strconv.FormatUint(m.GetId(), 10))
		switch m.BodyType {
		case protocol.Message_PLAIN_TEXT:
			fmt.Printf("\n\"%s\"\n", string(m.GetBody()))
		case protocol.Message_JSON:
			//Unmarshal body to json
		case protocol.Message_ERROR:
			err := fmt.Errorf("Error:%s", string(m.GetBody()))
			fmt.Printf(err.Error())
		}
	}
}
