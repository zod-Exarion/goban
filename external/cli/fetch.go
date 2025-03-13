package cli

import (
	"fmt"
	"goban/internal/database"
	"goban/internal/service"
	"strconv"

	"github.com/spf13/cobra"
)

const shortenBy = 15

func FetchCommand(db *database.Database) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch [state]",
		Short: "Fetch tasks with optional state filter",
		Long: `Fetch all tasks, or filter by state.

Valid states:
  0: Pending
  1: Working
  2: Finished`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				state, err := strconv.Atoi(args[0])
				if err != nil || state < 0 || state > service.FinishedState {
					return fmt.Errorf("invalid state: %v. Must be between 0 and %d", state, service.FinishedState)
				}
				return fetchTasksByState(db, state)
			}
			return fetchTasks(db)
		},
	}
}

func fetchTasks(db *database.Database) error {
	fetchedTasks := make([][]database.Task, service.FinishedState+1)

	for i := 0; i <= service.FinishedState; i++ {
		tasks, err := db.GetTaskByState(i)
		if err != nil {
			return fmt.Errorf("error fetching tasks: %v", err)
		}
		fetchedTasks[i] = tasks
	}

	maxLen := 0
	for _, tasks := range fetchedTasks {
		if len(tasks) > maxLen {
			maxLen = len(tasks)
		}
	}

	// HACK: Header for printing
	fmt.Println("\tPENDING \t\t\tWORKING \t\t\tFINISHED")

	// INFO: Column format printing
	for i := 0; i < maxLen; i++ {
		if i < len(fetchedTasks[service.PendingState]) {
			t := fetchedTasks[service.PendingState][i]
			fmt.Printf("%d. %s \t\t", t.ID, shortenString(t.TEXT, shortenBy))
		} else {
			fmt.Print("\t\t\t")
		}

		if i < len(fetchedTasks[service.WorkingState]) {
			t := fetchedTasks[service.WorkingState][i]
			fmt.Printf("%d. %s \t\t", t.ID, shortenString(t.TEXT, shortenBy))
		} else {
			fmt.Print("\t\t\t")
		}

		if i < len(fetchedTasks[service.FinishedState]) {
			t := fetchedTasks[service.FinishedState][i]
			fmt.Printf("%d. %s", t.ID, shortenString(t.TEXT, shortenBy))
		}

		fmt.Println()
	}

	return nil
}

func fetchTasksByState(db *database.Database, state int) error {
	stateNames := []string{"PENDING", "WORKING", "FINISHED"}
	fmt.Printf("State: %s\n\n\n", stateNames[state])

	tasks, err := db.GetTaskByState(state)
	if err != nil {
		return fmt.Errorf("error fetching tasks by state %d: %v", state, err)
	}

	for _, task := range tasks {
		fmt.Printf("%d. %s\n", task.ID, task.TEXT)
	}

	return nil
}

// PERF: Shorten task text for printing preview
func shortenString(s string, width int) string {
	if len(s) > width {
		return fmt.Sprintf("%-*s", width-3, s[:width-3]) + "..."
	}
	return fmt.Sprintf("%-*s", width, s) // Ensures alignment
}
