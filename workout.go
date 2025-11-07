package flexcreek

type Workout struct {
	//todo
}

type WorkoutService interface {
	//todo
	Latest() ([]*Workout, error) //get the X latest workouts
}
