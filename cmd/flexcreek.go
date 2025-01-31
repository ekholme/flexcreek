package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ekholme/flexcreek/server"
	"github.com/ekholme/flexcreek/sqlite"

	_ "modernc.org/sqlite"
)

const addr = ":8080"
const dsn = "./flexcreek.db"

func main() {
	//do this more elegantly later
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	ms := sqlite.NewMovementService(db)

	s := server.NewServer(addr, ms)

	fmt.Println("Running Flex Creek on port 8080")

	err = s.Run()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
