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
	var sortOption int

	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch tasks with optional state filter and sorting",
		Long: `Fetch all tasks, or filter by state.

Valid states:
  0: Pending
  1: Working
  2: Finished

Sorting Options:
  --sort 1  -> Newest First
  --sort 2  -> Shortest Text First
  --sort 3  -> Longest Text First`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var state *int // Pointer to track whether a state was provided
			if len(args) > 0 {
				s, err := strconv.Atoi(args[0])
				if err != nil || s < 0 || s > service.FinishedState {
					return fmt.Errorf("invalid state: %v. Must be between 0 and %d", s, service.FinishedState)
				}
				state = &s
			}

			// Fetch tasks with sorting
			return fetchTasks(db, state, sortOption)
		},
	}

	// Add sorting flag
	cmd.Flags().IntVar(&sortOption, "sort", 0, "Sort tasks (1=Newest, 2=Shortest, 3=Longest)")

	return cmd
}

func fetchTasks(db *database.Database, state *int, sortOption int) error {
	var tasks []database.Task
	var err error

	if state != nil {
		tasks, err = db.GetTaskByStateSorted(*state, sortOption)
	} else {
		tasks, err = db.GetAllTasksSorted(sortOption) // ✅ Apply sorting for all tasks
	}

	if err != nil {
		return fmt.Errorf("error fetching tasks: %v", err)
	}

	if state != nil {
		fmt.Printf("State: %s\n\n\n", []string{"PENDING", "WORKING", "FINISHED"}[*state])
		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.TEXT)
		}
		return nil
	}

	// Column-wise print for all states
	fetchedTasks := make([][]database.Task, service.FinishedState+1)
	for _, task := range tasks {
		fetchedTasks[task.STATE] = append(fetchedTasks[task.STATE], task)
	}

	maxLen := 0
	for _, t := range fetchedTasks {
		if len(t) > maxLen {
			maxLen = len(t)
		}
	}

	fmt.Println("\tPENDING \t\t\tWORKING \t\t\tFINISHED")

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

func shortenString(s string, width int) string {
	if len(s) > width {
		return fmt.Sprintf("%-*s", width-3, s[:width-3]) + "..."
	}
	return fmt.Sprintf("%-*s", width, s) // Ensures alignment
}
