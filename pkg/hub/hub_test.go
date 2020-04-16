package hub

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/eferhatg/go-proto-tcp-hub/pkg/protocol"
	"github.com/golang/protobuf/proto"
)

const (
	minTCPPort         = 0
	maxTCPPort         = 65535
	maxReservedTCPPort = 1024
	maxRandTCPPort     = maxTCPPort - (maxReservedTCPPort + 1)
)

var (
	tcpPortRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func IsTCPPortAvailable(port int) bool {
	if port < minTCPPort || port > maxTCPPort {
		return false
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func FindAvailablePort() int {
	for i := maxReservedTCPPort; i < maxTCPPort; i++ {
		p := tcpPortRand.Intn(maxRandTCPPort) + maxReservedTCPPort + 1
		if IsTCPPortAvailable(p) {
			return p
		}
	}
	return -1
}

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
	port := FindAvailablePort()
	h := NewHub()
	go h.Start(port)
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))

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

func TestHub_MultiClientIdentity(t *testing.T) {

	clientCount := 5
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
				log.Printf("Read %s bytes from IDENTITY RESPONSE. Client: ", strconv.Itoa(nn))
			}
		}
	}

	/*
		Init New Hub
		Connect clients
	*/

	port := FindAvailablePort()
	h := NewHub()
	go h.Start(port)
	time.Sleep(time.Second)
	conns := make([]*net.Conn, clientCount)
	for i := 0; i < clientCount; i++ {
		conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
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

func TestHub_MultiClientList(t *testing.T) {

	clientCount := 5
	/*
		Read handler function
	*/
	handler := func(conn net.Conn) {
		for {
			out := make([]byte, 1024*1024)
			_, err := conn.Read(out)
			if err != nil {
				conn.Close()
				if err != io.EOF {
					t.Error("Failed")
				}
			}
		}
	}

	/*
		Init New Hub
		Connect clients
	*/

	port := FindAvailablePort()
	h := NewHub()
	go h.Start(port)
	time.Sleep(time.Second)
	conns := make([]*net.Conn, clientCount)
	for i := 0; i < clientCount; i++ {
		conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			t.Error("Failed to connect to tcp server on address ", err)
			conn.Close()

		}
		conns[i] = &conn
	}

	connSender, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		t.Error("Failed to connect to tcp server on address ", err)
		connSender.Close()

	}

	m := protocol.Message{
		Command: protocol.Message_LIST,
	}
	bt, _ := proto.Marshal(&m)
	if _, err := connSender.Write(bt); err != nil {
		t.Error("Could not write to hub error:", err)
	}

	/*
		Handle all connections read ops
	*/
	for i := 0; i < clientCount; i++ {
		go handler(*conns[i])
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	/*Building sender connection*/
	go func(conn net.Conn, waitgroup *sync.WaitGroup) {
		for {
			out := make([]byte, 1024*1024)
			len, err := conn.Read(out)
			if err != nil {
				conn.Close()
				if err != io.EOF {
					t.Error("Failed")
				}
			} else {
				m := &protocol.Message{}
				err := proto.Unmarshal(out[:len], m)
				if err != nil {
					t.Error(err.Error())
				}
				cIds := m.GetConnectedClientIds()
				sort.Slice(cIds, func(i, j int) bool { return cIds[i] < cIds[j] })

				expected := []uint64{1, 2, 3, 4, 5}
				if !reflect.DeepEqual(expected, cIds) {
					t.Error("Wrong connected ids")
				}
				waitgroup.Done()
			}
		}
	}(connSender, &waitgroup)
	waitgroup.Wait()

}

func TestHub_MultiClientRelay(t *testing.T) {

	clientCount := 5
	/*
		Read handler function
	*/
	handler := func(conn net.Conn, waitgroup *sync.WaitGroup) {
		for {
			out := make([]byte, 1024*1024)
			len, err := conn.Read(out)
			if err != nil {
				conn.Close()
				if err != io.EOF {
					t.Error("Failed")
				} else {
					m := &protocol.Message{}
					err := proto.Unmarshal(out[:len], m)
					if err != nil {
						t.Error(err.Error())
					}
					expected := "This can be anything up to 1mb"
					if string(m.GetBody()) != expected {
						t.Error("Wrong body")
					}
				}
			}
			waitgroup.Done()
		}
	}

	/*
		Init New Hub
		Connect clients
	*/

	port := FindAvailablePort()
	h := NewHub()
	go h.Start(port)
	time.Sleep(time.Second)
	conns := make([]*net.Conn, clientCount)
	for i := 0; i < clientCount; i++ {
		conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			t.Error("Failed to connect to tcp server on address ", err)
			conn.Close()

		}
		conns[i] = &conn
	}

	connSender, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		t.Error("Failed to connect to tcp server on address ", err)
		connSender.Close()

	}

	to := []uint64{1, 2, 3, 4, 5, 6}
	bodyType := protocol.Message_PLAIN_TEXT
	body := []byte("This can be anything up to 1mb")
	m := protocol.Message{
		Command:  protocol.Message_RELAY,
		Id:       6,
		RelayTo:  to,
		BodyType: bodyType,
		Body:     body,
	}

	bt, _ := proto.Marshal(&m)
	if _, err := connSender.Write(bt); err != nil {
		t.Error("Could not write to hub error:", err)
	}

	/*
		Handle all connections read ops
	*/
	var waitgroup sync.WaitGroup
	waitgroup.Add(clientCount)
	for i := 0; i < clientCount; i++ {
		go handler(*conns[i], &waitgroup)
	}

	waitgroup.Wait()
	connSender.Close()

}
