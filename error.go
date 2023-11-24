package mklog

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// DetailedError represents an error with additional stack trace information.
type DetailedError struct {
	Err       error  // Original error.
	StackInfo string // Stack trace information.
}

// NewDetailedError creates a new DetailedError with the given error and captures the stack trace.
func NewDetailedError(err error) DetailedError {
	stackInfo := getStackInfo()
	return DetailedError{Err: err, StackInfo: stackInfo}
}

// Error returns the string representation of the DetailedError, including the original error and stack trace.
func (de DetailedError) Error() string {
	return fmt.Sprintf("%v %v", de.Err.Error(), de.StackInfo)
}

// getStackInfo retrieves the stack trace information for the current goroutine.
func getStackInfo() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	stackInfo := "Stack Trace:\n"
	for {
		frame, more := frames.Next()
		stackInfo += fmt.Sprintf("  %s:%d %s\n", filepath.Base(frame.File), frame.Line, frame.Function)
		if !more {
			break
		}
	}

	return stackInfo
}
