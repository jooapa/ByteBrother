package process

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	FILE_NOTIFY_CHANGE_FILE_NAME = 0x1
	FILE_LIST_DIRECTORY          = 0x1
)

func MonitorDirectory() {
	dirPath := "C:\\Windows\\System32"
	buffer := make([]byte, 1024)

	hDir, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(dirPath),
		FILE_LIST_DIRECTORY,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}
	defer syscall.CloseHandle(hDir)

	var overlapped syscall.Overlapped

	for {
		var bytesReturned uint32
		err := syscall.ReadDirectoryChanges(
			syscall.Handle(hDir),
			&buffer[0],
			uint32(len(buffer)),
			true,
			FILE_NOTIFY_CHANGE_FILE_NAME,
			&bytesReturned,
			&overlapped,
			0,
		)
		if err != nil {
			fmt.Println("Error reading directory changes:", err)
			return
		}

		var offset uint32
		for offset < bytesReturned {
			info := (*syscall.FileNotifyInformation)(unsafe.Pointer(&buffer[offset]))

			filename := syscall.UTF16ToString((*[syscall.MAX_PATH]uint16)(unsafe.Pointer(&info.FileName))[:])
			if isExe(filename) {
				fmt.Println("Executable file opened:", filename)
			}

			offset += info.NextEntryOffset
		}
	}
}

func isExe(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".exe"
}
