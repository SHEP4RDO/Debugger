package mklog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DetailedError represents an error with additional information about the stack trace.
type DetailedError struct {
	Err              error     // Original error.
	StackInfo        string    // Stack trace information.
	Time             time.Time // Time when the error occurred.
	FunctionName     string    // Name of the function that caused the error.
	CallingArguments string    // Arguments passed to the function.
	File             string    // File where the error occurred.
	Line             int       // Line number of the code where the error occurred.
}

// NewDetailedError creates a new DetailedError, capturing the original error and contextual information.
func NewDetailedError(err error, args ...interface{}) DetailedError {
	stackInfo := getStackInfo()                                            // Gather stack trace information.
	functionName, file, line, callingArguments := getFunctionInfo(args...) // Get function info and arguments.
	return DetailedError{
		Err:              err,
		StackInfo:        stackInfo,
		Time:             time.Now(),
		FunctionName:     functionName,
		CallingArguments: callingArguments,
		File:             file,
		Line:             line,
	}
}

// ErrorStack returns a string representation of the DetailedError, including the original error and stack trace.
func (de DetailedError) ErrorStack() string {
	return fmt.Sprintf("\nTime: %s\nFile: %s:%d\nFunction: %s\nArguments: %s\n%s",
		de.Time.Format("2006-01-02 15:04:05"),
		de.File,
		de.Line,
		de.FunctionName,
		de.CallingArguments,
		de.StackInfo)
}

// Error returns the string representation of the original error.
func (de DetailedError) Error() string {
	return de.Err.Error()
}

// getStackInfo collects stack trace information for the current goroutine.
func getStackInfo() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // Skip 3 levels to start from the actual call.
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	sb.WriteString("Stack Trace:\n")

	for {
		frame, more := frames.Next()
		// Exclude Go system functions and mklog functions.
		if !strings.Contains(frame.Function, "runtime.") && !strings.Contains(frame.Function, "mklog.") {
			sb.WriteString(fmt.Sprintf("  %s:%d %s\n", filepath.Base(frame.File), frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}

	return sb.String()
}

// getFunctionInfo returns information about the function, arguments, file, and line of the call.
func getFunctionInfo(args ...interface{}) (string, string, int, string) {
	pc, file, line, ok := runtime.Caller(2) // Get information about the calling function.
	if !ok {
		return "unknown", "unknown", 0, ""
	}
	function := runtime.FuncForPC(pc)
	if function == nil {
		return "unknown", filepath.Base(file), line, ""
	}
	return function.Name(), filepath.Base(file), line, fmt.Sprintf("%v", args) // Format the arguments into a string.
}
