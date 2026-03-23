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
	qry := `
		SELECT id,
		username,
		created_at
		FROM users
		WHERE username = ?	
	`

	var u flexcreek.User

	res := us.db.QueryRowContext(ctx, qry, username)

	//todo later -- define a custom error type to handle cases where no rows are returned & return this instead
	if err := res.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us userService) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {
	qry := `
		SELECT id,
		username,
		created_at
		FROM users
		WHERE id = ?	
	`

	var u flexcreek.User

	res := us.db.QueryRowContext(ctx, qry, id)

	//todo later -- define a custom error type to handle cases where no rows are returned & return this instead
	if err := res.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil

}

func (us userService) GetAllUsernames(ctx context.Context) ([]string, error) {
	qry := `
		SELECT username
		FROM users	
	`

	rows, err := us.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	//check for any errors that occur during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}

func (us userService) DeleteUser(ctx context.Context, id int) error {
	qry := `
		DELETE FROM users
		WHERE id = ?	
	`

	res, err := us.db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	// Check if any rows were actually deleted.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, it means the user with that ID was not found.
	if rowsAffected == 0 {
		// In the future, this could return a custom flexcreek.ErrNotFound.
		return sql.ErrNoRows
	}

	return nil
}
