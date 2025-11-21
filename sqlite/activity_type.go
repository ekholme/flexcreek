package sqlite

import (
	"context"
	"database/sql"

	"github.com/ekholme/flexcreek"
)

type activityTypeService struct {
	db *sql.DB
}

func NewActivityTypeService(db *sql.DB) flexcreek.ActivityTypeService {
	return &activityTypeService{db: db}
}

func (ats *activityTypeService) CreateActivityType(ctx context.Context, at *flexcreek.ActivityType) (int, error) {
	qry := `INSERT INTO activity_types (name) VALUES (?)`
	res, err := ats.db.ExecContext(ctx, qry, at.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (ats *activityTypeService) GetAllActivityTypes(ctx context.Context) ([]*flexcreek.ActivityType, error) {
	qry := `SELECT id, name FROM activity_types`
	rows, err := ats.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*flexcreek.ActivityType
	for rows.Next() {
		var at flexcreek.ActivityType
		if err := rows.Scan(&at.ID, &at.Name); err != nil {
			return nil, err
		}
		types = append(types, &at)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}

func (ats *activityTypeService) GetActivityTypeByID(ctx context.Context, id int) (*flexcreek.ActivityType, error) {
	qry := `SELECT id, name FROM activity_types WHERE id = ?`
	var at flexcreek.ActivityType
	err := ats.db.QueryRowContext(ctx, qry, id).Scan(&at.ID, &at.Name)
	if err != nil {
		return nil, err
	}
	return &at, nil
}
