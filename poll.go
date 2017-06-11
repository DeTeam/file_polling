package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("Work in progress")

	changes := Poll(1*time.Second, "test")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case r := <-changes:
			if r.Error != nil {
				fmt.Println("Got some errors", r.Error)
			} else {
				fmt.Println("Got something: ", r.Message)
			}

		case <-signalChan:
			fmt.Println("Interrupted, quiting")
			return
		}
	}

}

type Result struct {
	Error   error
	Message string
}

func ErrorResult(err error) Result {
	return Result{Error: err, Message: ""}
}
func SuccessfulResult(message string) Result {
	return Result{Error: nil, Message: message}
}

func Poll(d time.Duration, name string) <-chan Result {
	changes := make(chan Result, 0)
	go pollForChanges(d, name, changes)
	return changes
}

func pollForChanges(d time.Duration, name string, changes chan Result) {
	var (
		lastModified time.Time
		lastSize     int64
	)

	for {
		stat, err := os.Stat(name)

		if err != nil {
			changes <- ErrorResult(err)
		} else {
			isNew := stat.ModTime().After(lastModified)
			sizeChanged := stat.Size() != lastSize

			if isNew || sizeChanged {
				lastModified = stat.ModTime()
				lastSize = stat.Size()

				changes <- SuccessfulResult("Change occured")
			}
		}

		select {
		case <-time.After(d):
			fmt.Println("Done waiting")
		}
	}
}
