package sqlite

import (
	"database/sql"

	"github.com/ekholme/flexcreek"
)

type movementService struct {
	db  *sql.DB
	mus flexcreek.MuscleService
}

func NewMovementService(db *sql.DB, mus flexcreek.MuscleService) flexcreek.MovementService {

	return &movementService{
		db:  db,
		mus: mus,
	}
}

// methods ---------
func (ms movementService) CreateMovement(m *flexcreek.Movement) (int, error) {
	//insert movement
	mvStmt, err := ms.db.Prepare("INSERT INTO movements (name) VALUES (?)")

	if err != nil {
		return 0, err
	}

	mvRes, err := mvStmt.Exec(m.Name)

	if err != nil {
		return 0, err
	}

	//retrieve the auto-incremented id from the last write
	mvID, err := mvRes.LastInsertId()

	if err != nil {
		return 0, err
	}

	//insert muscles
	for _, m := range m.Muscles {

		msID, err := ms.mus.CreateMuscle(m)

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
	movement := &flexcreek.Movement{}

	err := ms.db.QueryRow("SELECT * FROM movements WHERE movement_id = ?", id).Scan(&movement.ID, &movement.Name, &movement.CreatedAt, &movement.UpdatedAt)

	if err != nil {
		return nil, err
	}

	muscles, err := ms.GetMovementMuscles(movement)

	if err != nil {
		return nil, err
	}

	movement.Muscles = muscles

	return movement, nil
}

func (ms movementService) GetMovementByName(name string) (*flexcreek.Movement, error) {
	movement := &flexcreek.Movement{}

	err := ms.db.QueryRow("SELECT * FROM movements WHERE name = ?", name).Scan(&movement.ID, &movement.Name, &movement.CreatedAt, &movement.UpdatedAt)

	if err != nil {
		return nil, err
	}

	muscles, err := ms.GetMovementMuscles(movement)

	if err != nil {
		return nil, err
	}

	movement.Muscles = muscles

	return movement, nil
}

func (ms movementService) GetAllMovements() ([]*flexcreek.Movement, error) {
	var movements []*flexcreek.Movement

	qry := "SELECT * FROM movements"

	rows, err := ms.db.Query(qry)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		movement := &flexcreek.Movement{}

		err = rows.Scan(&movement.ID, &movement.Name, &movement.CreatedAt, &movement.UpdatedAt)

		if err != nil {
			return nil, err
		}

		muscles, err := ms.GetMovementMuscles(movement)

		if err != nil {
			return nil, err
		}

		movement.Muscles = muscles

		movements = append(movements, movement)
	}

	return movements, nil
}

func (ms movementService) DeleteMovement(id int) (int, error) {
	//TODO
	return 0, nil
}

// utility to get all muscles associated with a movement
// this isn't currently in the MovementService interface, but I think that's ok for now?
func (ms movementService) GetMovementMuscles(m *flexcreek.Movement) ([]*flexcreek.Muscle, error) {
	qry := `
		SELECT mu.muscle_id,
		mu.name,
		mu.created_at,
		mu.updated_at
		FROM muscles mu
		INNER JOIN (
			SELECT muscle_id
			FROM movement_muscles
			WHERE movement_id = ?
		) mm
		ON mu.muscle_id = mm.muscle_id
	`
	rows, err := ms.db.Query(qry, m.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var muscles []*flexcreek.Muscle

	for rows.Next() {
		muscle := &flexcreek.Muscle{}

		err = rows.Scan(&muscle.ID, &muscle.Name, &muscle.CreatedAt, &muscle.UpdatedAt)

		if err != nil {
			return nil, err
		}

		muscles = append(muscles, muscle)
	}

	return muscles, nil

}
