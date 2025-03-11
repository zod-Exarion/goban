package goban

import (
	"fmt"
	"goban/external/cli"
	"goban/internal/database"
	"os"

	"github.com/spf13/cobra"
)

type App struct {
	RootCmd *cobra.Command
	DB      *database.Database
}

func NewApp(db *database.Database) *App {
	app := &App{
		DB: db,
		RootCmd: &cobra.Command{
			Use:   "goban",
			Short: "Kanban CLI/TUI in Go",
			Long:  `Goban is Kanban Task Manager with both a CLI and TUI twist!`,
		},
	}

	app.addCommands()

	return app
}

func (a *App) addCommands() {
	a.RootCmd.AddCommand(cli.FetchCommand(a.DB))
}

func (app *App) Execute() {
	if err := app.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
