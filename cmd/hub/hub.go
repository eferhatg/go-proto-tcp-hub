package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/eferhatg/uinty-assignment/pkg/client"
)

// only needed below for sample processing

func main() {

	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":5555")

	conn, _ := ln.Accept()
	c := client.NewClient(conn)

	for {

		b, err := c.Read()
		message := string(b)
		fmt.Print("Got:", message)

		newmessage := strings.ToUpper(message)

		c.Write([]byte(newmessage + "\n"))

		if err == io.EOF {
			break
		}

	}
}
