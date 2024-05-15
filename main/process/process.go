package process

import (
	"fmt"
	"github.com/mitchellh/go-ps"
)

func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func GetProcesses() []string {
	var process_names []string

	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("Failed to retrieve processes:", err)
		return nil
	}

	for _, process := range processes {
		if Contains(process_names, process.Executable()) || process.Executable() == "[System Process]" {
			continue
		}

		process_names = append(process_names, process.Executable())
	}

	return process_names
}
