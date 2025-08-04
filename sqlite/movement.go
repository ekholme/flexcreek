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
	//todo
	return nil, nil
}

func (ms *movementService) GetMovementByName(ctx context.Context, name string) (*flexcreek.Movement, error) {
	//todo
	return nil, nil
}

func (ms *movementService) GetAllMovements(ctx context.Context) ([]*flexcreek.Movement, error) {
	//todo
	return nil, nil
}

func (ms *movementService) UpdateMovement(ctx context.Context, m *flexcreek.Movement) error {
	//todo
	return nil
}

func (ms *movementService) DeleteMovement(ctx context.Context, id int) error {
	//todo
	return nil
}
