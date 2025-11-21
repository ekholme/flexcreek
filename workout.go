package flexcreek

import (
	"context"
	"time"
)

type Workout struct {
	ID             int          `db:"id"`
	UserID         int          `db:"user_id"`
	ActivityType   ActivityType `db:"activity_type"` // This will be populated by a JOIN
	DurationMins   float64      `db:"duration_minutes"`
	DistanceMiles  float64      `db:"distance_miles"`
	WorkoutDetails string       `db:"workout_details"`
	WorkoutDate    time.Time    `db:"workout_date"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at"`
	// ActivityTypeID is useful for create/update operations without fetching the whole ActivityType object.
	ActivityTypeID int `db:"activity_type_id"`
}

type ActivityType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, w *Workout) (int, error)
	GetWorkoutByID(ctx context.Context, id int) (*Workout, error)
	// Get the n latest workouts for a specific user
	LatestWorkouts(ctx context.Context, userID, n int) ([]*Workout, error)
	// Get all workouts for a specific user
	GetWorkoutsByUser(ctx context.Context, userID int) ([]*Workout, error)
	GetWorkoutsByActivityType(ctx context.Context, userID int, activityTypeID int) ([]*Workout, error)
	UpdateWorkout(ctx context.Context, w *Workout) error
	DeleteWorkout(ctx context.Context, id int) error
}

type ActivityTypeService interface {
	CreateActivityType(ctx context.Context, at *ActivityType) (int, error)
	GetAllActivityTypes(ctx context.Context) ([]*ActivityType, error)
	GetActivityTypeByID(ctx context.Context, id int) (*ActivityType, error)
}
