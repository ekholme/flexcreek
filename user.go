package flexcreek

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	HashedPw  string    `json:"hashedPw"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserService interface {
	CreateUser(ctx context.Context, u *User) (int, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int) error
	Login(ctx context.Context, email, password string) (*User, error)
}
