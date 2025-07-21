package flexcreek

import "time"

type MovementInstance struct {
	ID         int           `json:"id"`
	UserID     int           `json:"userId"`
	WorkoutID  int           `json:"workoutId"`
	Movement   *Movement     `json:"movement"`
	Notes      string        `json:"notes"`
	Sets       int           `json:"sets"`       //for workouts with sets
	RepsPerSet int           `json:"repsPerSet"` //for workouts with sets
	Duration   time.Duration `json:"duration"`   //for workouts with time
	Distance   float64       `json:"distance"`   //for workout with distance
	RPE        int           `json:"rpe"`        //relative perceived exertion (1-10)
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

type MovementInstanceService interface {
	//TODO
}
