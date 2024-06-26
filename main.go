package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	t "time"

	"github.com/gofrs/flock"

	a "bytebrother/main/archive"
	bt "bytebrother/main/bigtime"
	cp "bytebrother/main/clipboard"
	fl "bytebrother/main/filer"
	g "bytebrother/main/global"
	h "bytebrother/main/hook"
	nt "bytebrother/main/network"
	ps "bytebrother/main/process"
	ss "bytebrother/main/screenshot"
	st "bytebrother/main/settings"
	f "fmt"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "setup" {
			g.IsSetup = true
		}
	}

	fl.MakeDir(fl.Folder)
	// Create a lock file to prevent multiple instances of the application from running
	// Create a new file lock
	lock := flock.New(fl.Folder + fl.LockFile)
	locked, err := lock.TryLock()
	if err != nil {
		log.Fatalf("Failed to lock: %v", err)
	}
	if !locked {
		log.Fatalf("Another instance of the application is already running")
	}

	fl.MakeNecessaryFiles()
	ps.ResetCurrentlyOpened()
	settings := st.LoadSettings()

	g.ProcessInterval = t.Duration(settings.ProcessInterval)
	g.SaveProcessInformationInFile = t.Duration(settings.SaveProcessInformationInFile)
	g.ScreenshotInterval = t.Duration(settings.ScreenshotInterval)

	g.ChosenIndex = settings.NetworkIndexToMonitor
	g.ArchiveRowCount = settings.NumRowsInArchive

	g.Keylogging = settings.KeyloggerEnabled
	g.ClipboardTracking_text = settings.ClipboardTextTrackingEnabled
	g.ClipboardTracking_image = settings.ClipboardImageTrackingEnabled
	g.CanTakeScreenshot = settings.CanTakeScreenshot

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive interrupt and termination signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run a goroutine that waits for the application to be interrupted or terminated
	go func() {
		<-sigs

		lock.Unlock()
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
			if g.ProcessSaving {
				ps.Processes(false, true)
				g.ProcessSaving = false
			} else {
				ps.Processes(false, false)
			}
			bt.Reset()

			// Sleep for 10 seconds
			t.Sleep(g.ProcessInterval * t.Millisecond)
		}
	}()

	// PROCESS SAVING
	go func() {
		for {
			t.Sleep(g.ProcessSaveInterval * t.Millisecond)
			g.ProcessSaving = true
		}
	}()

	// NETWORK
	go func() {
		// if couldn't load wpcap.dll, dont start network
		err := nt.Start()
		if err != nil {
			log.Fatalf("Failed to start network: %v\nDownload WinPcap from https://www.winpcap.org/install/default.htm", err)
		}
	}()

	if g.Keylogging {
		// GLOBAL HOOK
		go func() {
			h.StartHook()
		}()
	}

	// CLIPBOARD Image
	if g.ClipboardTracking_image {
		go func() {
			cp.Image()
		}()
	}

	// CLIPBOARD Text
	if g.ClipboardTracking_text {
		go func() {
			cp.Text()
		}()
	}

	// SCREENSHOT
	if g.CanTakeScreenshot {
		go func() {
			for {
				// ARCHIVE SCREENSHOTS IF NEEDED
				// if g.ArchiveScreenshotsAfter {
				// can, folderName := ss.CanArchiveOlderFolder()
				// if can {
				// 	err := a.ArchiveFolder_sevenzip(fl.ScreenshotFolder + folderName)
				// 	if err != nil {
				// 		f.Printf("Error archiving the folder: %v\n", err)
				// 		os.Remove(fl.ScreenshotFolder + folderName + ".7z")
				// 		os.Rename(fl.ScreenshotFolder+folderName, fl.ScreenshotFolder+folderName+"_archived_failed")
				// 	}
				// }
				// }

				ss.TakeScreenshot()
				t.Sleep(g.ScreenshotInterval_sec * t.Second)
			}
		}()
	}

	select {}
}
