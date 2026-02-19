package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
// create a new user in the database
func (us *userService) CreateUser(ctx context.Context, username string) (int, error) {
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

// get a user from the database given a username
func (us *userService) GetUserByUsername(ctx context.Context, username string) (*flexcreek.User, error) {
	qry := `
	SELECT id,
	username,
	created_at
	FROM users
	WHERE username = ?
	`

	var u flexcreek.User

	err := us.db.QueryRowContext(ctx, qry, username).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, err
	}

	return &u, nil
}

func (us *userService) GetUserById(ctx context.Context, id int) (*flexcreek.User, error) {
	qry := `
	SELECT id,
	username,
	created_at
	FROM users
	WHERE id = ?
	`

	var u flexcreek.User

	err := us.db.QueryRowContext(ctx, qry, id).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %v not found", id)
		}
		return nil, err
	}

	return &u, nil
}

// get all usernames in the database
func (us *userService) GetAllUsernames(ctx context.Context) ([]string, error) {
	qry := "SELECT username FROM users"

	var ss []string

	rows, err := us.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s string

		if err := rows.Scan(&s); err != nil {
			return nil, err
		}

		ss = append(ss, s)
	}

	// check for errors during iteration
	return ss, rows.Err()
}


// get all users in the database
func (us *userService) GetAllUsers(ctx context.Context) ([]flexcreek.User, error) {
	qry := "SELECT id, username, created_at FROM users"

	var users []flexcreek.User

	rows, err := us.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u flexcreek.User

		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	// check for errors during iteration
	return users, rows.Err()
}

func (us *userService) DeleteUser(ctx context.Context, id int) error {
	qry := `DELETE FROM users WHERE id = ?`
	_, err := us.db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}
	return nil
}
