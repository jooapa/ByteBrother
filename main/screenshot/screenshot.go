package screenshot

import (
	"fmt"
	"image/png"
	"log"
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

		if os.IsNotExist(os.MkdirAll(fl.ScreenshotFolder+fl.Today(), 0755)) {
			fl.MakeDir(fl.ScreenshotFolder + fl.Today())
		}

		fileName := fmt.Sprintf(fl.ScreenshotFolder+fl.Today()+"/"+fl.CurrentTime("-")+"_%d.png", i)
		file, _ := os.Create(fileName)
		defer file.Close()

		// Use a lower compression level to reduce the file size
		encoder := png.Encoder{CompressionLevel: png.BestCompression}

		if err := encoder.Encode(file, img); err != nil {
			log.Fatalf("Failed to encode image: %v", err)
		}
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
