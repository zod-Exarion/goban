package goban

import (
	"goban/internal/database"
	"goban/test"
	"log"
	"os"
)

func RunGoban() {
	// TODO: Proper location for database on Linux file system
	db, err := database.InitDB("database.db")
	if err != nil {
		log.Fatalf("Unable to initialize the database: %v", err)
	}
	defer db.Close()

	// FIX: Remove Tests before final release
	test.RunTests(db)

	if len(os.Args) > 1 {
		gobancli := NewApp(db)
		gobancli.Execute()
	} else {
		RunTUI()
	}
}
