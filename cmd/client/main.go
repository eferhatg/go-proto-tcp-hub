package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/eferhatg/uinty-assignment/pkg/client"
	"github.com/eferhatg/uinty-assignment/pkg/protocol"
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
	body := []byte("This can be a string up to 1mb")
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
		log.Println(nm)

		if err == io.EOF {
			break
		}
	}
}
