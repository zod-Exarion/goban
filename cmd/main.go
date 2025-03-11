package main

import (
	"fmt"
	"goban/cmd/goban"
	"log"
	"os"
)

func main() {
	// HACK: Can't use stdOut due to occupation by Bubble Tea
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer func() {
		logFile.Close()

		fileInfo, err := os.Stat("debug.log")
		if err == nil && fileInfo.Size() == 0 {
			os.Remove("debug.log")
		} else {
			fmt.Println("Errors stored in debug.log file")
		}
	}()

	log.SetOutput(logFile)

	goban.RunGoban()
}
