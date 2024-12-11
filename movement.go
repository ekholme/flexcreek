package flexcreek

//a movement is a essentially an exercise, e.g. a kettlebell swing
type Movement struct {
	ID      int
	Name    string
	Muscles []string
}

type MovementService interface {
	CreateMovement(name string, muscles []string) int //return the id number
	GetMovementByID(id int) *Movement
	GetMovementByName(name string) *Movement
	DeleteMovement(id int)
	//TODO ADD OTHER METHODS
}
