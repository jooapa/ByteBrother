package global

import t "time"

var (
	SaveProcessInformationInFile t.Duration = 5000
	ProcessSaveInterval          t.Duration = 5000
	ProcessInterval              t.Duration = 1000
	ScreenshotInterval           t.Duration = 60
	ScreenshotInterval_sec       t.Duration = 60
	ChosenIndex                             = 69420
	NetworkIndexToMonitor                   = 69420
	ArchiveRowCount                         = 6000
	NumRowsInArchive                        = 6000
	ArchiveScreenshotsAfter                 = true
	CanTakeScreenshot                       = true
	ClipboardTextTrackingEnabled            = true
	ClipboardImageTrackingEnable            = true
	KeyloggerEnabled                        = true
	ClipboardTracking_text                  = false
	ClipboardTracking_image                 = false
	Keylogging                              = false
	ProcessSaving                           = false
	IsSetup                                 = false
)
