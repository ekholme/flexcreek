package sqlite_test

import (
	"testing"

	"github.com/ekholme/flexcreek"
	"github.com/ekholme/flexcreek/sqlite"

	_ "modernc.org/sqlite"
)

func TestMovementService_CreateMovement(t *testing.T) {
	sqlFile := "migration.sql"

	db, cleanup := sqlite.CreateTestDb(t, sqlFile)

	defer cleanup() //closes database after the test

	mus := sqlite.NewMuscleService(db)

	//create a new movementservice
	ms := sqlite.NewMovementService(db, mus)

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
