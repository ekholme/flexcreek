package flexcreek

import "time"

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
	//TODO
}
