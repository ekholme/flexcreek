package flexcreek

import "time"

//a movement is a essentially an exercise, e.g. a kettlebell swing
type Movement struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Muscles []string `json:"muscles"`

	//timestamps for creation and update
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MovementService interface {
	CreateMovement(m *Movement) (int, error) //return the id number
	GetMovementByID(id int) (*Movement, error)
	GetMovementByName(name string) (*Movement, error)
	GetAllMovements() ([]*Movement, error)
	DeleteMovement(id int) (int, error)
	//TODO ADD OTHER METHODS
}
