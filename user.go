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
	Weight    float64   `db:"weight"`
	CreatedAt time.Time `db:"created_at"`
}

type UserService interface {
	CreateUser(ctx context.Context, u *User, password string) (int, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int) error
	Login(ctx context.Context, email, password string) (*User, error)
}
