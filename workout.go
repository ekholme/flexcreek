package flexcreek

import (
	"context"
	"time"
)

type Workout struct {
	ID             int          `db:"id"`
	UserID         int          `db:"user_id"`
	ActivityType   ActivityType `db:"activity_type"` //note that this isn't the actual name of this col in the db, but it should be given this alias during a query
	DurationMins   float64      `db:"duration_minutes"`
	DistanceMiles  float64      `db:"distance_miles"`
	WorkoutDetails string       `db:"workout_details"`
	WorkoutDate    time.Time    `db:"workout_date"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at"`
}

type ActivityType string

type WorkoutService interface {
	CreateWorkout(ctx context.Context, w *Workout)
	Latest() ([]*Workout, error) //get the n latest workouts
}
