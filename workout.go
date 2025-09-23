package flexcreek

import (
	"context"
	"time"
)

type Workout struct {
	ID                int                 `json:"id"`
	UserID            int                 `json:"userId"`
	Date              time.Time           `json:"date"`
	MovementInstances []*MovementInstance `json:"movementInstances"`
	Notes             string              `json:"notes"`
	Duration          time.Duration       `json:"duration"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, w *Workout) (int, error)
	GetWorkoutByID(ctx context.Context, id int) (*Workout, error)
	GetAllWorkoutsByUser(ctx context.Context, userID int) ([]*Workout, error)
	GetWorkoutsByDate(ctx context.Context, userID int, d time.Time) ([]*Workout, error)
	GetWorkoutsByDateRange(ctx context.Context, userID int, start time.Time, end time.Time) ([]*Workout, error)
	UpdateWorkout(ctx context.Context, w *Workout) error
	DeleteWorkout(ctx context.Context, id int) error
}
