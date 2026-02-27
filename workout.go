package flexcreek

import "time"

type Workout struct {
	ID               int       `db:"id"`
	UserID           int       `db:"user_id"`
	ShortDescription string    `db:"short_description"`
	LongDescription  string    `db:"long_description"`
	WorkoutDate      time.Time `db:"workout_date"`
	CreatedAt        time.Time `db:"created_at"`
}
