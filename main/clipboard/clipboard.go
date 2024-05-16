package clipboard

import (
	fl "bytebrother/main/filer"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.design/x/clipboard"
	"image"
	"image/png"
	"os"
)

func Image() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtImage)
	for data := range ch {
		// image to folder
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := os.Create(fl.ClipboardFolder + fl.Today() + "/clipboard_" + fl.CurrentTime() + ".png")

		if err != nil {
			fmt.Println(err)
			return
		}

		err = png.Encode(f, img)
		if err != nil {
			fmt.Println(err)
			return
		}

		f.Close()
	}

}

// Text watches the clipboard for text changes
// and appends it to single json file

type ClipboardData struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

func Text() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		// text to folder
		clipboardData := readClipboard()
		clipboardData = append(clipboardData, ClipboardData{
			Time: fl.CurrentTime(),
			Text: string(data),
		})

		newContent, err := json.MarshalIndent(clipboardData, "", "    ")
		if err != nil || len(clipboardData) == 0 {
			fmt.Printf("Failed to marshal log entries: %v\n", err)
		}

		err = os.WriteFile(fl.ClipboardFolder+fl.Today()+"/"+fl.ClipboardFile, newContent, 0644)
		if err != nil {
			fmt.Printf("Failed to write to file: %v\n", err)
		}
	}
}

func readClipboard() []ClipboardData {
	var clipboardData []ClipboardData
	if fl.IfFileExists(fl.ClipboardFolder + fl.Today() + "/" + fl.ClipboardFile) {
		oldContent, err := fl.ReadFilePath(fl.ClipboardFolder + fl.Today() + "/clipboard.json")
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
		} else {

			err := json.Unmarshal(oldContent, &clipboardData)
			if err != nil {
				fmt.Printf("Failed to unmarshal existing log: %v\n", err)
			}
		}
	} else {
		fmt.Println("File doesn't exist")
	}
	return clipboardData
}
