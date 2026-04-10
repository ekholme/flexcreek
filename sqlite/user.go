package sqlite

import (
	"context"
	"database/sql"

	"github.com/ekholme/flexcreek"
)

// Create a new user in the users table of the database
func (s *Storage) CreateUser(ctx context.Context, username string) (int, error) {
	qry := `
		INSERT INTO users (username)
		VALUES (?)	
	`

	res, err := s.db.ExecContext(ctx, qry, username)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*flexcreek.User, error) {
	qry := `
		SELECT id,
		username,
		created_at
		FROM users
		WHERE username = ?	
	`

	var u flexcreek.User

	res := s.db.QueryRowContext(ctx, qry, username)

	//todo later -- define a custom error type to handle cases where no rows are returned & return this instead
	if err := res.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {
	qry := `
		SELECT id,
		username,
		created_at
		FROM users
		WHERE id = ?	
	`

	var u flexcreek.User

	res := s.db.QueryRowContext(ctx, qry, id)

	//todo later -- define a custom error type to handle cases where no rows are returned & return this instead
	if err := res.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil

}

func (s *Storage) GetAllUsers(ctx context.Context) ([]*flexcreek.User, error) {
	qry := `
		SELECT id,
		username,
		created_at
		FROM users	
	`

	rows, err := s.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*flexcreek.User
	for rows.Next() {
		var user flexcreek.User
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	//check for any errors that occur during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	qry := `
		DELETE FROM users
		WHERE id = ?	
	`

	res, err := s.db.ExecContext(ctx, qry, id)
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
