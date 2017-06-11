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

type Change struct {
	Error error
	// TODO: expose more info than just a message stating an update
	Message string
}

func ErrorChange(err error) Change {
	return Change{Error: err, Message: ""}
}
func SuccessfulChange(message string) Change {
	return Change{Error: nil, Message: message}
}

func Poll(delay time.Duration, name string) <-chan Change {
	// TODO: should we also accept quit channel as an arg?
	changes := make(chan Change, 0)
	go pollForChanges(delay, name, changes)
	return changes
}

func pollForChanges(delay time.Duration, name string, changes chan Change) {
	var (
		lastModified time.Time
		lastSize     int64
	)

	for {
		stat, err := os.Stat(name)

		// TODO: make it nicer somehow?
		if err != nil {
			// TODO: Break if it happens too often?
			changes <- ErrorChange(err)
		} else {
			isNew := stat.ModTime().After(lastModified)
			sizeChanged := stat.Size() != lastSize

			if isNew || sizeChanged {
				lastModified = stat.ModTime()
				lastSize = stat.Size()

				changes <- SuccessfulChange("Change occured")
			}
		}

		select {
		case <-time.After(delay):
		}
	}
}
