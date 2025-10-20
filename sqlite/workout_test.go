package sqlite

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/ekholme/flexcreek"
	_ "modernc.org/sqlite" // Import the sqlite driver for database/sql
)

// newTestWorkoutService creates a new WorkoutService with a clean in-memory database.
func newTestWorkoutService(t *testing.T) (flexcreek.WorkoutService, *sql.DB, func()) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	// Enable foreign key support in SQLite. It's off by default.
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	createAllTablesTestSchema(t, db)

	teardown := func() {
		db.Close()
	}

	return NewWorkoutService(db), db, teardown
}

// mustSeedWorkoutPrereqs pre-populates the DB with a user and movements.
func mustSeedWorkoutPrereqs(t *testing.T, db *sql.DB) (userID, squatID, benchID int) {
	t.Helper()

	res, err := db.Exec(`INSERT INTO users (email) VALUES ('test@user.com')`)
	if err != nil {
		t.Fatal(err)
	}
	uid, _ := res.LastInsertId()

	res, err = db.Exec(`INSERT INTO movements (name, movement_type) VALUES ('Squat', 'strength')`)
	if err != nil {
		t.Fatal(err)
	}
	sid, _ := res.LastInsertId()

	res, err = db.Exec(`INSERT INTO movements (name, movement_type) VALUES ('Bench Press', 'strength')`)
	if err != nil {
		t.Fatal(err)
	}
	bid, _ := res.LastInsertId()

	return int(uid), int(sid), int(bid)
}

func TestWorkoutService_CreateAndGetWorkout(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, squatID, benchID := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	workout := &flexcreek.Workout{
		UserID:   userID,
		Date:     time.Now(),
		Notes:    "A good session",
		Duration: 3600 * time.Second,
		MovementInstances: []*flexcreek.MovementInstance{
			{
				Movement: &flexcreek.Movement{ID: squatID, MovementType: flexcreek.StrengthMovement},
				Strength: &flexcreek.StrengthLog{Sets: []flexcreek.Set{{Reps: 5, Weight: 225}}},
			},
			{
				Movement: &flexcreek.Movement{ID: benchID, MovementType: flexcreek.StrengthMovement},
				Strength: &flexcreek.StrengthLog{Sets: []flexcreek.Set{{Reps: 8, Weight: 135}}},
			},
		},
	}

	// Test Create
	id, err := service.CreateWorkout(ctx, workout)
	if err != nil {
		t.Fatalf("CreateWorkout() error = %v, want nil", err)
	}
	if id == 0 {
		t.Fatal("CreateWorkout() returned id 0")
	}

	// Test Get by ID
	fetched, err := service.GetWorkoutByID(ctx, id)
	if err != nil {
		t.Fatalf("GetWorkoutByID() error = %v, want nil", err)
	}

	if fetched.ID != id {
		t.Errorf("ID mismatch: got %d, want %d", fetched.ID, id)
	}
	if fetched.Notes != workout.Notes {
		t.Errorf("Notes mismatch: got %s, want %s", fetched.Notes, workout.Notes)
	}
	if len(fetched.MovementInstances) != 2 {
		t.Fatalf("expected 2 movement instances, got %d", len(fetched.MovementInstances))
	}
	if fetched.MovementInstances[0].Movement.ID != squatID {
		t.Error("first movement instance is not Squat")
	}
	if fetched.MovementInstances[1].Movement.ID != benchID {
		t.Error("second movement instance is not Bench Press")
	}

	// Test Get Not Found
	_, err = service.GetWorkoutByID(ctx, 999)
	if err != sql.ErrNoRows {
		t.Errorf("GetWorkoutByID() with non-existent ID, error = %v, want %v", err, sql.ErrNoRows)
	}
}

func TestWorkoutService_GetAllWorkoutsByUser(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, squatID, _ := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	// Create two workouts for the user
	w1 := &flexcreek.Workout{UserID: userID, MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}}}
	w2 := &flexcreek.Workout{UserID: userID, MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}}}
	_, _ = service.CreateWorkout(ctx, w1)
	_, _ = service.CreateWorkout(ctx, w2)

	// Create another user and workout that should not be fetched
	res, _ := db.Exec(`INSERT INTO users (email) VALUES ('other@user.com')`)
	otherUserID, _ := res.LastInsertId()
	w3 := &flexcreek.Workout{UserID: int(otherUserID), MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}}}
	_, _ = service.CreateWorkout(ctx, w3)

	workouts, err := service.GetAllWorkoutsByUser(ctx, userID)
	if err != nil {
		t.Fatalf("GetAllWorkoutsByUser() error = %v", err)
	}

	if len(workouts) != 2 {
		t.Errorf("expected 2 workouts for user, got %d", len(workouts))
	}
}

