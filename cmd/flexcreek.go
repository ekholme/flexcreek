package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ekholme/flexcreek/sqlite"
	_ "modernc.org/sqlite"
)

const (
	dsn = "flexcreek.db"
)

func main() {
	fmt.Println("Hello from flexcreek!")

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("Couldn't open the database: %s", err)
	}

	storage := sqlite.NewStorage(db)

}
