package cli

import (
	"fmt"
	"goban/internal/database"
	"goban/internal/service"
	"strconv"

	"github.com/spf13/cobra"
)

var state int

func FetchCommand(db *database.Database) *cobra.Command {
	FetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch all tasks",
		Long: `Fetch all tasks.

You can optionally filter tasks by state by providing a state number as an argument.
Valid states are:
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

				fetchTasksByState(db, state)
			} else {
				fetchTasks(db)
			}

			return nil
		},
	}

	return FetchCmd
}

func fetchTasks(db *database.Database) error {
	var err error
	fetchedTasks := make([][]database.Task, service.FinishedState+1)

	for i := 0; i <= service.FinishedState; i++ {
		fetchedTasks[i], err = db.GetTaskByState(i)
	}

	lengths := make([]int, len(fetchedTasks))
	for i, tasks := range fetchedTasks {
		lengths[i] = len(tasks)
	}

	maxLen := lengths[0]
	for _, length := range lengths[1:] {
		if length > maxLen {
			maxLen = length
		}
	}

	fmt.Println("\tPENDING \t\t\tWORKING \t\t\tFINISHED")
	shortenBy := 15
	for i := 0; i < maxLen; i++ {
		if i < len(fetchedTasks[service.PendingState]) {
			t0 := fetchedTasks[service.PendingState][i]
			fmt.Printf("%v. %v \t\t", t0.ID, shortenString(t0.TEXT, shortenBy))
		}

		if i < len(fetchedTasks[service.WorkingState]) {
			t1 := fetchedTasks[service.WorkingState][i]
			fmt.Printf("%v. %v \t\t", t1.ID, shortenString(t1.TEXT, shortenBy))
		}

		if i < len(fetchedTasks[service.FinishedState]) {
			t2 := fetchedTasks[service.FinishedState][i]
			fmt.Printf("%v. %v \t\t", t2.ID, shortenString(t2.TEXT, shortenBy))
		}

		fmt.Println()
	}

	return err
}

func fetchTasksByState(db *database.Database, state int) error {
	var currentstate string

	switch state {
	case service.PendingState:
		currentstate = "PENDING"
	case service.WorkingState:
		currentstate = "WORKING"
	case service.FinishedState:
		currentstate = "FINISHED"
	}
	fmt.Println("State:", currentstate)

	tasks, err := db.GetTaskByState(state)
	if err != nil {
		return fmt.Errorf("error getting tasks by state %v: %v", state, err)
	}

	for _, task := range tasks {
		fmt.Printf("%v. %v\n", task.ID, task.TEXT)
	}

	return nil
}

func shortenString(s string, shortenWidth int) string {
	if len(s) > shortenWidth {
		s = s[:shortenWidth] + "..."
	}

	return s
}
