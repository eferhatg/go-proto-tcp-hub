package client

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	c := NewClient(nil)

	if c == nil {
		t.Error("Client init error")
	}

	if reflect.TypeOf(c).String() != "*client.Client" {
		t.Error("Wrong type error ", reflect.TypeOf(c).String(), "*client.Client")
	}
}

func TestClient_Read(t *testing.T) {

	fi, err := os.Open("testdata/readtest.dat")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(fi)

	c := NewClient(nil)
	c.Reader = r

	b, err := c.Read()
	expected := "readtest"
	if string(b) != expected {
		t.Error("Read error")
	}
}

func TestClient_Write(t *testing.T) {

	fo, err := os.Create("testdata/writetest.dat")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)

	c := NewClient(nil)
	c.Writer = w
	expected := "writetest"
	c.Write([]byte(expected))

	fi, err := os.Open("testdata/writetest.dat")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	len, err := r.Read(buf)
	if err != nil {
		log.Printf(err.Error())

	}

	if string(buf[:len]) != expected {
		t.Error("Read error")
	}
}
