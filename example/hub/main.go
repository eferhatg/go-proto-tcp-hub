package main

import (
	"fmt"

	"github.com/eferhatg/uinty-assignment/pkg/hub"
)

func main() {

	fmt.Println("Listening port:1087 ")
	h := hub.NewHub()
	h.Start(1087)
}
