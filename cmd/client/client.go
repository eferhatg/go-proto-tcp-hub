package main

import (
	"bufio"
	"fmt"
	"log"
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

		b := make([]byte, 10000)
		_, err := conn.Read(b)
		if err != nil {
			log.Printf(err.Error())
		}

		fmt.Print("Response: " + string(b))
	}
}
