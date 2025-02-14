package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ekholme/flexcreek"
)

type muscleService struct {
	db *sql.DB
}

func NewMuscleService(db *sql.DB) flexcreek.MuscleService {
	return &muscleService{
		db: db,
	}
}

//methods ------------

func (mus muscleService) CreateMuscle(muscle *flexcreek.Muscle) (int, error) {
	musStmt, err := mus.db.Prepare("INSERT INTO muscles (name) VALUES (?)")

	if err != nil {
		return 0, err
	}

	musRes, err := musStmt.Exec(muscle.Name)

	if err != nil {
		return 0, err
	}

	musID, err := musRes.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(musID), nil
}

func (mus muscleService) GetMuscleByID(id int) (*flexcreek.Muscle, error) {
	muscle := &flexcreek.Muscle{}

	err := mus.db.QueryRow("SELECT * FROM muscles WHERE muscle_id = ?", id).Scan(&muscle.ID, &muscle.Name, &muscle.CreatedAt, &muscle.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return muscle, nil
}

func (mus muscleService) GetMuscleByName(name string) (*flexcreek.Muscle, error) {

	muscle := &flexcreek.Muscle{}

	err := mus.db.QueryRow("SELECT * FROM muscles WHERE name = ?", name).Scan(&muscle.ID, &muscle.Name, &muscle.CreatedAt, &muscle.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return muscle, nil
}

func (mus muscleService) GetAllMuscles() ([]*flexcreek.Muscle, error) {
	rows, err := mus.db.Query("SELECT * FROM muscles")

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

func (mus muscleService) DeleteMuscle(id int) (int, error) {
	musStmt, err := mus.db.Prepare("DELETE FROM muscles WHERE id =  ?")

	if err != nil {
		return 0, err
	}

	musRes, err := musStmt.Exec(id)

	if err != nil {
		return 0, err
	}

	//return the number of rows affected
	rowsAffected, err := musRes.RowsAffected()

	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("muscle with ID %d not found", id)
	}

	return int(rowsAffected), nil
}
