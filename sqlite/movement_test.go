package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ekholme/flexcreek"
	_ "modernc.org/sqlite" // Import the sqlite driver for database/sql
)

// createMovementTestSchema creates the necessary tables for the movement tests.
func createMovementTestSchema(t *testing.T, db *sql.DB) {
	t.Helper()
	// Only the movements table is strictly necessary for this service's tests.
	schema := `
	CREATE TABLE IF NOT EXISTS movements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		movement_type TEXT, 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
}

// newTestMovementService creates a new MovementService with a clean in-memory database.
// It returns the service and a teardown function to close the DB.
func newTestMovementService(t *testing.T) (flexcreek.MovementService, func()) {
	t.Helper()

	// mustOpenDB is defined in user_test.go but we can redefine it here
	// for package-level test isolation or move it to a test helper file.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	createMovementTestSchema(t, db)

	// Teardown function to clean up after tests.
	teardown := func() {
		db.Close()
	}

	return NewMovementService(db), teardown
}

func TestMovementService_CreateMovement(t *testing.T) {
	t.Parallel()
	service, teardown := newTestMovementService(t)
	defer teardown()

	ctx := context.Background()
	movement := &flexcreek.Movement{
		Name:         "Squat",
		MovementType: flexcreek.StrengthMovement,
	}

	id, err := service.CreateMovement(ctx, movement)
	if err != nil {
		t.Fatalf("CreateMovement() error = %v, want nil", err)
	}

	if id != 1 {
		t.Errorf("CreateMovement() id = %d, want 1", id)
	}

	// Verify movement was actually created by fetching it
	createdMovement, err := service.GetMovementByID(ctx, id)
	if err != nil {
		t.Fatalf("GetMovementByID() after create error = %v", err)
	}

	if createdMovement.Name != movement.Name {
		t.Errorf("created movement name = %s, want %s", createdMovement.Name, movement.Name)
	}
}

func TestMovementService_GetMovement(t *testing.T) {
	t.Parallel()
	service, teardown := newTestMovementService(t)
	defer teardown()

	ctx := context.Background()
	movement := &flexcreek.Movement{
		Name:         "Deadlift",
		MovementType: flexcreek.StrengthMovement,
	}

	id, err := service.CreateMovement(ctx, movement)
	if err != nil {
		t.Fatalf("failed to create movement for testing: %v", err)
	}

	t.Run("GetMovementByID", func(t *testing.T) {
		found, err := service.GetMovementByID(ctx, id)
		if err != nil {
			t.Fatalf("GetMovementByID() error = %v, want nil", err)
		}
		if found.ID != id || found.Name != movement.Name {
			t.Errorf("GetMovementByID() returned wrong data")
		}
	})

	t.Run("GetMovementByID_NotFound", func(t *testing.T) {
		_, err := service.GetMovementByID(ctx, 999)
		if err != sql.ErrNoRows {
			t.Errorf("GetMovementByID() with non-existent ID, error = %v, want %v", err, sql.ErrNoRows)
		}
	})

	t.Run("GetMovementByName", func(t *testing.T) {
		found, err := service.GetMovementByName(ctx, "Deadlift")
		if err != nil {
			t.Fatalf("GetMovementByName() error = %v, want nil", err)
		}
		if found.ID != id || found.Name != movement.Name {
			t.Errorf("GetMovementByName() returned wrong data")
		}
	})

	t.Run("GetMovementByName_NotFound", func(t *testing.T) {
		_, err := service.GetMovementByName(ctx, "NonExistent")
		if err != sql.ErrNoRows {
			t.Errorf("GetMovementByName() with non-existent name, error = %v, want %v", err, sql.ErrNoRows)
		}
	})
}

func TestMovementService_GetAllMovements(t *testing.T) {
	t.Parallel()
	service, teardown := newTestMovementService(t)
	defer teardown()

	ctx := context.Background()

	// Create two movements
	m1 := &flexcreek.Movement{Name: "Bench Press", MovementType: flexcreek.StrengthMovement}
	m2 := &flexcreek.Movement{Name: "Running", MovementType: flexcreek.CardioMovement}
	_, _ = service.CreateMovement(ctx, m1)
	_, _ = service.CreateMovement(ctx, m2)

	t.Run("GetAllMovements", func(t *testing.T) {
		movements, err := service.GetAllMovements(ctx)
		if err != nil {
			t.Fatalf("GetAllMovements() error = %v, want nil", err)
		}
		if len(movements) != 2 {
			t.Errorf("GetAllMovements() count = %d, want 2", len(movements))
		}
	})

	t.Run("GetAllMovementsByType", func(t *testing.T) {
		strengthMovements, err := service.GetAllMovementsByType(ctx, flexcreek.StrengthMovement)
		if err != nil {
			t.Fatalf("GetAllMovementsByType() error = %v, want nil", err)
		}
		if len(strengthMovements) != 1 {
			t.Errorf("GetAllMovementsByType('strength') count = %d, want 1", len(strengthMovements))
		}
		if strengthMovements[0].Name != "Bench Press" {
			t.Errorf("Incorrect movement returned for type 'strength'")
		}
	})
}

func TestMovementService_UpdateMovement(t *testing.T) {
	t.Parallel()
	service, teardown := newTestMovementService(t)
	defer teardown()

	ctx := context.Background()
	movement := &flexcreek.Movement{Name: "Pull-up", MovementType: flexcreek.StrengthMovement}
	id, err := service.CreateMovement(ctx, movement)
	if err != nil {
		t.Fatalf("failed to create movement for update test: %v", err)
	}

	updatedMovement := &flexcreek.Movement{ID: id, Name: "Weighted Pull-up", MovementType: flexcreek.StrengthMovement}
	err = service.UpdateMovement(ctx, updatedMovement)
	if err != nil {
		t.Fatalf("UpdateMovement() error = %v, want nil", err)
	}

	fetched, err := service.GetMovementByID(ctx, id)
	if err != nil {
		t.Fatalf("GetMovementByID() after update error = %v", err)
	}

	if fetched.Name != "Weighted Pull-up" {
		t.Errorf("movement was not updated correctly. Got name: %s", fetched.Name)
	}
}

func TestMovementService_DeleteMovement(t *testing.T) {
	t.Parallel()
	service, teardown := newTestMovementService(t)
	defer teardown()

	ctx := context.Background()
	movement := &flexcreek.Movement{Name: "To Be Deleted", MovementType: flexcreek.CardioMovement}
	id, err := service.CreateMovement(ctx, movement)
	if err != nil {
		t.Fatalf("failed to create movement for delete test: %v", err)
	}

	err = service.DeleteMovement(ctx, id)
	if err != nil {
		t.Fatalf("DeleteMovement() error = %v, want nil", err)
	}

	_, err = service.GetMovementByID(ctx, id)
	if err != sql.ErrNoRows {
		t.Errorf("GetMovementByID() after delete, error = %v, want %v", err, sql.ErrNoRows)
	}
}
