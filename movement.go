package flexcreek

import (
	"context"
	"time"
)

// MovementType defines the category of a movement.
// This helps determine what kind of metrics are tracked for it.
type MovementType string

const (
	StrengthMovement MovementType = "strength"
	CardioMovement   MovementType = "cardio"
	AmrapMovement    MovementType = "amrap" // As Many Rounds As Possible
	EmomMovement     MovementType = "emom"  // Every Minute On the Minute
	// Add other types as needed, e.g., "flexibility", "tabata"
)

type Movement struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	MovementType MovementType `json:"movementType"` // e.g., 'strength', 'cardio'
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type MovementService interface {
	CreateMovement(ctx context.Context, m *Movement) (int, error)
	GetMovementByID(ctx context.Context, id int) (*Movement, error)
	GetMovementByName(ctx context.Context, name string) (*Movement, error)
	GetAllMovements(ctx context.Context) ([]*Movement, error)
	GetAllMovementsByType(ctx context.Context, movementType MovementType) ([]*Movement, error)
	UpdateMovement(ctx context.Context, m *Movement) error
	DeleteMovement(ctx context.Context, id int) error
}
