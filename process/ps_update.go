package process

import (
	bt "bigbro/bigtime"
	fl "bigbro/filer"
	f "fmt"
	"strings"
	t "time"
)

func Processes() {
	processes := GetProcesses()
	var newContent string
	var oldContent string

	// read the file to get the old content
	if fl.IfFileExists(fl.Folder + "/" + fl.ExeLog) {
		oldContent2, err := fl.ReadFilePath(fl.Folder + "/" + fl.ExeLog)
		if err != nil {
			f.Printf("Failed to read file: %v\n", err)
		}
		oldContent = string(oldContent2) // Convert []byte to string
	} else {
		f.Printf("File doesn't exist\n")
	}

	oldContentSlice := strings.Split(string(oldContent), "\n")

	// if process is already in the file, update the time, else add the process and related time
	updatedProcesses := make([]string, 0, len(oldContentSlice))

	for i := 0; i < len(processes); i++ {
		found := false
		for j := 0; j < len(oldContentSlice); j++ {
			// if the process is already in the file add the elapsed time to the time
			if strings.Contains(oldContentSlice[j], processes[i]) {
				newTime := getTimeFromProcess(oldContentSlice[j])
				newTimeDuration, err := t.ParseDuration(newTime)
				if err != nil {
					f.Printf("Failed to parse duration: %v\n", err)
				}
				newTimeDuration = bt.AddTimeTogether(newTimeDuration, bt.ElapsedTime())
				updatedProcesses = append(updatedProcesses, processes[i]+"|"+newTimeDuration.String())
				oldContentSlice[j] = "" // clear the old process so it won't be added back later
				found = true
				break
			}
		}
		// if the process is not in the file, add the process and the time
		if !found {
			updatedProcesses = append(updatedProcesses, processes[i]+"|"+bt.ElapsedTime().String())
		}
	}

	// Append any old processes that were not updated
	for _, oldProcess := range oldContentSlice {
		if oldProcess != "" {
			updatedProcesses = append(updatedProcesses, oldProcess)
		}
	}

	newContent = strings.Join(updatedProcesses, "\n")

	// write the new content to the file
	err := fl.WriteFilePath(fl.Folder+"/"+fl.ExeLog, newContent)
	if err != nil {
		f.Printf("Failed to write to file: %v\n", err)
	}

}

func getTimeFromProcess(process string) string {
	split := strings.Split(process, "|")
	if len(split) > 1 {
		return split[1]
	}
	return t.Now().String()
}
