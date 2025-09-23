package flexcreek

import (
	"context"
	"time"
)

// Set represents a single set of a strength-based movement,
// detailing the number of repetitions and the weight used.
type Set struct {
	Reps   int     `json:"reps"`
	Weight float64 `json:"weight"` // Using float64 for weight to allow for fractional values.
}

// StrengthLog holds the metrics for a strength-based movement instance,
// such as squats or bench press. It consists of one or more sets.
type StrengthLog struct {
	Sets []Set `json:"sets"`
}

// CardioLog holds the metrics for a cardio-based movement instance,
// such as running or cycling.
type CardioLog struct {
	Distance float64       `json:"distance,omitempty"` // Distance in a standard unit (e.g., meters or miles).
	Duration time.Duration `json:"duration,omitempty"` // The total time of the activity.
}

// AmrapLog holds the metrics for an AMRAP (As Many Rounds/Reps As Possible) workout.
type AmrapLog struct {
	Duration       time.Duration `json:"duration"`       // How long the AMRAP lasted.
	Rounds         int           `json:"rounds"`         // Total rounds completed.
	ExtraReps      int           `json:"extraReps"`      // Any additional reps completed in the final partial round.
	PrescribedWork string        `json:"prescribedWork"` // Description of the work per round.
}

// EmomLog holds the metrics for an EMOM (Every Minute On the Minute) workout.
type EmomLog struct {
	Duration      time.Duration `json:"duration"`      // Total duration of the EMOM.
	WorkPerMinute string        `json:"workPerMinute"` // Description of the work to be done each minute.
}

// MovementInstance represents a specific performance of a Movement within a Workout.
// For example, 3 sets of 5 reps of Squats, or a 30-minute run.
// It uses composition to hold specific metrics for different types of activities.
type MovementInstance struct {
	ID        int       `json:"id"`
	WorkoutID int       `json:"workoutId"`
	Movement  *Movement `json:"movement"`
	Notes     string    `json:"notes,omitempty"`
	RPE       *int      `json:"rpe,omitempty"` // Pointer to allow for null value (optional field). 1-10 scale.

	// Depending on the type of movement, one of the following will be populated.
	Strength *StrengthLog `json:"strength,omitempty"`
	Cardio   *CardioLog   `json:"cardio,omitempty"`
	Amrap    *AmrapLog    `json:"amrap,omitempty"`
	Emom     *EmomLog     `json:"emom,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MovementInstanceService interface {
	CreateMovementInstance(ctx context.Context, mi *MovementInstance) (int, error)
	GetMovementInstanceByID(ctx context.Context, id int) (*MovementInstance, error)
	GetAllMovementInstancesByWorkoutID(ctx context.Context, workoutID int) ([]*MovementInstance, error)
	GetAllMovementInstancesForMovement(ctx context.Context, userID int, movementID int) ([]*MovementInstance, error)
	GetAllMovementInstancesByUser(ctx context.Context, userID int) ([]*MovementInstance, error)
	UpdateMovementInstance(ctx context.Context, mi *MovementInstance) error
	DeleteMovementInstance(ctx context.Context, id int) error
}
