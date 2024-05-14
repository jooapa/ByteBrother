package main

import (
	bt "bigbro/bigtime"
	fl "bigbro/filer"
	ps "bigbro/process"
	f "fmt"
	t "time"
)

var processInterval t.Duration = 1

func main() {
	// time

	fl.MakeNessesaryFiles()

	// Run the function in the background
	go func() {
		for {
			// Call your function here
			ps.Processes()
			f.Printf("Current time: %v\n", bt.ElapsedTime())
			bt.Reset()

			// Sleep for 10 seconds
			t.Sleep(processInterval * t.Second)
		}
	}()

	// Keep the main goroutine alive
	// This is necessary to prevent the program from exiting immediately
	// since the other goroutine is running in the background.
	select {}

}
