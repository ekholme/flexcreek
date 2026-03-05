package flexcreek

import (
	"context"
	"time"
)

type Workout struct {
	ID               int       `db:"id"`
	UserID           int       `db:"user_id"`
	ShortDescription string    `db:"short_description"`
	LongDescription  string    `db:"long_description"`
	WorkoutDate      time.Time `db:"workout_date"`
	CreatedAt        time.Time `db:"created_at"`
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, w *Workout) (int, error)
	GetWorkoutByID(ctx context.Context, id int) (*Workout, error)
	GetWorkoutByDate(ctx context.Context, date time.Time) (*Workout, error)
	GetLatestWorkouts(ctx context.Context, n int) ([]*Workout, error)
	UpdateWorkout(ctx context.Context, w *Workout) error
	DeleteWorkout(ctx context.Context, id int) error
}
