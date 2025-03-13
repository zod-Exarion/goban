package test

import (
	"fmt"
	"goban/internal/database"
	"goban/internal/service"
	"log"
)

func RunTests(db *database.Database) {
	// PrintTasks(db)
	// db.NukeDB()
	// CreateDummyTasks(db)
	// PrintTasks(db)
}

func CreateDummyTasks(db *database.Database) {
	for i := 0; i < 10; i++ {
		task := service.CreateTask("Build database", i%3)

		err := db.SaveTask(&task, true)
		if err != nil {
			log.Printf("Error saving task: %v", err)
		}
	}
}

func PrintTasks(db *database.Database) {
	for i := 0; i < 3; i++ {
		tasks, err := db.GetTaskByState(i)
		if err != nil {
			log.Fatalf("Error retrieving tasks: %v", err)
		}

		fmt.Printf("\n%v State\n", i)
		for j := range tasks {
			fmt.Println(tasks[j])
		}
	}

	fmt.Println("\nFinsihed Fetching")
}
