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
			var state *int
			if len(args) > 0 {
				s, err := strconv.Atoi(args[0])
				if err != nil || s < 0 || s > service.FinishedState {
					return fmt.Errorf("invalid state: %v. Must be between 0 and %d", s, service.FinishedState)
				}
				state = &s
			}
			return fetchTasks(db, state, sortOption)
		},
	}

	cmd.Flags().IntVar(&sortOption, "sort", 0, "Sort tasks (1=Newest, 2=Shortest, 3=Longest)")
	return cmd
}

func fetchTasks(db *database.Database, state *int, sortOption int) error {
	var tasks []database.Task
	var err error

	if state != nil {
		tasks, err = db.GetTaskByStateSorted(*state, sortOption)
	} else {
		tasks, err = db.GetAllTasksSorted(sortOption)
	}

	if err != nil {
		return err
	}

	if state != nil {
		fmt.Printf("State: %s\n\n\n", []string{"PENDING", "WORKING", "FINISHED"}[*state])
		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.TEXT)
		}
		return nil
	}

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
		for state := 0; state <= service.FinishedState; state++ {
			tasks := fetchedTasks[state]
			if i < len(tasks) {
				fmt.Printf("%d. %s \t\t", tasks[i].ID, shortenString(tasks[i].TEXT, shortenBy))
			} else {
				if state < service.FinishedState {
					fmt.Print("\t\t\t")
				}
			}
		}
		fmt.Println()
	}

	return nil
}

func shortenString(s string, width int) string {
	if len(s) > width {
		return fmt.Sprintf("%-*s", width-3, s[:width-3]) + "..."
	}
	return fmt.Sprintf("%-*s", width, s)
}
