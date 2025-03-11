package goban

import (
	"goban/internal/database"
	"goban/test"
	"log"
	"os"
)

func RunGoban() {
	// HACK: Can't use stdOut due to occupation by Bubble Tea
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

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
