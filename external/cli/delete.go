package cli

import (
	"fmt"
	"goban/internal/database"
	"log"

	"github.com/spf13/cobra"
)

func DeleteTaskCommand(db *database.Database) *cobra.Command {
	var id int

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a task",
		Long:  `Delete a task using the corresponding id from the task list.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := db.DeleteTask(id)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Task deleted successfully.")
			return nil
		},
	}

	deleteCmd.Flags().IntVar(&id, "id", 0, "Task ID")
	deleteCmd.MarkFlagRequired("id")

	return deleteCmd
}
