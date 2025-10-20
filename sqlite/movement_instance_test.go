package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ekholme/flexcreek"
	_ "modernc.org/sqlite" // Import the sqlite driver for database/sql
)

// newTestMovementInstanceService creates a new MovementInstanceService with a clean in-memory database.
func newTestMovementInstanceService(t *testing.T) (flexcreek.MovementInstanceService, *sql.DB, func()) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	createAllTablesTestSchema(t, db)

	teardown := func() {
		db.Close()
	}

	return NewMovementInstanceService(db), db, teardown
}

// mustSeedDB is a helper to pre-populate the database with necessary records for tests.
func mustSeedDB(t *testing.T, db *sql.DB) (userID, movementID, workoutID int) {
	t.Helper()

	// Create a user
	res, err := db.Exec(`INSERT INTO users (email) VALUES ('test@user.com')`)
	if err != nil {
		t.Fatal(err)
	}
	uid, _ := res.LastInsertId()

	// Create a movement
	res, err = db.Exec(`INSERT INTO movements (name, movement_type) VALUES ('Squat', 'strength')`)
	if err != nil {
		t.Fatal(err)
	}
	mid, _ := res.LastInsertId()

	// Create a workout
	res, err = db.Exec(`INSERT INTO workouts (user_id, workout_date) VALUES (?, ?)`, uid, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	wid, _ := res.LastInsertId()

	return int(uid), int(mid), int(wid)
}

func TestMovementInstanceService_CreateAndGet(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestMovementInstanceService(t)
	defer teardown()

	_, movementID, workoutID := mustSeedDB(t, db)
	ctx := context.Background()

	mi := &flexcreek.MovementInstance{
		WorkoutID: workoutID,
		Movement:  &flexcreek.Movement{ID: movementID, MovementType: flexcreek.StrengthMovement},
		Notes:     "Felt strong",
		Strength: &flexcreek.StrengthLog{
			Sets: []flexcreek.Set{{Reps: 5, Weight: 100}},
		},
	}

	// Test Create
	id, err := service.CreateMovementInstance(ctx, mi)
	if err != nil {
		t.Fatalf("CreateMovementInstance() error = %v, want nil", err)
	}
	if id == 0 {
		t.Fatal("CreateMovementInstance() returned id 0")
	}

	// Test Get by ID
	fetched, err := service.GetMovementInstanceByID(ctx, id)
	if err != nil {
		t.Fatalf("GetMovementInstanceByID() error = %v, want nil", err)
	}

	if fetched.ID != id {
		t.Errorf("ID mismatch: got %d, want %d", fetched.ID, id)
	}
	if fetched.Notes != mi.Notes {
		t.Errorf("Notes mismatch: got %s, want %s", fetched.Notes, mi.Notes)
	}
	if fetched.Movement == nil || fetched.Movement.ID != movementID {
		t.Errorf("Movement ID mismatch")
	}
	if fetched.Strength == nil || len(fetched.Strength.Sets) != 1 || fetched.Strength.Sets[0].Weight != 100 {
		t.Errorf("Strength log data was not unmarshaled correctly. Got: %+v", fetched.Strength)
	}

	// Test Get Not Found
	_, err = service.GetMovementInstanceByID(ctx, 999)
	if err != sql.ErrNoRows {
		t.Errorf("GetMovementInstanceByID() with non-existent ID, error = %v, want %v", err, sql.ErrNoRows)
	}
}

func TestMovementInstanceService_GetAllByWorkoutID(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestMovementInstanceService(t)
	defer teardown()

	_, movementID, workoutID := mustSeedDB(t, db)
	ctx := context.Background()

	// Create two instances for the same workout
	mi1 := &flexcreek.MovementInstance{WorkoutID: workoutID, Movement: &flexcreek.Movement{ID: movementID}}
	mi2 := &flexcreek.MovementInstance{WorkoutID: workoutID, Movement: &flexcreek.Movement{ID: movementID}}
	_, _ = service.CreateMovementInstance(ctx, mi1)
	_, _ = service.CreateMovementInstance(ctx, mi2)

	// Create another workout and instance to ensure we don't fetch it
	res, err := db.Exec(`INSERT INTO workouts (user_id) VALUES (1)`)
	if err != nil {
		t.Fatal(err)
	}
	otherWorkoutID, _ := res.LastInsertId()
	mi3 := &flexcreek.MovementInstance{WorkoutID: int(otherWorkoutID), Movement: &flexcreek.Movement{ID: movementID}}
	_, _ = service.CreateMovementInstance(ctx, mi3)

	instances, err := service.GetAllMovementInstancesByWorkoutID(ctx, workoutID)
	if err != nil {
		t.Fatalf("GetAllMovementInstancesByWorkoutID() error = %v", err)
	}

	if len(instances) != 2 {
		t.Errorf("expected 2 movement instances, got %d", len(instances))
	}
}

func TestMovementInstanceService_GetAllForMovement(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestMovementInstanceService(t)
	defer teardown()

	userID, squatID, workoutID1 := mustSeedDB(t, db)
	ctx := context.Background()

	// Create another movement (Bench)
	res, err := db.Exec(`INSERT INTO movements (name, movement_type) VALUES ('Bench Press', 'strength')`)
	if err != nil {
		t.Fatal(err)
	}
	benchID, _ := res.LastInsertId()

	// Create another workout
	res, err = db.Exec(`INSERT INTO workouts (user_id) VALUES (?)`, userID)
	if err != nil {
		t.Fatal(err)
	}
	workoutID2, _ := res.LastInsertId()

	// Instance 1: Squat in workout 1
	_, _ = service.CreateMovementInstance(ctx, &flexcreek.MovementInstance{WorkoutID: workoutID1, Movement: &flexcreek.Movement{ID: squatID}})
	// Instance 2: Bench in workout 1
	_, _ = service.CreateMovementInstance(ctx, &flexcreek.MovementInstance{WorkoutID: workoutID1, Movement: &flexcreek.Movement{ID: int(benchID)}})
	// Instance 3: Squat in workout 2
	_, _ = service.CreateMovementInstance(ctx, &flexcreek.MovementInstance{WorkoutID: int(workoutID2), Movement: &flexcreek.Movement{ID: squatID}})

	// Fetch all instances of Squat for the user
	instances, err := service.GetAllMovementInstancesForMovement(ctx, userID, squatID)
	if err != nil {
		t.Fatalf("GetAllMovementInstancesForMovement() error = %v", err)
	}

	if len(instances) != 2 {
		t.Errorf("expected 2 squat instances, got %d", len(instances))
	}
	for _, mi := range instances {
		if mi.Movement.ID != squatID {
			t.Error("fetched instance for the wrong movement")
		}
	}
}

func TestMovementInstanceService_Update(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestMovementInstanceService(t)
	defer teardown()

	_, movementID, workoutID := mustSeedDB(t, db)
	ctx := context.Background()

	rpe := 5
	mi := &flexcreek.MovementInstance{
		WorkoutID: workoutID,
		Movement:  &flexcreek.Movement{ID: movementID},
		Notes:     "Initial notes",
		RPE:       &rpe,
	}

	id, err := service.CreateMovementInstance(ctx, mi)
	if err != nil {
		t.Fatal(err)
	}

	updatedRPE := 7
	updatedMI := &flexcreek.MovementInstance{
		ID:        id,
		WorkoutID: workoutID,
		Movement:  &flexcreek.Movement{ID: movementID},
		Notes:     "Updated notes",
		RPE:       &updatedRPE,
	}

	err = service.UpdateMovementInstance(ctx, updatedMI)
	if err != nil {
		t.Fatalf("UpdateMovementInstance() error = %v", err)
	}

	fetched, err := service.GetMovementInstanceByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if fetched.Notes != "Updated notes" {
		t.Errorf("Notes were not updated. Got: %s", fetched.Notes)
	}
	if fetched.RPE == nil || *fetched.RPE != 7 {
		t.Errorf("RPE was not updated. Got: %v", fetched.RPE)
	}
}

func TestMovementInstanceService_Delete(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestMovementInstanceService(t)
	defer teardown()

	_, movementID, workoutID := mustSeedDB(t, db)
	ctx := context.Background()

	mi := &flexcreek.MovementInstance{WorkoutID: workoutID, Movement: &flexcreek.Movement{ID: movementID}}
	id, err := service.CreateMovementInstance(ctx, mi)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteMovementInstance(ctx, id)
	if err != nil {
		t.Fatalf("DeleteMovementInstance() error = %v", err)
	}

	_, err = service.GetMovementInstanceByID(ctx, id)
	if err != sql.ErrNoRows {
		t.Errorf("GetMovementInstanceByID() after delete, error = %v, want %v", err, sql.ErrNoRows)
	}
}
