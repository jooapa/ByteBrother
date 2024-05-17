package filer

import (
	f "fmt"
	"os"
	"path/filepath"
	"time"
)

var Folder string = "bytebrother/"
var LogFolder string = Folder + "logs/"
var ProcessFolder string = LogFolder + "exes/"
var ExeLog string = "exes.json"

var SettingsFile string = "settings.json"
var NetworkLogs string = "network_logs.log"
var KeyLogs string = "key_logs.log"
var ClipboardFolder string = Folder + "clipboard/"
var ClipboardFile string = "clipboard.json"

func WalkDir(root string, walkFn func(path string, info os.FileInfo, err error) error) error {
	return filepath.Walk(root, walkFn)
}

func MakeDir(path string) error {
	return os.Mkdir(path, 0755)
}

func MakeFile(path string) (*os.File, error) {
	return os.Create(path)
}

func IfFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func IfDirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Pwd() (string, error) {
	return os.Getwd()
}

func WriteFile(file *os.File, data []byte) (int, error) {
	return file.Write(data)
}

func WriteFilePath(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func AppendToFile(filename, content string) error {
	// Open the file in append mode. Create the file if it doesn't exist.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func CloseFile(file *os.File) error {
	return file.Close()
}

func ReadFile(file *os.File) ([]byte, error) {
	return os.ReadFile(file.Name())
}

func ReadFilePath(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func MakeNecessaryFiles() {
	// Create a new directory
	if !IfDirExists(Folder) {
		err := MakeDir(Folder)
		if err != nil {
			f.Printf("Failed to create directory: %v\n", err)
		}
	}

	// Create a new directory for logs
	if !IfDirExists(LogFolder) {
		err := MakeDir(LogFolder)
		if err != nil {
			f.Printf("Failed to create directory: %v\n", err)
		}
	}

	// Create a new directory for process logs
	if !IfDirExists(ProcessFolder) {
		err := MakeDir(ProcessFolder)
		if err != nil {
			f.Printf("Failed to create directory: %v\n", err)
		}
	}

	// Create a new file for settings
	if !IfFileExists(Folder + SettingsFile) {
		_, err := MakeFile(Folder + SettingsFile)
		if err != nil {
			f.Printf("Failed to create file: %v\n", err)
		}
	}

	if !IfDirExists(ClipboardFolder) {
		err := MakeDir(ClipboardFolder)
		if err != nil {
			f.Printf("Failed to create directory: %v\n", err)
		}
	}

	MakeTodayClipboardFolder()
}

func MakeTodayClipboardFolder() {
	if !IfDirExists(TodaysClipboardFolder()) {
		err := MakeDir(TodaysClipboardFolder())
		if err != nil {
			f.Printf("Failed to create directory: %v\n", err)
		}
	}
}

func TodaysClipboardFolder() string {
	return ClipboardFolder + Today() + "/"
}

func Today() string {
	return time.Now().Format("02-01-2006")
}

func CurrentTime(op string) string {
	return time.Now().Format("15" + op + "04" + op + "05")
}
