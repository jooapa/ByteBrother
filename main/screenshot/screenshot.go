package screenshot

import (
	"fmt"
	"image/png"
	"os"
	// "path/filepath"

	fl "bytebrother/main/filer"
	// g "bytebrother/main/global"
	ss "github.com/kbinani/screenshot"
)

func TakeScreenshot() {
	n := ss.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := ss.GetDisplayBounds(i)

		img, err := ss.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}

		fileName := fmt.Sprintf(fl.ScreenshotFolder+fl.Today()+"/"+fl.CurrentTime("-")+"_%d.png", i)
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)
	}
}

func CanArchiveOlderFolder() (bool, string) {
	today := fl.Today()

	files, err := os.ReadDir(fl.ScreenshotFolder)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	for _, file := range files {
		if file.Name() != today {
			if len(file.Name()) == 10 {
				return true, file.Name()
			}
		}
	}

	return false, ""
}

// func CanArchiveCurrentFiles() bool {
// 	today := fl.Today()

// 	dirSize, err := fl.GetFilesInDir(fl.ScreenshotFolder + today)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return false
// }
