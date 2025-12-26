package flexcreek

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type UserService interface {
	CreateUser(ctx context.Context, username string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}
