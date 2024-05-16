package main

import (
	a "bytebrother/main/archive"
	bt "bytebrother/main/bigtime"
	fl "bytebrother/main/filer"
	nt "bytebrother/main/network"
	ps "bytebrother/main/process"
	st "bytebrother/main/settings"
	f "fmt"
	"os"
	"os/signal"
	"syscall"
	t "time"
)

var processInterval t.Duration = 1000

func main() {
	fl.MakeNecessaryFiles()
	ps.ResetCurrentlyOpened()
	settings := st.LoadSettings()
	processInterval = t.Duration(settings.ProcessInterval)
	nt.ChosenIndex = settings.NetworkIndexToMonitor
	a.ArchiveRowCount = settings.NumRowsInArchive

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

	f.Println("Starting the main loop...")
	shouldArchive := a.ShouldArchive(fl.LogFolder + fl.ExeLog)

	if shouldArchive {
		err := a.Archive(fl.LogFolder + fl.ExeLog)
		if err != nil {
			f.Printf("Error archiving the file: %v\n", err)
		}
	}

	go func() {

		for {
			if nt.ChosenIndex != 69420 {
				// Call your first function here
				ps.Processes(false)
				f.Printf("Current time: %v\n", bt.ElapsedTime())
				bt.Reset()

				// Sleep for 10 seconds
				t.Sleep(processInterval * t.Millisecond)
			}
		}
	}()

	go func() {
		nt.Start()
	}()

	select {}
}
