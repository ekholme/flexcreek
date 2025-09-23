package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ekholme/flexcreek"
)

type movementInstanceService struct {
	db querier
}

func NewMovementInstanceService(db *sql.DB) flexcreek.MovementInstanceService {
	return &movementInstanceService{
		db: db, // A *sql.DB satisfies the querier interface
	}
}

// Methods ---------------

func (s *movementInstanceService) CreateMovementInstance(ctx context.Context, mi *flexcreek.MovementInstance) (int, error) {
	if mi.Movement == nil || mi.Movement.ID == 0 {
		return 0, fmt.Errorf("movement with a valid ID is required for a movement instance")
	}

	logData, err := marshalLogData(mi)
	if err != nil {
		return 0, err
	}

	qry := `
		INSERT INTO movement_instances (
			workout_id, movement_id, notes, rpe, log_data
		) VALUES (?, ?, ?, ?, ?)`

	// This method no longer manages its own transaction. It simply executes
	// its statement on the querier it was initialized with.
	res, err := s.db.ExecContext(ctx, qry,
		mi.WorkoutID,
		mi.Movement.ID,
		mi.Notes,
		mi.RPE,
		logData,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *movementInstanceService) GetMovementInstanceByID(ctx context.Context, id int) (*flexcreek.MovementInstance, error) {
	mi := &flexcreek.MovementInstance{
		Movement: &flexcreek.Movement{},
	}
	var logData sql.NullString

	qry := `
		SELECT
			mi.id, mi.workout_id, mi.notes, mi.rpe,
			mi.log_data, mi.created_at, mi.updated_at,
			m.id, m.name, m.movement_type, m.description, m.created_at, m.updated_at
		FROM movement_instances mi
		JOIN movements m ON mi.movement_id = m.id
		WHERE mi.id = ?`

	err := s.db.QueryRowContext(ctx, qry, id).Scan(
		&mi.ID,
		&mi.WorkoutID,
		&mi.Notes,
		&mi.RPE,
		&logData,
		&mi.CreatedAt,
		&mi.UpdatedAt,
		&mi.Movement.ID,
		&mi.Movement.Name,
		&mi.Movement.MovementType,
		&mi.Movement.Description,
		&mi.Movement.CreatedAt,
		&mi.Movement.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := unmarshalLogData(mi, logData, string(mi.Movement.MovementType)); err != nil {
		return nil, err
	}

	return mi, nil
}

func (s *movementInstanceService) GetAllMovementInstancesByWorkoutID(ctx context.Context, workoutID int) ([]*flexcreek.MovementInstance, error) {
	qry := `
		SELECT
		mi.id, mi.workout_id, mi.notes, mi.rpe, mi.log_data, mi.created_at, mi.updated_at, m.id, m.name, m.movement_type, m.description, m.created_at, m.updated_at
		FROM movement_instances mi
		JOIN movements m ON mi.movement_id = m.id
		WHERE mi.workout_id = ?	
	`
	rows, err := s.db.QueryContext(ctx, qry, workoutID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	instances := []*flexcreek.MovementInstance{}

	for rows.Next() {
		mi := &flexcreek.MovementInstance{
			Movement: &flexcreek.Movement{},
		}
		var logData sql.NullString

		if err := rows.Scan(
			&mi.ID, &mi.WorkoutID, &mi.Notes, &mi.RPE,
			&logData, &mi.CreatedAt, &mi.UpdatedAt,
			&mi.Movement.ID,
			&mi.Movement.Name,
			&mi.Movement.MovementType,
			&mi.Movement.Description,
			&mi.Movement.CreatedAt,
			&mi.Movement.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err = unmarshalLogData(mi, logData, string(mi.Movement.MovementType)); err != nil {
			return nil, err
		}

		instances = append(instances, mi)

	}

	return instances, rows.Err()
}

func (s *movementInstanceService) GetAllMovementInstancesForMovement(ctx context.Context, userID int, movementID int) ([]*flexcreek.MovementInstance, error) {
	qry := `
		SELECT
			mi.id, mi.workout_id, mi.notes, mi.rpe, mi.log_data,
			mi.created_at, mi.updated_at, m.id, m.name, m.movement_type, m.description, m.created_at, m.updated_at
		FROM movement_instances mi
		JOIN workouts w ON mi.workout_id = w.id
		JOIN movements m ON mi.movement_id = m.id
		WHERE m.id = ?
		  AND w.user_id = ?
		ORDER BY mi.created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, qry, movementID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instances := []*flexcreek.MovementInstance{}
	for rows.Next() {
		mi := &flexcreek.MovementInstance{
			Movement: &flexcreek.Movement{},
		}
		var logData sql.NullString

		if err := rows.Scan(
			&mi.ID, &mi.WorkoutID, &mi.Notes, &mi.RPE,
			&logData, &mi.CreatedAt, &mi.UpdatedAt,
			&mi.Movement.ID,
			&mi.Movement.Name,
			&mi.Movement.MovementType,
			&mi.Movement.Description,
			&mi.Movement.CreatedAt,
			&mi.Movement.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := unmarshalLogData(mi, logData, string(mi.Movement.MovementType)); err != nil {
			return nil, err
		}

		instances = append(instances, mi)
	}

	return instances, rows.Err()

}

func (s *movementInstanceService) GetAllMovementInstancesByUser(ctx context.Context, userID int) ([]*flexcreek.MovementInstance, error) {
	qry := `
		SELECT
			mi.id, mi.workout_id, mi.notes, mi.rpe, mi.log_data,
			mi.created_at, mi.updated_at, m.id, m.name, m.movement_type, m.description, m.created_at, m.updated_at
		FROM movement_instances mi
		JOIN workouts w ON mi.workout_id = w.id
		JOIN movements m ON mi.movement_id = m.id
		WHERE w.user_id = ?`

	rows, err := s.db.QueryContext(ctx, qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instances := []*flexcreek.MovementInstance{}
	for rows.Next() {
		mi := &flexcreek.MovementInstance{
			Movement: &flexcreek.Movement{},
		}
		var logData sql.NullString

		if err := rows.Scan(
			&mi.ID, &mi.WorkoutID, &mi.Notes, &mi.RPE,
			&logData, &mi.CreatedAt, &mi.UpdatedAt,
			&mi.Movement.ID,
			&mi.Movement.Name,
			&mi.Movement.MovementType,
			&mi.Movement.Description,
			&mi.Movement.CreatedAt,
			&mi.Movement.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := unmarshalLogData(mi, logData, string(mi.Movement.MovementType)); err != nil {
			return nil, err
		}

		instances = append(instances, mi)
	}

	return instances, rows.Err()
}

func (s *movementInstanceService) UpdateMovementInstance(ctx context.Context, mi *flexcreek.MovementInstance) error {
	logData, err := marshalLogData(mi)
	if err != nil {
		return err
	}

	if mi.Movement == nil || mi.Movement.ID == 0 {
		return fmt.Errorf("movement with a valid ID is required to update a movement instance")
	}

	qry := `
		UPDATE movement_instances
		SET
			workout_id = ?, movement_id = ?, notes = ?, rpe = ?, log_data = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	// Note: This method now relies on the caller to wrap it in a transaction
	// if atomic operations are needed.
	res, err := s.db.ExecContext(ctx, qry,
		mi.WorkoutID, mi.Movement.ID, mi.Notes, mi.RPE,
		logData, mi.ID,
	)
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

func (s *movementInstanceService) DeleteMovementInstance(ctx context.Context, id int) error {
	qry := `DELETE FROM movement_instances WHERE id = ?`
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

// Helpers ---------------

// marshalLogData serializes the log data from a MovementInstance into a single sql.NullString.
func marshalLogData(mi *flexcreek.MovementInstance) (sql.NullString, error) {
	var data any
	if mi.Strength != nil {
		data = mi.Strength
	} else if mi.Cardio != nil {
		data = mi.Cardio
	} else if mi.Amrap != nil {
		data = mi.Amrap
	} else if mi.Emom != nil {
		data = mi.Emom
	}

	if data == nil {
		return sql.NullString{}, nil // No log data to marshal
	}

	b, err := json.Marshal(data)
	if err != nil {
		return sql.NullString{}, fmt.Errorf("failed to marshal movement instance log: %w", err)
	}

	return sql.NullString{String: string(b), Valid: true}, nil
}

// unmarshalLogData deserializes JSON from log_data into the correct log struct in a MovementInstance,
// based on the movement's type. This assumes the movement type strings in the database
// are "strength", "cardio", "amrap", and "emom".
func unmarshalLogData(mi *flexcreek.MovementInstance, logData sql.NullString, movementType string) error {
	if !logData.Valid || logData.String == "" {
		return nil // No data to unmarshal
	}

	dataBytes := []byte(logData.String)

	switch movementType {
	case "strength":
		mi.Strength = &flexcreek.StrengthLog{}
		if err := json.Unmarshal(dataBytes, mi.Strength); err != nil {
			return fmt.Errorf("failed to unmarshal strength log: %w", err)
		}
	case "cardio":
		mi.Cardio = &flexcreek.CardioLog{}
		if err := json.Unmarshal(dataBytes, mi.Cardio); err != nil {
			return fmt.Errorf("failed to unmarshal cardio log: %w", err)
		}
	case "amrap":
		mi.Amrap = &flexcreek.AmrapLog{}
		if err := json.Unmarshal(dataBytes, mi.Amrap); err != nil {
			return fmt.Errorf("failed to unmarshal amrap log: %w", err)
		}
	case "emom":
		mi.Emom = &flexcreek.EmomLog{}
		if err := json.Unmarshal(dataBytes, mi.Emom); err != nil {
			return fmt.Errorf("failed to unmarshal emom log: %w", err)
		}
	}
	return nil
}
