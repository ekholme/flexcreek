package flexcreek

import "time"

type Muscle struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	//timestamps for creation and update
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MuscleService interface {
	CreateMuscle(muscle *Muscle) (int, error)
	GetMuscleByID(id int) (*Muscle, error)
	GetMuscleByName(name string) (*Muscle, error)
	GetAllMuscles() ([]*Muscle, error)
	DeleteMuscle(id int) (int, error)
	//TODO ADD OTHER METHODS
}
