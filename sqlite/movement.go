package sqlite

import (
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

// methods ---------
func (ms movementService) CreateMovement(name string, muscles []string) (int, error) {
	//insert movement
	mvStmt, err := ms.db.Prepare("INSERT INTO movements (name) VALUES (?)")

	if err != nil {
		return 0, err
	}

	mvRes, err := mvStmt.Exec(name)

	if err != nil {
		return 0, err
	}

	//retrieve the auto-incremented id from the last write
	mvID, err := mvRes.LastInsertId()

	if err != nil {
		return 0, err
	}

	//insert muscles
	for _, m := range muscles {
		msStmt, err := ms.db.Prepare("INSERT INTO muscles (name) VALUES (?)")

		if err != nil {
			return 0, err
		}

		msRes, err := msStmt.Exec(m)

		if err != nil {
			return 0, err
		}

		msID, err := msRes.LastInsertId()

		if err != nil {
			return 0, err
		}

		stmt, err := ms.db.Prepare("INSERT INTO movement_muscles (movement_id, muscle_id) VALUES (?, ?)")

		if err != nil {
			return 0, err
		}

		_, err = stmt.Exec(mvID, msID)

		if err != nil {
			return 0, err
		}

	}

	return int(mvID), nil
}

func (ms movementService) GetMovementByID(id int) (*flexcreek.Movement, error) {
	//TODO
	return nil, nil
}

func (ms movementService) GetMovementByName(name string) (*flexcreek.Movement, error) {
	//TODO
	return nil, nil
}

func (ms movementService) GetAllMovements() ([]*flexcreek.Movement, error) {
	//TODO
	return nil, nil
}

func (ms movementService) DeleteMovement(id int) (int, error) {
	//TODO
	return 0, nil
}
