package cli

import (
	"fmt"
	"goban/internal/database"
	"log"

	"github.com/spf13/cobra"
)

func MarkTaskCommand(db *database.Database) *cobra.Command {
	var id int

	markCmd := &cobra.Command{
		Use:   "mark",
		Short: "Mark a task",
		Long:  `Mark a task with a specific id .`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := db.MarkTask(id)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Task marked successfully!")
			return nil
		},
	}

	// Define flags
	markCmd.Flags().IntVar(&id, "id", 0, "ID of the Task")

	// Mark "text" flag as required
	markCmd.MarkFlagRequired("id")

	return markCmd
}
