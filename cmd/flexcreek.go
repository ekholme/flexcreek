package main

import (
	"database/sql"
	"log"

	"github.com/ekholme/flexcreek/sqlite"
)

const (
	logFilePath = "logs/demo_log.json"
	dsn         = "./flexcreek.db"
)

func main() {

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("couldn't open database: %s", err.Error())
	}

	userService := sqlite.NewUserService(db)
	workoutService := sqlite.NewWorkoutService(db)
	activityTypeService := sqlite.NewActivityTypeService(db)

}
