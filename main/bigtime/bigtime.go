package bigtime

import (
	"time"
)

var startTime time.Time // Variable to store the start time of the program

func init() {
	startTime = time.Now() // Initialize the start time when the package is initialized
}

// reset
func Reset() {
	startTime = time.Now()
}

// ElapsedTime returns the duration since the program started running
func ElapsedTime() time.Duration {
	return time.Since(startTime)
}

func AddTimeTogether(d1, d2 time.Duration) time.Duration {
	return d1 + d2
}

func StringToDuration(str string) (time.Duration, error) {
	return time.ParseDuration(str)
}
