package cli

import (
	"goban/internal/database"

	"github.com/spf13/cobra"
)

func AddTaskCommand(db *database.Database) *cobra.Command {
	AddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a task",
		Long:  `Add a task.`,
		Args:  cobra.ExactArgs(3),

		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return AddCmd
}
