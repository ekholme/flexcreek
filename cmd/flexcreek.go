package main

import "github.com/ekholme/flexcreek/server"

const addr = ":8080"

func main() {
	s := server.NewServer(addr)
	s.Run()
}
