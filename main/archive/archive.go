package archive

import (
	fl "bytebrother/main/filer"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var ArchiveRowCount int = 1000

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ShouldArchive(path string) bool {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return false
	}
	defer file.Close()

	// Count the number of lines in the file
	count, err := LineCounter(file)
	if err != nil {
		fmt.Printf("Failed to count lines: %v\n", err)
		return false
	}

	// If the number of lines is greater than 100, return true
	if count > ArchiveRowCount {
		fmt.Printf("Should archive: %v\n", count)
		return true
	}

	return false
}

func Archive(path string) error {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Walk through the file directory and see if there are other archives
	// If there are, increment the number of the archive

	// Create a new file with the archive number
	archiveNum := 1
	archivePath := path + ".archive" + strconv.Itoa(archiveNum)

	for fl.IfFileExists(archivePath) {
		archiveNum++
		archivePath = path + ".archive" + strconv.Itoa(archiveNum)
	}

	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// Copy the contents of the original file to the archive file
	_, err = io.Copy(archiveFile, file)
	if err != nil {
		return err
	}

	// delete the original file
	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
