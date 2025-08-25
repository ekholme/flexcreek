package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ekholme/flexcreek"
)

type workoutService struct {
	db *sql.DB
}

func NewWorkoutService(db *sql.DB) flexcreek.WorkoutService {
	return &workoutService{
		db: db,
	}
}

func (s *workoutService) CreateWorkout(ctx context.Context, w *flexcreek.Workout) (int, error) {
	// Start a new transaction. This is the "atomic" boundary for the operation.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback. If everything succeeds, we commit and the rollback is a no-op.
	// If any error occurs, the rollback will execute.
	defer tx.Rollback()

	// 1. Insert the parent Workout record within the transaction.
	qry := `INSERT INTO workouts (user_id, date, notes, duration) VALUES (?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, qry, w.UserID, w.Date, w.Notes, w.Duration)
	if err != nil {
		return 0, fmt.Errorf("failed to insert workout: %w", err)
	}
	workoutID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id for workout: %w", err)
	}

	// 2. Create a new movement instance service that operates on our transaction.
	// This is where the dependency injection happens at the method level.
	txMovementInstanceSvc := &movementInstanceService{db: tx}

	// 3. Loop through and create each MovementInstance using the transactional service.
	for _, mi := range w.MovementInstances {
		mi.WorkoutID = int(workoutID)
		if _, err := txMovementInstanceSvc.CreateMovementInstance(ctx, mi); err != nil {
			// The deferred Rollback will handle cleanup.
			return 0, fmt.Errorf("failed to create movement instance for workout: %w", err)
		}
	}

	// 4. If all operations were successful, commit the transaction.
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(workoutID), nil
}

func (s *workoutService) GetWorkoutByID(ctx context.Context, id int) (*flexcreek.Workout, error) {
	//todo
	return nil, nil
}

func (s *workoutService) GetAllWorkoutsByUser(ctx context.Context, user *flexcreek.User) ([]*flexcreek.Workout, error) {
	//todo
	return nil, nil
}

func (s *workoutService) GetWorkoutByDate(ctx context.Context, user *flexcreek.User, d time.Time) (*flexcreek.Workout, error) {
	//todo
	return nil, nil
}

func (s *workoutService) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	//todo
	return nil
}

func (s *workoutService) DeleteWorkout(ctx context.Context, id int) error {
	//todo
	return nil
}
