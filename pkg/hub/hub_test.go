package hub

import (
	"net"
	"reflect"
	"testing"
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

func TestHub_Start(t *testing.T) {
	h := NewHub()
	go h.Start(9999)
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		t.Error(err.Error())
	}
	if conn == nil {
		t.Error("Nil connection error")
	}
	h.terminate <- true

}
