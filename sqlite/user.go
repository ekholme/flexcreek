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
	return userService{
		db: db,
	}
}

// Create a new user in the users table of the database
func (us userService) CreateUser(ctx context.Context, username string) (int, error) {
	qry := `
		INSERT INTO users (username)
		VALUES (?)	
	`

	res, err := us.db.ExecContext(ctx, qry, username)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (us userService) GetUserByUsername(ctx context.Context, username string) (*flexcreek.User, error) {
	return nil, nil
}

func (us userService) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {
	return nil, nil
}

func (us userService) GetAllUsernames(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (us userService) DeleteUser(ctx context.Context, id int) error {
	return nil
}
