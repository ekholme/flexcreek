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
	dsn               = "flexcreek.db"
	workoutListLength = 10
	testingID         = 1
)

func main() {

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("Couldn't open the database: %s", err)
	}

	defer db.Close()

	storage := sqlite.NewStorage(db)
	// userModel := ui.NewUserModel(storage)
	workoutModel := ui.NewWorkoutModel(storage, testingID, workoutListLength)
	p := tea.NewProgram(workoutModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
