package sqlite

import (
	"context"
	"database/sql"
	"fmt"

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

// methods ----------------
func (us *userService) CreateUser(ctx context.Context, u *flexcreek.User, password string) (int, error) {
	qry := `
	INSERT INTO users (username, email, hashed_password)
	VALUES (?, ?, ?)
	`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("could not hash password: %w", err)
	}

	u.PasswordHash = hashedPassword

	res, err := us.db.ExecContext(ctx, qry, u.UserName, u.Email, u.PasswordHash)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (us *userService) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {
	qry := `
	SELECT id,
	username,
	email,
	hashed_password,
	created_at,
	updated_at
	FROM users
	WHERE id = ?
	`

	var u = &flexcreek.User{}

	err := us.db.QueryRowContext(ctx, qry, id).Scan(&u.ID, &u.UserName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*flexcreek.User, error) {
	qry := `
	SELECT id,
	username,
	email,
	hashed_password,
	created_at,
	updated_at
	FROM users
	WHERE email = ?
	`

	var u = &flexcreek.User{}

	err := us.db.QueryRowContext(ctx, qry, email).Scan(&u.ID, &u.UserName, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil

}

func (us *userService) UpdateUser(ctx context.Context, u *flexcreek.User) error {
	qry := `
	UPDATE users
	SET username = ?,
	email = ?,
	hashed_password = ?
	WHERE id = ?
	`

	_, err := us.db.ExecContext(ctx, qry, u.UserName, u.Email, u.PasswordHash, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, id int) error {
	qry := `DELETE FROM users WHERE id = ?`
	_, err := us.db.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Login(ctx context.Context, email string, password string) (*flexcreek.User, error) {
	user, err := us.GetUserByEmail(ctx, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, flexcreek.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: %w", err)
	}

	ok, err := checkPasswordHash(password, string(user.PasswordHash))
	if err != nil {
		return nil, fmt.Errorf("login: could not check password hash: %w", err)
	}

	if !ok {
		return nil, flexcreek.ErrInvalidCredentials
	}

	return user, nil
}

// utility ---------
func checkPasswordHash(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err == nil {
		return true, nil
	}

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, fmt.Errorf("unexpected error during password comparision: %w", err)
}
