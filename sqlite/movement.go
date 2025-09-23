package sqlite

import (
	"context"
	"database/sql"

	"github.com/ekholme/flexcreek"
)

type movementService struct {
	db *sql.DB
}

func NewMovementService(db *sql.DB) flexcreek.MovementService {
	return &movementService{
		db: db,
	}
}

// Methods ---------

func (ms *movementService) CreateMovement(ctx context.Context, m *flexcreek.Movement) (int, error) {
	qry := `
		INSERT INTO movements (name, movement_type, movement_description) VALUES (?, ?, ?)	
	`

	stmt, err := ms.db.PrepareContext(ctx, qry)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, m.Name, m.MovementType, m.Description)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func (ms *movementService) GetMovementByID(ctx context.Context, id int) (*flexcreek.Movement, error) {
	mv := &flexcreek.Movement{}

	qry := `
		SELECT id,
		name,
		movement_type,
		movement_description,
		created_at,
		updated_at
		FROM movements
		WHERE id = ?
	`

	err := ms.db.QueryRowContext(ctx, qry, id).Scan(&mv.ID, &mv.Name, &mv.MovementType, &mv.Description, &mv.CreatedAt, &mv.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return mv, nil
}

func (ms *movementService) GetMovementByName(ctx context.Context, name string) (*flexcreek.Movement, error) {
	mv := &flexcreek.Movement{}

	qry := `
		SELECT id,
		name,
		movement_type,
		movement_description,
		created_at,
		updated_at
		FROM movements
		WHERE name = ?
	`

	err := ms.db.QueryRowContext(ctx, qry, name).Scan(&mv.ID, &mv.Name, &mv.MovementType, &mv.Description, &mv.CreatedAt, &mv.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return mv, nil
}

func (ms *movementService) GetAllMovements(ctx context.Context) ([]*flexcreek.Movement, error) {
	qry := `
		SELECT id,
		name,
		movement_type,
		movement_description,
		created_at,
		updated_at
		FROM movements
	`

	rows, err := ms.db.QueryContext(ctx, qry)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mvs := []*flexcreek.Movement{}

	for rows.Next() {
		mv := &flexcreek.Movement{}

		err = rows.Scan(&mv.ID, &mv.Name, &mv.MovementType, &mv.Description, &mv.CreatedAt, &mv.UpdatedAt)

		if err != nil {
			return nil, err
		}

		mvs = append(mvs, mv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mvs, nil
}

func (ms *movementService) GetAllMovementsByType(ctx context.Context, movementType flexcreek.MovementType) ([]*flexcreek.Movement, error) {
	qry := `
		SELECT id,
		name,
		movement_type,
		movement_description,
		created_at,
		updated_at
		FROM movements
		WHERE movement_type = ?
	`

	rows, err := ms.db.QueryContext(ctx, qry, movementType)

	if err != nil {
		return nil, err
	}

	mvs := []*flexcreek.Movement{}

	for rows.Next() {
		mv := &flexcreek.Movement{}

		err = rows.Scan(&mv.ID, &mv.Name, &mv.MovementType, &mv.Description, &mv.CreatedAt, &mv.UpdatedAt)

		if err != nil {
			return nil, err
		}

		mvs = append(mvs, mv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mvs, nil

}

func (ms *movementService) UpdateMovement(ctx context.Context, m *flexcreek.Movement) error {
	qry := `
		UPDATE movements
		SET
			name = ?,
			movement_type = ?,
			movement_description = ?
		WHERE id = ?
	`

	res, err := ms.db.ExecContext(ctx, qry, m.Name, m.MovementType, m.Description, m.ID)
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

func (ms *movementService) DeleteMovement(ctx context.Context, id int) error {
	qry := `DELETE FROM movements WHERE id = ?`

	res, err := ms.db.ExecContext(ctx, qry, id)
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
