package service

import (
	"goban/internal/database"
)

const (
	PendingState = iota
	WorkingState
	FinishedState
)

func CreateTask(text string, state uint) database.Task {
	return database.Task{
		TEXT:  text,
		STATE: state,
	}
}
