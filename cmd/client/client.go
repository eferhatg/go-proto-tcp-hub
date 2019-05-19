package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, _ := net.Dial("tcp", ":5555")
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Write text: ")
		text, _ := reader.ReadString('\n')

		connw := bufio.NewWriter(conn)
		connw.WriteString(text + "\n")
		connw.Flush()

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Response: " + message)
	}
}
