package sqlite

import (
	"context"
	"database/sql"

	"github.com/ekholme/flexcreek"
)

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) flexcreek.UserService {
	return &userService{
		db: db,
	}
}

// methods ----------------
func (us *userService) CreateUser(ctx context.Context, username string) (int, error) {
	qry := `
	INSERT INTO users (username)
	VALUES (?)
	`

	res, err := us.db.ExecContext(ctx, username)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (us *userService) GetUserByUsername(ctx context.Context, username string) (*flexcreek.User, error) {
	qry := `
	SELECT id,
	username,
	created_at
	FROM users
	WHERE username = ?
	`

	res, err := us.db.QueryRowContext(ctx, qry, username)

	// RESUME HERE
}

func (us *userService) DeleteUser(ctx context.Context, id int) error {
	qry := `DELETE FROM users WHERE id = ?`
	_, err := us.db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}
	return nil
}
