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

//methods ----------------

func (us userService) CreateUser(ctx context.Context, u *flexcreek.User) (int, error) {
	qry := `
	INSERT INTO users (username, email, hashed_password)
	VALUES (?, ?, ?)
	`

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

func (us userService) GetUserByID(ctx context.Context, id int) (*flexcreek.User, error) {
	return nil, nil
}

func (us userService) GetUserByEmail(ctx context.Context, email string) (*flexcreek.User, error) {
	return nil, nil
}

func (us userService) GetAllUsers(ctx context.Context) ([]*flexcreek.User, error) {
	return nil, nil
}

func (us userService) UpdateUser(ctx context.Context, u *flexcreek.User) error {
	return nil
}

func (us userService) DeleteUser(ctx context.Context, id int) error {
	return nil
}

func (us userService) Login(ctx context.Context, email string, password string) (*flexcreek.User, error) {
	return nil, nil
}

// utilities --------------

func hashPW(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil

}

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
