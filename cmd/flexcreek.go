package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek/sqlite"
	"github.com/ekholme/flexcreek/ui"
	_ "modernc.org/sqlite"
)

const (
	logFilePath = "logs/demo_log.json"
	dsn         = "./flexcreek.db"
)

func main() {

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("couldn't open database: %s", err.Error())
	}
	defer db.Close()

	userService := sqlite.NewUserService(db)
	workoutService := sqlite.NewWorkoutService(db)

	m := ui.New(userService, workoutService)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nFlexcreek exiting.")
}
