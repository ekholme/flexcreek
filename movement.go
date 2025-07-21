package flexcreek

import (
	"context"
	"time"
)

type Movement struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	MovementType string    `json:"movementType"` //eg 'kettlebell', 'barbell', 'cardio', etc. Not sure if I actually need this
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	//potentially add other fields?
}

type MovementService interface {
	CreateMovement(ctx context.Context, m *Movement) (int, error)
	GetMovementByID(ctx context.Context, id int) (*Movement, error)
	GetMovementByName(ctx context.Context, name string) (*Movement, error)
	GetAllMovements(ctx context.Context) ([]*Movement, error)
	UpdateMovement(ctx context.Context, m *Movement) error
	DeleteMovement(ctx context.Context, id int) error
}
