package sqlite

import (
	"context"
	"database/sql"

	"github.com/ekholme/flexcreek"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) flexcreek.UserService {
	return &userService{
		db: db,
	}
}

// defining methods
func (us *userService) CreateUser(ctx context.Context, user *flexcreek.User) (int, error) {
	qry := `
		INSERT INTO users (first_name, last_name, email, hashed_password)
		VALUES (?, ?, ?, ?)	
	`

	stmt, err := us.db.PrepareContext(ctx, qry)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.HashedPw)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func (us *userService) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {

	user := &flexcreek.User{}

	qry := `
		SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		created_at,
		updated_at	
		FROM users
		WHERE id = ?
	`

	err := us.db.QueryRowContext(ctx, qry, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPw, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*flexcreek.User, error) {
	user := &flexcreek.User{}

	qry := `
		SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		created_at,
		updated_at	
		FROM users
		WHERE email = ?
	`

	err := us.db.QueryRowContext(ctx, qry, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPw, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (us *userService) GetAllUsers(ctx context.Context) ([]*flexcreek.User, error) {
	users := []*flexcreek.User{}

	qry := `
		SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		created_at,
		updated_at	
		FROM users
	`

	rows, err := us.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &flexcreek.User{}

		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPw, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) UpdateUser(ctx context.Context, user *flexcreek.User) error {
	qry := `
		UPDATE users
		SET
			first_name = ?,
			last_name = ?,
			email = ?,
			hashed_password = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := us.db.ExecContext(ctx, qry, user.FirstName, user.LastName, user.Email, user.HashedPw, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, id int) error {
	qry := `DELETE FROM users WHERE id = ?`
	res, err := us.db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (us *userService) Login(ctx context.Context, email, password string) (*flexcreek.User, error) {
	// Get user by email. The GetUserByEmail method will return sql.ErrNoRows if not found.
	user, err := us.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			// To prevent user enumeration attacks, we return a generic error
			// for both "user not found" and "wrong password".
			return nil, flexcreek.ErrInvalidCredentials
		}
		// A different database error occurred.
		return nil, err
	}

	// Compare the provided password with the stored hash.
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPw), []byte(password))
	if err != nil {
		// This error is likely bcrypt.ErrMismatchedHashAndPassword, but could be others.
		// We return the same generic error for security.
		return nil, flexcreek.ErrInvalidCredentials
	}

	// Passwords match, login is successful.
	return user, nil
}
