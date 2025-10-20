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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case of error. The rollback will be a no-op if the
	// transaction is successfully committed.
	defer tx.Rollback()

	// Truncate the time part to ensure we're only storing the date.
	dateOnly := w.Date.Truncate(24 * time.Hour)

	// 1. Insert the parent Workout record within the transaction.
	qry := `INSERT INTO workouts (user_id, workout_date, notes, duration_seconds) VALUES (?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, qry, w.UserID, dateOnly, w.Notes, int64(w.Duration.Seconds()))
	if err != nil {
		return 0, fmt.Errorf("failed to insert workout: %w", err)
	}
	workoutID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id for workout: %w", err)
	}

	// 2. Create a movement instance service that operates on our transaction.
	txMovementInstanceSvc := &movementInstanceService{db: tx}

	// 3. Loop through and create each MovementInstance using the transactional service.
	for _, mi := range w.MovementInstances {
		mi.WorkoutID = int(workoutID)
		if _, err := txMovementInstanceSvc.CreateMovementInstance(ctx, mi); err != nil {
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
	// This query fetches the workout and all its movement instances in one go.
	qry := `
		SELECT
			w.id, w.user_id, w.workout_date, w.notes, w.duration_seconds, w.created_at, w.updated_at,
			mi.id, mi.notes, mi.rpe, mi.log_data, mi.created_at, mi.updated_at,
			m.id, m.name, m.movement_type, m.created_at, m.updated_at
		FROM workouts w
		LEFT JOIN movement_instances mi ON w.id = mi.workout_id
		LEFT JOIN movements m ON mi.movement_id = m.id
		WHERE w.id = ?
		ORDER BY mi.id
	`
	rows, err := s.db.QueryContext(ctx, qry, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query workout by id: %w", err)
	}
	defer rows.Close()

	var workout *flexcreek.Workout
	for rows.Next() {
		var mi flexcreek.MovementInstance
		var m flexcreek.Movement
		var miID, miRPE sql.NullInt64
		var miNotes, miLogData sql.NullString
		var miCreatedAt, miUpdatedAt sql.NullTime
		var movementID sql.NullInt64
		var movementName, movementType sql.NullString
		var movementCreatedAt, movementUpdatedAt sql.NullTime

		// Initialize workout only on the first row
		if workout == nil {
			workout = &flexcreek.Workout{}
		}

		err := rows.Scan(
			&workout.ID, &workout.UserID, &workout.Date, &workout.Notes, &workout.Duration, &workout.CreatedAt, &workout.UpdatedAt,
			&miID, &miNotes, &miRPE, &miLogData, &miCreatedAt, &miUpdatedAt,
			&movementID, &movementName, &movementType, &movementCreatedAt, &movementUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout data: %w", err)
		}

		// If miID is valid, it means we have a movement instance from the LEFT JOIN
		if miID.Valid {
			mi.ID = int(miID.Int64)
			mi.WorkoutID = workout.ID
			mi.Notes = miNotes.String
			if miRPE.Valid {
				rpe := int(miRPE.Int64)
				mi.RPE = &rpe
			}
			mi.CreatedAt = miCreatedAt.Time
			mi.UpdatedAt = miUpdatedAt.Time

			m.ID = int(movementID.Int64)
			m.Name = movementName.String
			m.MovementType = flexcreek.MovementType(movementType.String)
			m.CreatedAt = movementCreatedAt.Time
			m.UpdatedAt = movementUpdatedAt.Time
			mi.Movement = &m

			if err := unmarshalLogData(&mi, miLogData, string(m.MovementType)); err != nil {
				return nil, fmt.Errorf("failed to unmarshal log data: %w", err)
			}
			workout.MovementInstances = append(workout.MovementInstances, &mi)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if workout == nil {
		return nil, sql.ErrNoRows
	}

	return workout, nil
}

func (s *workoutService) GetAllWorkoutsByUser(ctx context.Context, userID int) ([]*flexcreek.Workout, error) {
	// This query is complex but efficient. It fetches all workouts and their associated
	// movement instances for a user in a single database call.
	qry := `
		SELECT
			w.id, w.user_id, w.workout_date, w.notes, w.duration_seconds, w.created_at, w.updated_at,
			mi.id, mi.notes, mi.rpe, mi.log_data, mi.created_at, mi.updated_at,
			m.id, m.name, m.movement_type, m.created_at, m.updated_at
		FROM workouts w
		LEFT JOIN movement_instances mi ON w.id = mi.workout_id
		LEFT JOIN movements m ON mi.movement_id = m.id
		WHERE w.user_id = ?
		ORDER BY w.workout_date DESC, w.id, mi.id
	`

	rows, err := s.db.QueryContext(ctx, qry, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workouts by user: %w", err)
	}
	defer rows.Close()

	// Use a map to aggregate movement instances into their parent workouts.
	workoutMap := make(map[int]*flexcreek.Workout)
	var workouts []*flexcreek.Workout // Keep order

	for rows.Next() {
		var w flexcreek.Workout
		var mi flexcreek.MovementInstance
		var m flexcreek.Movement
		var miID, miRPE sql.NullInt64
		var miNotes, miLogData sql.NullString
		var miCreatedAt, miUpdatedAt sql.NullTime
		var movementID sql.NullInt64
		var movementName, movementType sql.NullString
		var movementCreatedAt, movementUpdatedAt sql.NullTime

		err := rows.Scan(
			&w.ID, &w.UserID, &w.Date, &w.Notes, &w.Duration, &w.CreatedAt, &w.UpdatedAt,
			&miID, &miNotes, &miRPE, &miLogData, &miCreatedAt, &miUpdatedAt,
			&movementID, &movementName, &movementType, &movementCreatedAt, &movementUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout data: %w", err)
		}

		// If we haven't seen this workout ID yet, add it to our map and ordered slice.
		if _, ok := workoutMap[w.ID]; !ok {
			workoutMap[w.ID] = &w
			workouts = append(workouts, &w)
		}

		// If miID is valid, we have a movement instance to process.
		if miID.Valid {
			currentWorkout := workoutMap[w.ID]
			mi.ID = int(miID.Int64)
			mi.WorkoutID = w.ID
			mi.Notes = miNotes.String
			if miRPE.Valid {
				rpe := int(miRPE.Int64)
				mi.RPE = &rpe
			}
			mi.CreatedAt = miCreatedAt.Time
			mi.UpdatedAt = miUpdatedAt.Time

			m.ID = int(movementID.Int64)
			m.Name = movementName.String
			m.MovementType = flexcreek.MovementType(movementType.String)
			m.CreatedAt = movementCreatedAt.Time
			m.UpdatedAt = movementUpdatedAt.Time
			mi.Movement = &m

			if err := unmarshalLogData(&mi, miLogData, string(m.MovementType)); err != nil {
				return nil, fmt.Errorf("failed to unmarshal log data: %w", err)
			}
			currentWorkout.MovementInstances = append(currentWorkout.MovementInstances, &mi)
		}
	}

	return workouts, rows.Err()
}

func (s *workoutService) GetWorkoutsByDate(ctx context.Context, userID int, d time.Time) ([]*flexcreek.Workout, error) {
	// This query is optimized to only fetch workouts for the specified user and date.
	qry := `
		SELECT
			w.id, w.user_id, w.workout_date, w.notes, w.duration_seconds, w.created_at, w.updated_at,
			mi.id, mi.notes, mi.rpe, mi.log_data, mi.created_at, mi.updated_at,
			m.id, m.name, m.movement_type, m.created_at, m.updated_at
		FROM workouts w
		LEFT JOIN movement_instances mi ON w.id = mi.workout_id
		LEFT JOIN movements m ON mi.movement_id = m.id
		WHERE w.user_id = ? AND date(w.workout_date) = date(?)
		ORDER BY w.workout_date DESC, w.id, mi.id
	`

	rows, err := s.db.QueryContext(ctx, qry, userID, d.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query workouts by date: %w", err)
	}
	defer rows.Close()

	return scanWorkouts(rows)
}

func (s *workoutService) GetWorkoutsByDateRange(ctx context.Context, userID int, start time.Time, end time.Time) ([]*flexcreek.Workout, error) {
	// This query is optimized to only fetch workouts within the specified date range.
	qry := `
		SELECT
			w.id, w.user_id, w.workout_date, w.notes, w.duration_seconds, w.created_at, w.updated_at,
			mi.id, mi.notes, mi.rpe, mi.log_data, mi.created_at, mi.updated_at,
			m.id, m.name, m.movement_type, m.created_at, m.updated_at
		FROM workouts w
		LEFT JOIN movement_instances mi ON w.id = mi.workout_id
		LEFT JOIN movements m ON mi.movement_id = m.id
		WHERE w.user_id = ? AND w.workout_date BETWEEN ? AND ?
		ORDER BY w.workout_date DESC, w.id, mi.id
	`

	rows, err := s.db.QueryContext(ctx, qry, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query workouts by date range: %w", err)
	}
	defer rows.Close()

	return scanWorkouts(rows)
}

func (s *workoutService) UpdateWorkout(ctx context.Context, w *flexcreek.Workout) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Truncate the time part to ensure we're only storing the date.
	dateOnly := w.Date.Truncate(24 * time.Hour)

	// 1. Update the parent workout record.
	qry := `UPDATE workouts SET user_id = ?, workout_date = ?, notes = ?, duration_seconds = ? WHERE id = ?`
	res, err := tx.ExecContext(ctx, qry, w.UserID, dateOnly, w.Notes, int64(w.Duration.Seconds()), w.ID)
	if err != nil {
		return fmt.Errorf("failed to update workout: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// 2. Delete old movement instances for this workout.
	if _, err := tx.ExecContext(ctx, `DELETE FROM movement_instances WHERE workout_id = ?`, w.ID); err != nil {
		return fmt.Errorf("failed to delete old movement instances: %w", err)
	}

	// 3. Insert the new movement instances.
	txMovementInstanceSvc := &movementInstanceService{db: tx}
	for _, mi := range w.MovementInstances {
		mi.WorkoutID = w.ID
		if _, err := txMovementInstanceSvc.CreateMovementInstance(ctx, mi); err != nil {
			return fmt.Errorf("failed to create new movement instance: %w", err)
		}
	}

	return tx.Commit()
}

func (s *workoutService) DeleteWorkout(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// With `ON DELETE CASCADE` in the schema, we no longer need to manually
	// delete the child `movement_instances`. The database will handle it
	// automatically and atomically when the parent `workout` is deleted.
	// If you choose not to update the schema, the original code is correct.

	res, err := tx.ExecContext(ctx, `DELETE FROM workouts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

