package main

import (
	a "bytebrother/main/archive"
	bt "bytebrother/main/bigtime"
	cp "bytebrother/main/clipboard"
	fl "bytebrother/main/filer"
	h "bytebrother/main/hook"
	nt "bytebrother/main/network"
	ps "bytebrother/main/process"
	st "bytebrother/main/settings"
	f "fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	t "time"
)

var processInterval t.Duration = 1000
var processSaving = false
var processSaveInterval t.Duration = 5000

var keylogging = false
var clipboardTracking_text = false
var clipboardTracking_image = false

func main() {
	fl.MakeNecessaryFiles()
	ps.ResetCurrentlyOpened()
	settings := st.LoadSettings()
	processInterval = t.Duration(settings.ProcessInterval)
	processSaveInterval = t.Duration(settings.SaveProcessInformationInFile)

	nt.ChosenIndex = settings.NetworkIndexToMonitor
	a.ArchiveRowCount = settings.NumRowsInArchive

	keylogging = settings.KeyloggerEnabled
	clipboardTracking_text = settings.ClipboardTextTrackingEnabled
	clipboardTracking_image = settings.ClipboardImageTrackingEnabled

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive interrupt and termination signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run a goroutine that waits for the application to be interrupted or terminated
	go func() {
		<-sigs
		ps.Processes(true, true) // Process one more time to get the timestamp
		os.Exit(0)
	}()

	shouldArchive := a.ShouldArchive(fl.ProcessFolder + fl.ExeLog)

	if shouldArchive {
		err := a.Archive(fl.ProcessFolder + fl.ExeLog)
		if err != nil {
			f.Printf("Error archiving the file: %v\n", err)
		}
	}

	f.Println("Byte Brother is watching you")

	// PROCESS
	go func() {
		for {
			if nt.ChosenIndex != 69420 {

				if processSaving {
					ps.Processes(false, true)
					processSaving = false
				} else {
					ps.Processes(false, false)
				}
				bt.Reset()

				// Sleep for 10 seconds
				t.Sleep(processInterval * t.Millisecond)
			}
		}
	}()

	// PROCESS SAVING
	go func() {
		for {
			if nt.ChosenIndex != 69420 {
				t.Sleep(processSaveInterval * t.Millisecond)
				processSaving = true
			}
		}
	}()

	// NETWORK
	go func() {
		// if couldn't load wpcap.dll, dont start network
		err := nt.Start()
		if err != nil {
			log.Fatalf("Failed to start network: %v", err)
		}
	}()

	if keylogging {
		// GLOBAL HOOK
		go func() {
			h.StartHook()
		}()
	}

	// CLIPBOARD Image
	if clipboardTracking_image {
		go func() {
			cp.Image()
		}()
	}

	// CLIPBOARD Text
	if clipboardTracking_text {
		go func() {
			cp.Text()
		}()
	}

	select {}
}
