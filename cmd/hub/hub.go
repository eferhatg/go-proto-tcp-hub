package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func main() {

	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":5555")

	conn, _ := ln.Accept()

	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Print("Got:", string(message))

		newmessage := strings.ToUpper(message)

		conn.Write([]byte(newmessage + "\n"))
	}
}