func TestWorkoutService_GetWorkoutsByDate(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, squatID, _ := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	today := time.Now().Truncate(24 * time.Hour) // Ensure 'today' is at midnight for consistent comparison
	yesterday := today.AddDate(0, 0, -1)

	wToday := &flexcreek.Workout{UserID: userID, Date: today, MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}}}
	wYesterday := &flexcreek.Workout{UserID: userID, Date: yesterday, MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}}}
	_, _ = service.CreateWorkout(ctx, wToday)
	_, _ = service.CreateWorkout(ctx, wYesterday)

	// Test GetWorkoutsByDate
	workouts, err := service.GetWorkoutsByDate(ctx, userID, today)
	if err != nil {
		t.Fatalf("GetWorkoutsByDate() error = %v", err)
	}
	if len(workouts) != 1 {
		t.Errorf("expected 1 workout for today, got %d", len(workouts))
	}

	// Test GetWorkoutsByDateRange
	twoDaysAgo := today.AddDate(0, 0, -2)
	workouts, err = service.GetWorkoutsByDateRange(ctx, userID, twoDaysAgo, today)
	if err != nil {
		t.Fatalf("GetWorkoutsByDateRange() error = %v", err)
	}
	if len(workouts) != 2 {
		t.Errorf("expected 2 workouts in the date range, got %d", len(workouts))
	}
}

func TestWorkoutService_UpdateWorkout(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, squatID, benchID := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	// Create initial workout with one movement
	initialWorkout := &flexcreek.Workout{
		UserID:            userID,
		Notes:             "Initial Notes",
		MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}},
	}
	id, err := service.CreateWorkout(ctx, initialWorkout)
	if err != nil {
		t.Fatal(err)
	}

	// Update: change notes and replace movement instances
	updatedWorkout := &flexcreek.Workout{
		ID:     id,
		UserID: userID,
		Notes:  "Updated Notes",
		MovementInstances: []*flexcreek.MovementInstance{
			{Movement: &flexcreek.Movement{ID: benchID}}, // Replaced Squat with Bench
		},
	}

	err = service.UpdateWorkout(ctx, updatedWorkout)
	if err != nil {
		t.Fatalf("UpdateWorkout() error = %v", err)
	}

	fetched, err := service.GetWorkoutByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if fetched.Notes != "Updated Notes" {
		t.Errorf("Notes were not updated. Got: %s", fetched.Notes)
	}
	if len(fetched.MovementInstances) != 1 {
		t.Fatalf("Expected 1 movement instance after update, got %d", len(fetched.MovementInstances))
	}
	if fetched.MovementInstances[0].Movement.ID != benchID {
		t.Error("Movement instance was not updated correctly")
	}
}

func TestWorkoutService_DeleteWorkout(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, squatID, _ := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	workout := &flexcreek.Workout{
		UserID:            userID,
		MovementInstances: []*flexcreek.MovementInstance{{Movement: &flexcreek.Movement{ID: squatID}}},
	}
	id, err := service.CreateWorkout(ctx, workout)
	if err != nil {
		t.Fatal(err)
	}

	// Check that one movement instance was created
	var count int
	_ = db.QueryRowContext(ctx, `SELECT count(*) FROM movement_instances WHERE workout_id = ?`, id).Scan(&count)
	if count != 1 {
		t.Fatalf("expected 1 movement instance before delete, got %d", count)
	}

	// Delete the workout
	err = service.DeleteWorkout(ctx, id)
	if err != nil {
		t.Fatalf("DeleteWorkout() error = %v", err)
	}

	// Verify workout is gone
	_, err = service.GetWorkoutByID(ctx, id)
	if err != sql.ErrNoRows {
		t.Errorf("GetWorkoutByID() after delete, error = %v, want %v", err, sql.ErrNoRows)
	}

	// Verify movement instances are also gone (due to ON DELETE CASCADE)
	_ = db.QueryRowContext(ctx, `SELECT count(*) FROM movement_instances WHERE workout_id = ?`, id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 movement instances after delete, got %d", count)
	}
}

func TestWorkoutService_GetWorkoutByID_EmptyWorkout(t *testing.T) {
	t.Parallel()
	service, db, teardown := newTestWorkoutService(t)
	defer teardown()

	userID, _, _ := mustSeedWorkoutPrereqs(t, db)
	ctx := context.Background()

	// Create a workout with NO movement instances
	workout := &flexcreek.Workout{
		UserID: userID,
		Notes:  "Just notes, no movements",
	}

	id, err := service.CreateWorkout(ctx, workout)
	if err != nil {
		t.Fatalf("CreateWorkout() error = %v, want nil", err)
	}

	fetched, err := service.GetWorkoutByID(ctx, id)
	if err != nil {
		t.Fatalf("GetWorkoutByID() error = %v, want nil", err)
	}

	if fetched.ID != id {
		t.Errorf("ID mismatch: got %d, want %d", fetched.ID, id)
	}
	if len(fetched.MovementInstances) != 0 {
		t.Errorf("expected 0 movement instances for empty workout, got %d", len(fetched.MovementInstances))
	}
	if !reflect.DeepEqual(fetched.MovementInstances, workout.MovementInstances) {
		// A nil slice is not equal to an empty slice, but functionally they are the same for us.
		// The service initializes an empty slice, not a nil one.
		if !(workout.MovementInstances == nil && len(fetched.MovementInstances) == 0) {
			t.Errorf("MovementInstances mismatch: got %+v, want %+v", fetched.MovementInstances, workout.MovementInstances)
		}
	}
}
