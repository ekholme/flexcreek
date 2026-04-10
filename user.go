package flexcreek

import (
	"time"
)

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"create_at"`
}
