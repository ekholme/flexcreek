package main

import (
	"database/sql"
	"log"

	"github.com/ekholme/flexcreek/server"
	"github.com/ekholme/flexcreek/sqlite"
)

const (
	logFilePath   = "logs/demo_log.json"
	templatePaths = "templates/*.html"
	dsn           = "./flexcreek.db"
	addr          = ":8080"
)

func main() {

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("couldn't open database: %s", err.Error())
	}

	userService := sqlite.NewUserService(db)
	workoutService := sqlite.NewWorkoutService(db)
	activityTypeService := sqlite.NewActivityTypeService(db)

	s := server.NewServer(addr, userService, workoutService, activityTypeService, templatePaths, logFilePath)
	s.Run()
}
