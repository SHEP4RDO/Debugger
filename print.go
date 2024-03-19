package mklog

import (
	"fmt"
	"time"
)

// CustomDebug outputs a custom log message if debug mode.
// It includes the specified message and error and respects the custom log level names.
func (d *Debugger) CustomDebug(logLevel LogLevel, msg string, args ...interface{}) {
	if d.debugMode {
		d.printLogCustom(logLevel, msg, args...)
		d.printCustomLogToFile(logLevel, msg, args...)
	}
}

// Custom outputs a custom log message at the specified log level.
// It includes the specified message and error and respects the custom log level names.
func (d *Debugger) Custom(logLevel LogLevel, msg string, args ...interface{}) {
	d.printLogCustom(logLevel, msg, args...)
	d.printCustomLogToFile(logLevel, msg, args...)
}

// printLogCustom formats and prints the custom log message to the console.
// It uses the custom log level names if provided, otherwise falls back to the default log level names.
func (d *Debugger) printLogCustom(logLevel LogLevel, msg string, args ...interface{}) {
	var logLevelName string
	if d.customLogLevelNames != nil {
		logLevelName = d.customLogLevelNames[logLevel]
	} else {
		logLevelName = logLevelNames[logLevel]
	}
	logMessage := d.formatLog(msg, args...)
	fmt.Print(d.logFormatter.Format(
		logMessage,
		logLevelName,
		d.moduleName,
		d.submodules,
		time.Now().Format(d.dateFormat),
	))
}

// printCustomLogToFile writes the custom log message to the log file if logging to a file is enabled.
// It uses the custom log level names if provided, otherwise falls back to the default log level names.
func (d *Debugger) printCustomLogToFile(logLevel LogLevel, msg string, args ...interface{}) error {
	if d.log.isToFile {
		var logLevelName string
		if d.customLogLevelNames != nil {
			logLevelName = d.customLogLevelNames[logLevel]
		} else {
			logLevelName = logLevelNames[logLevel]
		}
		logMessage := d.formatLog(msg, args...)
		toPrint := d.logFormatter.Format(logMessage, logLevelName, d.moduleName, d.submodules, time.Now().Format(d.dateFormat))
		return d.writeLog(toPrint)
	}
	return nil
}

// Debug outputs a log message at the Debug level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
func (d *Debugger) Debug(msg string, args ...interface{}) {
	if d.debugMode {
		d.printLog(DebugLevel, msg, args...)
		d.printLogToFile(msg, args...)
	}
}

// Trace outputs a log message at the Trace level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
// Trace logs are only printed if the debug mode is enabled, and the log level is set to Trace.
func (d *Debugger) Trace(msg string, args ...interface{}) {
	if d.debugMode && d.logLevel == TraceLevel {
		d.printLog(TraceLevel, msg, args...)
		d.printLogToFile(msg, args...)
	}
}

// Info outputs a log message at the Info level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
func (d *Debugger) Info(msg string, args ...interface{}) {
	d.printLog(InfoLevel, msg, args...)
	d.printLogToFile(msg, args...)
}

// Warning outputs a log message at the Warning level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
func (d *Debugger) Warning(msg string, args ...interface{}) {
	d.printLog(WarningLevel, msg, args...)
	d.printLogToFile(msg, args...)
}

// Error outputs a log message at the Error level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
func (d *Debugger) Error(msg string, args ...interface{}) {
	d.printLog(ErrorLevel, msg, args...)
	d.printLogToFile(msg, args...)
}

// Fatal outputs a log message at the Fatal level with the specified message and error (if any).
// It prints the log message to the console and, if enabled, to the log file.
func (d *Debugger) Fatal(msg string, args ...interface{}) {
	d.printLog(FatalLevel, msg, args...)
	d.printLogToFile(msg, args...)
}

// printLog prints the log message to the console with the specified log level.
// It updates the log level of the debugger and formats the log message using the configured log formatter.
func (d *Debugger) printLog(logLevel LogLevel, msg string, args ...interface{}) {
	d.logLevel = logLevel
	logMessage := fmt.Sprintf(msg, args...)
	fmt.Print(d.logFormatter.Format(logMessage, logLevelNames[logLevel], d.moduleName, d.submodules, time.Now().Format(d.dateFormat)))
}

// printLogToFile writes the log message to the log file if logging to a file is enabled.
// It uses the log formatter to format the log message before writing it to the file.
func (d *Debugger) printLogToFile(msg string, args ...interface{}) error {
	if d.log.isToFile {
		logMessage := fmt.Sprintf(msg, args...)
		toPrint := d.logFormatter.Format(logMessage, logLevelNames[d.logLevel], d.moduleName, d.submodules, time.Now().Format(d.dateFormat))
		return d.writeLog(toPrint)
	}
	return nil
}

// formatLog creates a formatted log message combining the specified message and error (if any).
func (d *Debugger) formatLog(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// DebugDetailed outputs a detailed log message at the Debug level.
// It includes the specified message and error with additional stack trace information.
// DebugDetailed logs are printed to the console and, if enabled, to the log file.
func (d *Debugger) DebugDetailed(msg string, err error) {
	if d.debugMode {
		d.printLog(DebugLevel, msg, NewDetailedError(err))
		d.printLogToFile(msg, NewDetailedError(err))
	}
}

// TraceDetailed outputs a detailed log message at the Trace level.
// It includes the specified message and error with additional stack trace information.
// TraceDetailed logs are only printed if the debug mode is enabled, and the log level is set to Trace.
func (d *Debugger) TraceDetailed(msg string, err error) {
	if d.debugMode && d.logLevel == TraceLevel {
		d.printLog(TraceLevel, msg, NewDetailedError(err))
		d.printLogToFile(msg, NewDetailedError(err))
	}
}

// InfoDetailed outputs a detailed log message at the Info level.
// It includes the specified message and error with additional stack trace information.
// InfoDetailed logs are printed to the console and, if enabled, to the log file.
func (d *Debugger) InfoDetailed(msg string, err error) {
	d.printLog(InfoLevel, msg, NewDetailedError(err))
	d.printLogToFile(msg, NewDetailedError(err))
}

// WarningDetailed outputs a detailed log message at the Warning level.
// It includes the specified message and error with additional stack trace information.
// WarningDetailed logs are printed to the console and, if enabled, to the log file.
func (d *Debugger) WarningDetailed(msg string, err error) {
	d.printLog(WarningLevel, msg, NewDetailedError(err))
	d.printLogToFile(msg, NewDetailedError(err))
}

// ErrorDetailed outputs a detailed log message at the Error level.
// It includes the specified message and error with additional stack trace information.
// ErrorDetailed logs are printed to the console and, if enabled, to the log file.
func (d *Debugger) ErrorDetailed(msg string, err error) {
	d.printLog(ErrorLevel, msg, NewDetailedError(err))
	d.printLogToFile(msg, NewDetailedError(err))
}

// FatalDetailed outputs a detailed log message at the Fatal level.
// It includes the specified message and error with additional stack trace information.
// FatalDetailed logs are printed to the console and, if enabled, to the log file.
func (d *Debugger) FatalDetailed(msg string, err error) {
	d.printLog(FatalLevel, msg, NewDetailedError(err))
	d.printLogToFile(msg, NewDetailedError(err))
}
