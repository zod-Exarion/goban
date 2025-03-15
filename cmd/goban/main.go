package goban

import (
	"goban/internal/database"
	"log"
	"os"
	"path/filepath"
)

func RunGoban() {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Resolve the executable path to its absolute location
	realpath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		panic(err)
	}

	// Get the directory of the realpath
	dir := filepath.Dir(realpath)

	db, err := database.InitDB(filepath.Join(dir, "database.db"))
	if err != nil {
		log.Fatalf("Unable to initialize the database: %v", err)
	}
	defer db.Close()

	if len(os.Args) > 1 {
		gobancli := NewApp(db)
		gobancli.Execute()
	} else {
		RunTUI()
	}
}
