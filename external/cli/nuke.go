package cli

import (
	"fmt"
	"goban/internal/database"
	"log"

	"github.com/spf13/cobra"
)

func NukeCommand(db *database.Database) *cobra.Command {
	nukeCmd := &cobra.Command{
		Use:   "nuke",
		Short: "Nuke the database",
		Long:  `Clears the entire task table, and removes all id counters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := db.NukeDB()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Database nuked successfully!")
			return nil
		},
	}
	return nukeCmd
}
