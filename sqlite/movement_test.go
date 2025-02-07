package sqlite_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ekholme/flexcreek"
	"github.com/ekholme/flexcreek/sqlite"

	_ "modernc.org/sqlite"
)

func TestMovementService_CreateMovement(t *testing.T) {
	//use memory for testing
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	//read in sql file to create tables for testing
	schema, err := os.ReadFile("migration.sql")

	if err != nil {
		t.Fatalf("Error reading migration.sql file: %v", err)
	}

	//execute the .sql creation script
	_, err = db.Exec(string(schema))

	if err != nil {
		t.Fatalf("Error creating tables %v", err)
	}

	//create a new movementservice
	ms := sqlite.NewMovementService(db)

	//just writing one for now, but eventually this will be more
	testCase := struct {
		name        string
		movement    flexcreek.Movement
		expectedID  int
		expectedErr string
	}{
		name: "Valid Movement Creation",
		movement: flexcreek.Movement{
			Name:    "Squat",
			Muscles: []string{"quads", "hamstrings", "glutes"},
		},
		expectedID:  1,
		expectedErr: "",
	}

	t.Run(testCase.name, func(t *testing.T) {
		id, err := ms.CreateMovement(&testCase.movement)

		if id != testCase.expectedID {
			t.Errorf("Expected ID: %d, got %d", testCase.expectedID, id)
		}

		if err != nil {
			t.Fatal(err)
		}
	})

	//continue to add in more testing, include testing to make sure the muscles are being written correctly
}
