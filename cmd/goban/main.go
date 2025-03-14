package goban

import (
	"goban/internal/database"
	"goban/test"
	"log"
	"os"
)

func RunGoban() {
	db, err := database.InitDB("database.db")
	if err != nil {
		log.Fatalf("Unable to initialize the database: %v", err)
	}
	defer db.Close()

	test.RunTests(db)

	if len(os.Args) > 1 {
		gobancli := NewApp(db)
		gobancli.Execute()
	} else {
		RunTUI()
	}
}
