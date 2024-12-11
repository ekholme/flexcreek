package main

import (
	"fmt"
	"log"

	"github.com/ekholme/flexcreek/server"
)

const addr = ":8080"

func main() {
	s := server.NewServer(addr)

	fmt.Println("Running Flex Creek on port 8080")

	err := s.Run()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
