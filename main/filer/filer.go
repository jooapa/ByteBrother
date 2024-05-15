package filer

import (
	f "fmt"
	"os"
	"path/filepath"
)

var Folder string = "ByteBrother"
var ExeLog string = "exes.json"

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
	} else {
		f.Printf("Directory already exists\n")
	}
}