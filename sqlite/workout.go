package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/ekholme/flexcreek"
)

type workoutService struct {
	db *sql.DB
}

func NewWorkoutService(db *sql.DB) flexcreek.WorkoutService {
	return workoutService{
		db: db,
	}
}

func (ws workoutService) CreateWorkout(ctx context.Context, w *flexcreek.Workout) (int, error) {
	qry := `
		INSERT INTO workouts (
			user_id,
			short_description,
			long_description,
			workout_date	
		) 
		VALUES (?, ?, ?, ?)
	`

	res, err := ws.db.ExecContext(ctx, qry, w.UserID, w.ShortDescription, w.LongDescription, w.WorkoutDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (ws workoutService) GetWorkoutByID(ctx context.Context, id int, userID int) (*flexcreek.Workout, error) {
	qry := `
		SELECT id,
		user_id,
		short_description,
		long_description,
		workout_date,
		created_at
		FROM workouts
		WHERE id = ?
		  AND user_id = ?
	`

	var w flexcreek.Workout

	if err := ws.db.QueryRowContext(ctx, qry, id, userID).Scan(&w.ID, &w.UserID, &w.ShortDescription, &w.LongDescription, &w.WorkoutDate, &w.CreatedAt); err != nil {
		return nil, err
	}

	return &w, nil
}

func (ws workoutService) GetWorkoutByDate(ctx context.Context, date time.Time, userID int) (*flexcreek.Workout, error) {
	qry := `
		SELECT id,
		user_id,
		short_description,
		long_description,
		workout_date,
		created_at
		FROM workouts
		WHERE id = ?
		  AND date = ?
	`

	formattedDate := date.Format("2026-04-01")

	var w flexcreek.Workout

	if err := ws.db.QueryRowContext(ctx, qry, formattedDate).Scan(&w.ID, &w.UserID, &w.ShortDescription, &w.LongDescription, &w.WorkoutDate, &w.CreatedAt); err != nil {
		return nil, err
	}

	return &w, nil
}

func (ws workoutService) GetLatestWorkouts(ctx context.Context, n int, userID int) ([]*flexcreek.Workout, error) {
	return nil, nil
}

func (ws workoutService) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	return nil
}

func (ws workoutService) DeleteWorkout(ctx context.Context, id int) error {
	return nil
}
