package hub

import (
	"bytes"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/eferhatg/uinty-assignment/pkg/protocol"
	"github.com/golang/protobuf/proto"
)

func TestHub_NewHub(t *testing.T) {
	h := NewHub()

	if h == nil {
		t.Error("NewHub init error")
	}

	if reflect.TypeOf(h).String() != "*hub.Hub" {
		t.Error("Wrong type error ", reflect.TypeOf(h).String(), "*hub.Hub")
	}
}

func TestHub_SimpleClientMsgSendRecieve(t *testing.T) {
	h := NewHub()
	go h.Start(8731)
	conn, err := net.Dial("tcp", ":8731")

	if err != nil {
		t.Error(err.Error())
	}
	if conn == nil {
		t.Error("Nil connection error")
	}
	defer conn.Close()

	m := protocol.Message{
		Command: protocol.Message_IDENTITY,
	}
	bt, _ := proto.Marshal(&m)

	if _, err := conn.Write(bt); err != nil {

		t.Error("Could not write to hub error:", err)
	}

	out := make([]byte, 1024)
	if _, err := conn.Read(out); err == nil {
		if bytes.Compare(out, bt) == 0 {
			t.Error("Response did match error")
		}
	} else {
		t.Error("Could not read error")
	}

}

func TestHub_1000ClientIdentity(t *testing.T) {

	clientCount := 1000
	/*
		Read handler function
	*/
	handler := func(conn net.Conn) {
		for {
			out := make([]byte, 1024)
			nn, err := conn.Read(out)
			if err != nil {
				conn.Close()
				if err != io.EOF {
					t.Error("Failed")
				}
			} else {
				log.Printf("Read %s bytes from client ", strconv.Itoa(nn))
			}
		}
	}

	/*
		Init New Hub
		Connect clients
	*/
	h := NewHub()
	go h.Start(4566)
	conns := make([]*net.Conn, clientCount)
	for i := 0; i < clientCount; i++ {
		conn, err := net.Dial("tcp", ":4566")
		if err != nil {
			t.Error("Failed to connect to tcp server on address ", err)
			conn.Close()

		}
		conns[i] = &conn
	}

	/*
		Write Identity message to socket
	*/
	for i := 0; i < clientCount; i++ {
		m := protocol.Message{
			Command: protocol.Message_IDENTITY,
		}
		bt, _ := proto.Marshal(&m)
		if _, err := (*conns[i]).Write(bt); err != nil {
			t.Error("Could not write to hub error:", err)
		}
	}

	/*
		Handle all connections read ops
	*/
	for i := 0; i < clientCount; i++ {
		go handler(*conns[i])
	}

}

// func TestHub_1000_Dial(t *testing.T) {
// 	h := NewHub()
// 	go h.Start(6789)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			conn, err := net.Dial("tcp", ":6789")

// 			if err != nil {
// 				t.Error(err.Error())
// 				return
// 			}
// 			if conn == nil {
// 				t.Error("Nil connection error")
// 				return
// 			}

// 			defer conn.Close()
// 		}()

// 	}

// }