// scanWorkouts is a helper function to reduce code duplication between GetAllWorkoutsByUser,
// GetWorkoutsByDate, and GetWorkoutsByDateRange. It processes the rows from a complex
// workout query and reconstructs the slice of Workout objects.
func scanWorkouts(rows *sql.Rows) ([]*flexcreek.Workout, error) {
	workoutMap := make(map[int]*flexcreek.Workout)
	var workouts []*flexcreek.Workout // Keep order

	for rows.Next() {
		var w flexcreek.Workout
		var mi flexcreek.MovementInstance
		var m flexcreek.Movement
		var miID, miRPE sql.NullInt64
		var miNotes, miLogData sql.NullString
		var miCreatedAt, miUpdatedAt sql.NullTime
		var movementID sql.NullInt64
		var movementName, movementType sql.NullString
		var movementCreatedAt, movementUpdatedAt sql.NullTime

		err := rows.Scan(
			&w.ID, &w.UserID, &w.Date, &w.Notes, &w.Duration, &w.CreatedAt, &w.UpdatedAt,
			&miID, &miNotes, &miRPE, &miLogData, &miCreatedAt, &miUpdatedAt,
			&movementID, &movementName, &movementType, &movementCreatedAt, &movementUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout data: %w", err)
		}

		if _, ok := workoutMap[w.ID]; !ok {
			workoutMap[w.ID] = &w
			workouts = append(workouts, &w)
		}

		if miID.Valid {
			currentWorkout := workoutMap[w.ID]
			mi.ID = int(miID.Int64)
			mi.WorkoutID = w.ID
			mi.Notes = miNotes.String
			if miRPE.Valid {
				rpe := int(miRPE.Int64)
				mi.RPE = &rpe
			}
			mi.CreatedAt = miCreatedAt.Time
			mi.UpdatedAt = miUpdatedAt.Time

			m.ID = int(movementID.Int64)
			m.Name = movementName.String
			m.MovementType = flexcreek.MovementType(movementType.String)
			m.CreatedAt = movementCreatedAt.Time
			m.UpdatedAt = movementUpdatedAt.Time
			mi.Movement = &m

			if err := unmarshalLogData(&mi, miLogData, string(m.MovementType)); err != nil {
				return nil, fmt.Errorf("failed to unmarshal log data: %w", err)
			}
			currentWorkout.MovementInstances = append(currentWorkout.MovementInstances, &mi)
		}
	}

	return workouts, rows.Err()
}
