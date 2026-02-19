package flexcreek

import (
	"context"
	"time"
)

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type UserService interface {
	CreateUser(ctx context.Context, username string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetAllUsernames(ctx context.Context) ([]string, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}
