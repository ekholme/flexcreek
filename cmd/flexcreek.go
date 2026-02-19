package main

import (
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek/sqlite"
	"github.com/ekholme/flexcreek/ui/userselect"
	_ "modernc.org/sqlite"
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

	defer db.Close()

	userService := sqlite.NewUserService(db)

	userModel := userselect.New(userService)

	p := tea.NewProgram(userModel)

}
