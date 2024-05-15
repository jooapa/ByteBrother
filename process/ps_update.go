package process

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	bt "bigbro/bigtime"
	fl "bigbro/filer"
)

// ProcessEntry represents a single entry in the process log.
type ProcessEntry struct {
	Name             string   `json:"name"`
	TotalTime        string   `json:"total_time"`
	CurrentlyRunning bool     `json:"currently_opened"`
	Opened           []string `json:"opened"`
	Closed           []string `json:"closed"`
}

// Processes updates the process log.
func Processes(forceNotRunning bool) {
	currentProcesses := GetProcesses()
	logEntries := readLogFile()

	// Create a map for quick lookup of current processes
	currentProcessesMap := make(map[string]bool)
	for _, process := range currentProcesses {
		currentProcessesMap[process] = true
	}

	// Update or add new process entries
	for _, process := range currentProcesses {
		found := false
		for i, entry := range logEntries {
			if entry.Name == process {
				newTime := getTimeDuration(entry.TotalTime)
				newTime = bt.AddTimeTogether(newTime, bt.ElapsedTime())
				logEntries[i].TotalTime = newTime.String()
				logEntries[i].CurrentlyRunning = true
				if !entry.CurrentlyRunning {
					// If the process was not running before, add a timestamp to the Opened array
					logEntries[i].Opened = append(logEntries[i].Opened, time.Now().Format(time.RFC3339))
				}
				found = true
				break
			}
		}
		if !found {
			logEntries = append(logEntries, ProcessEntry{
				Name:             process,
				TotalTime:        bt.ElapsedTime().String(),
				CurrentlyRunning: true,
				Opened:           []string{time.Now().Format(time.RFC3339)}, // Add a timestamp to the Opened array
			})
		}
	}

	// Check for processes that have stopped running
	for i, entry := range logEntries {
		if _, exists := currentProcessesMap[entry.Name]; !exists || forceNotRunning {
			logEntries[i].CurrentlyRunning = false // Set the CurrentlyRunning field to false
			// If the process was running before, add a timestamp to the Closed array
			if entry.CurrentlyRunning || forceNotRunning && entry.CurrentlyRunning {
				logEntries[i].Closed = append(logEntries[i].Closed, time.Now().Format(time.RFC3339))
			}
		}
	}

	// Write the new log entries to the file
	if makeLogEntries(logEntries) {
		return
	}
}

func readLogFile() []ProcessEntry {
	var logEntries []ProcessEntry
	if fl.IfFileExists(fl.Folder + "/" + fl.ExeLog) {
		oldContent, err := fl.ReadFilePath(fl.Folder + "/" + fl.ExeLog)
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
		} else {

			err := json.Unmarshal(oldContent, &logEntries)
			if err != nil {
				fmt.Printf("Failed to unmarshal existing log: %v\n", err)
			}
		}
	} else {
		fmt.Println("File doesn't exist")
	}
	return logEntries
}

func makeLogEntries(logEntries []ProcessEntry) bool {
	newContent, err := json.MarshalIndent(logEntries, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal log entries: %v\n", err)
		return true
	}

	err = os.WriteFile(fl.Folder+"/"+fl.ExeLog, newContent, 0644)
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
	}
	return false
}

// getTimeDuration converts a duration string to a time.Duration object.
func getTimeDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		fmt.Printf("Failed to parse duration: %v\n", err)
		return 0
	}
	return duration
}

func ResetCurrentlyOpened() {
	// Read the existing log file
	logEntries := readLogFile()

	// Set all currently_opened fields to false
	for i := range logEntries {
		logEntries[i].CurrentlyRunning = false
	}

	// Write the updated log entries back to the file
	if makeLogEntries(logEntries) {
		return
	}
}
