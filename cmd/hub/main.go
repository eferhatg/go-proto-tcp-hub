package main

import (
	"fmt"

	"github.com/eferhatg/uinty-assignment/pkg/hub"
)

func main() {

	fmt.Println("Launching server...")

	h := hub.NewHub()
	h.Listen(1087)
}
