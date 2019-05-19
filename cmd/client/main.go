package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/eferhatg/uinty-assignment/pkg/client"
)

func main() {

	conn, _ := net.Dial("tcp", ":1087")
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Write text: ")
		text, _ := reader.ReadString('\n')
		//Writing to hub
		c := client.NewClient(conn, conn, conn)
		c.Write([]byte(text))

		//Reading hub response
		b, err := c.Read()
		if err != nil {
			log.Printf(err.Error())
		}

		fmt.Print("Response: " + string(b) + "\n")
	}
}
