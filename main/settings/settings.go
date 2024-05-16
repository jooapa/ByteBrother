package settings

import (
	"encoding/json"
	"fmt"
	"os"

	fl "bytebrother/main/filer"
)

// Settings represents the settings of the application.
type Settings struct {
	ProcessInterval       int `json:"process_interval_milliseconds"`
	NetworkIndexToMonitor int `json:"network_index_to_monitor"`
	NumRowsInArchive      int `json:"num_rows_in_archive"`
}

// LoadSettings loads the settings from the settings file.
func LoadSettings() Settings {
	settings := Settings{
		ProcessInterval:       1000,
		NetworkIndexToMonitor: 69420,
		NumRowsInArchive:      5000,
	}

	if fl.IfFileExists(fl.Folder + fl.SettingsFile) {
		file, err := os.ReadFile(fl.Folder + fl.SettingsFile)
		if err != nil {
			fmt.Println("Error reading settings file:", err)
			return settings
		}

		if len(file) == 0 {
			fmt.Println("Settings file is empty. Making a new one with default settings.")
			SaveSettings(settings)
			return settings
		}

		err = json.Unmarshal(file, &settings)
		if err != nil {
			fmt.Println("Error unmarshalling settings file:", err)
			return settings
		}
	} else {
		fmt.Println("Settings file does not exist. Making a new one with default settings.")
		SaveSettings(settings)
	}

	return settings
}

// SaveSettings saves the settings to the settings file.
func SaveSettings(settings Settings) {
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling settings:", err)
		return
	}

	err = os.WriteFile(fl.Folder+fl.SettingsFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing settings file:", err)
	}
}
