package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	userModel := userselect.New(userService)

	p := tea.NewProgram(userModel, tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := m.(userselect.Model); ok {
		if m.SelectedUser != nil {
			log.Printf("User %s selected.", m.SelectedUser.Username)
		}
	}

	fmt.Println("\nFlexcreek exiting.")
}
