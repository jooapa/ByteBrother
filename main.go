package main

import (
	bt "bytebrother/main/bigtime"
	fl "bytebrother/main/filer"
	ps "bytebrother/main/process"
	f "fmt"
	"os"
	"os/signal"
	"syscall"
	t "time"
)

var processInterval t.Duration = 400

func main() {
	fl.MakeNecessaryFiles()

	ps.ResetCurrentlyOpened()

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive interrupt and termination signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run a goroutine that waits for the application to be interrupted or terminated
	go func() {
		<-sigs
		ps.Processes(true) // Process one more time to get the timestamp
		os.Exit(0)
	}()

	// Run the function in the background
	go func() {
		for {
			// Call your first function here
			ps.Processes(false)
			f.Printf("Current time: %v\n", bt.ElapsedTime())
			bt.Reset()

			// Sleep for 10 seconds
			t.Sleep(processInterval * t.Millisecond)
		}
	}()

	// Keep the main function running indefinitely
	select {}
}
