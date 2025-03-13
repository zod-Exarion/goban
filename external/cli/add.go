package cli

import (
	"fmt"
	"goban/internal/database"
	"goban/internal/service"
	"log"

	"github.com/spf13/cobra"
)

func AddTaskCommand(db *database.Database) *cobra.Command {
	var text string
	var state int

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a task",
		Long:  `Add a new task with a specific text and state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if state < 0 || state > service.FinishedState {
				return fmt.Errorf("invalid state: %d. Must be between 0 and %d", state, service.FinishedState)
			}

			task := service.CreateTask(text, state)
			if err := db.SaveTask(&task, false); err != nil {
				log.Fatal(err)
			}

			fmt.Println("Task added successfully!")
			return nil
		},
	}

	// Define flags
	addCmd.Flags().StringVar(&text, "text", "", "Task description")
	addCmd.Flags().IntVar(&state, "state", 0, "Task state (0=Pending, 1=Working, 2=Finished)")

	// Mark "text" flag as required
	addCmd.MarkFlagRequired("text")

	return addCmd
}
