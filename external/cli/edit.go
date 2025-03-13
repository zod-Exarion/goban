package cli

import (
	"fmt"
	"goban/internal/database"
	"log"

	"github.com/spf13/cobra"
)

func EditTaskCommand(db *database.Database) *cobra.Command {
	var id int
	var text string

	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a task",
		Long:  `Edit a task with a specific id providing text.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := db.EditTask(id, text)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Task edited successfully!")
			return nil
		},
	}

	// Define flags
	editCmd.Flags().StringVar(&text, "text", "", "Task description")
	editCmd.Flags().IntVar(&id, "id", 0, "ID of the Task")

	// Mark "text" flag as required
	editCmd.MarkFlagRequired("id")
	editCmd.MarkFlagRequired("text")

	return editCmd
}
