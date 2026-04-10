package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ekholme/flexcreek"
)

func (s *Storage) CreateWorkout(ctx context.Context, w *flexcreek.Workout) (int, error) {
	qry := `
		INSERT INTO workouts (
			user_id,
			short_description,
			long_description,
			workout_date	
		) 
		VALUES (?, ?, ?, ?)
	`

	res, err := s.db.ExecContext(ctx, qry, w.UserID, w.ShortDescription, w.LongDescription, w.WorkoutDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetWorkoutByID(ctx context.Context, id int, userID int) (*flexcreek.Workout, error) {
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

	if err := s.db.QueryRowContext(ctx, qry, id, userID).Scan(&w.ID, &w.UserID, &w.ShortDescription, &w.LongDescription, &w.WorkoutDate, &w.CreatedAt); err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Storage) GetWorkoutByDate(ctx context.Context, date time.Time, userID int) (*flexcreek.Workout, error) {
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

	if err := s.db.QueryRowContext(ctx, qry, formattedDate).Scan(&w.ID, &w.UserID, &w.ShortDescription, &w.LongDescription, &w.WorkoutDate, &w.CreatedAt); err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Storage) GetLatestWorkouts(ctx context.Context, n int, userID int) ([]*flexcreek.Workout, error) {
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

	rows, err := s.db.QueryContext(ctx, qry, userID, n)

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

func (s *Storage) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	qry := `
		UPDATE workouts
		SET short_description = ?,
		long_description = ?,
		workout_date = ?
		WHERE id = ?
		  AND user_id = ?
	`

	res, err := s.db.ExecContext(ctx, qry, w.ShortDescription, w.LongDescription, w.WorkoutDate, w.ID, w.UserID)

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

func (s *Storage) DeleteWorkout(ctx context.Context, id int) error {
	qry := `
		DELETE FROM workouts WHERE id = ?	
	`

	res, err := s.db.ExecContext(ctx, qry, id)

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
