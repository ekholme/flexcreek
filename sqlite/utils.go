package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func CreateTestDb(t *testing.T, sqlScript string) (*sql.DB, func()) {
	t.Helper() //marks this function as a test helper

	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatalf("Failed to open in-memory db for testing: %v", err)
	}

	//read in sql file to create tables for testing
	schema, err := os.ReadFile(sqlScript)

	if err != nil {
		t.Fatalf("Error reading sqlScript file: %v", err)
	}

	//execute the .sql creation script
	_, err = db.Exec(string(schema))

	if err != nil {
		t.Fatalf("Error creating tables %v", err)
	}

	cleanup := func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing test database:", err) //use fmt.Println rather than t.Fatalf so test doesn't fail
		}
	}

	return db, cleanup
}
