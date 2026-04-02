package sqlite

import (
	"context"
	"database/sql"
	"errors"
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
		WHERE user_id = ?
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
	qry := `
		SELECT id,
		user_id,
		short_description,
		long_description,
		workout_date,
		created_at
		FROM workouts
		WHERE user_id = ?	
		ORDER BY workout_date desc
		LIMIT ?;
	`

	rows, err := ws.db.QueryContext(ctx, qry, userID, n)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workouts []*flexcreek.Workout

	for rows.Next() {
		var w flexcreek.Workout

		err = rows.Scan(&w.ID, &w.UserID, &w.ShortDescription, &w.LongDescription, &w.WorkoutDate, &w.CreatedAt)

		if err != nil {
			return nil, err
		}

		workouts = append(workouts, &w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ws workoutService) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	qry := `
		UPDATE workouts
		SET short_description = ?,
		long_description = ?,
		workout_date = ?
		WHERE id = ?
		  AND user_id = ?
	`

	res, err := ws.db.ExecContext(ctx, qry, w.ShortDescription, w.LongDescription, w.WorkoutDate, w.ID, w.UserID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("no workout found with that ID for this user")
	}

	return nil
}

func (ws workoutService) DeleteWorkout(ctx context.Context, id int) error {
	qry := `
		DELETE FROM workouts WHERE id = ?	
	`

	res, err := ws.db.ExecContext(ctx, qry, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
