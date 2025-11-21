package sqlite

import (
	"context"
	"database/sql"

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

//methods ------------

func (ws *workoutService) CreateWorkout(ctx context.Context, w *flexcreek.Workout) (int, error) {
	qry := `
	INSERT INTO workouts (user_id, activity_type_id, duration_minutes, distance_miles, workout_details, workout_date)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := ws.db.ExecContext(ctx, qry, w.UserID, w.ActivityTypeID, w.DurationMins, w.DistanceMiles, w.WorkoutDetails, w.WorkoutDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (ws *workoutService) GetWorkoutByID(ctx context.Context, id int) (*flexcreek.Workout, error) {
	qry := `
	SELECT
		w.id,
		w.user_id,
		w.duration_minutes,
		w.distance_miles,
		w.workout_details,
		w.workout_date,
		w.created_at,
		w.updated_at,
		w.activity_type_id,
		a.id,
		a.name
	FROM
		workouts w
	JOIN
		activity_types a ON w.activity_type_id = a.id
	WHERE
		w.id = ?
	`

	var w flexcreek.Workout

	err := ws.db.QueryRowContext(ctx, qry, id).Scan(
		&w.ID,
		&w.UserID,
		&w.DurationMins,
		&w.DistanceMiles,
		&w.WorkoutDetails,
		&w.WorkoutDate,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.ActivityTypeID,
		&w.ActivityType.ID,
		&w.ActivityType.Name,
	)

	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (ws *workoutService) LatestWorkouts(ctx context.Context, userID, n int) ([]*flexcreek.Workout, error) {
	qry := `
	SELECT
		w.id,
		w.user_id,
		w.duration_minutes,
		w.distance_miles,
		w.workout_details,
		w.workout_date,
		w.created_at,
		w.updated_at,
		w.activity_type_id,
		a.id,
		a.name
	FROM
		workouts w
	JOIN
		activity_types a ON w.activity_type_id = a.id
	WHERE
		w.user_id = ?
	ORDER BY workout_date DESC
	LIMIT ?;
	`
	rows, err := ws.db.QueryContext(ctx, qry, userID, n)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workouts := make([]*flexcreek.Workout, 0, n)

	for rows.Next() {
		var w flexcreek.Workout

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.DurationMins,
			&w.DistanceMiles,
			&w.WorkoutDetails,
			&w.WorkoutDate,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.ActivityTypeID,
			&w.ActivityType.ID,
			&w.ActivityType.Name,
		); err != nil {
			return nil, err
		}

		workouts = append(workouts, &w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ws *workoutService) GetWorkoutsByUser(ctx context.Context, userID int) ([]*flexcreek.Workout, error) {
	qry := `
	SELECT
		w.id,
		w.user_id,
		w.duration_minutes,
		w.distance_miles,
		w.workout_details,
		w.workout_date,
		w.created_at,
		w.updated_at,
		w.activity_type_id,
		a.id,
		a.name
	FROM
		workouts w
	JOIN
		activity_types a ON w.activity_type_id = a.id
	WHERE
		w.user_id = ?
	`

	rows, err := ws.db.QueryContext(ctx, qry, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workouts := make([]*flexcreek.Workout, 0)

	for rows.Next() {
		var w flexcreek.Workout

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.DurationMins,
			&w.DistanceMiles,
			&w.WorkoutDetails,
			&w.WorkoutDate,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.ActivityTypeID,
			&w.ActivityType.ID,
			&w.ActivityType.Name,
		); err != nil {
			return nil, err
		}

		workouts = append(workouts, &w)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ws *workoutService) GetWorkoutsByActivityType(ctx context.Context, userID int, activityTypeID int) ([]*flexcreek.Workout, error) {
	qry := `
	SELECT
		w.id,
		w.user_id,
		w.duration_minutes,
		w.distance_miles,
		w.workout_details,
		w.workout_date,
		w.created_at,
		w.updated_at,
		w.activity_type_id,
		a.id,
		a.name
	FROM
		workouts w
	JOIN
		activity_types a ON w.activity_type_id = a.id
	WHERE
		w.user_id = ?
		AND w.activity_type_id = ?
	`

	rows, err := ws.db.QueryContext(ctx, qry, userID, activityTypeID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workouts := make([]*flexcreek.Workout, 0)

	for rows.Next() {
		var w flexcreek.Workout

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.DurationMins,
			&w.DistanceMiles,
			&w.WorkoutDetails,
			&w.WorkoutDate,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.ActivityTypeID,
			&w.ActivityType.ID,
			&w.ActivityType.Name,
		); err != nil {
			return nil, err
		}

		workouts = append(workouts, &w)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (ws *workoutService) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	qry := `
	UPDATE workouts
	SET activity_type_id = ?,
	duration_minutes = ?,
	distance_miles = ?,
	workout_details = ?,
	workout_date = ?
	WHERE id = ?
	`

	_, err := ws.db.ExecContext(ctx, qry, w.ActivityTypeID, w.DurationMins, w.DistanceMiles, w.WorkoutDetails, w.WorkoutDate, w.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ws *workoutService) DeleteWorkout(ctx context.Context, id int) error {
	qry := `DELETE FROM workouts WHERE id = ?`

	_, err := ws.db.ExecContext(ctx, qry, id)

	if err != nil {
		return err
	}

	return nil
}
