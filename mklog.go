package mklog

import (
	"os"
	"os/signal"
)

// LogLevel represents the severity levels of log messages.
type LogLevel int

const (
	DebugLevel   LogLevel = iota // DebugLevel represents debug log messages.
	InfoLevel                    // InfoLevel represents informational log messages.
	WarningLevel                 // WarningLevel represents warning log messages.
	ErrorLevel                   // ErrorLevel represents error log messages.
	FatalLevel                   // FatalLevel represents fatal log messages.
	TraceLevel                   // TraceLevel represents trace log messages.
)

// logLevelNames maps LogLevel values to their corresponding string representations.
var logLevelNames = map[LogLevel]string{
	DebugLevel:   "DEBUG",
	InfoLevel:    "INFO",
	WarningLevel: "WARNING",
	ErrorLevel:   "ERROR",
	FatalLevel:   "FATAL",
	TraceLevel:   "TRACE",
}

// Debugger is a logging utility that provides various configuration options.
type Debugger struct {
	moduleName          string              // Name of the module using the debugger.
	submodules          []string            // Submodules associated with the module.
	customLogLevelNames map[LogLevel]string // Custom log level names provided by the user.
	debugMode           bool                // Indicates whether debug mode is enabled.
	dateFormat          string              // Format for log timestamps.

	isConsoleOutput     bool           // Indicates whether console output is enabled.
	detailedErrorOutput bool           // Indicates whether detailed error output is enabled.
	logFinishChannel    chan struct{}  // Channel to signal log finishing.
	signalChannel       chan os.Signal // Channel to handle OS signals.

	logLevel     LogLevel     // Minimum log level to be recorded.
	logFormatter LogFormatter // Formatter for log messages.
	log          log          // Internal log instance.
}

// NewDebugLogger creates a new instance of Debugger with the specified module name and optional submodules.
// Default parameters:
//   - logLevel: InfoLevel
//   - logFormatter: Default log formatter set to PlainTextFormatter with the date format "02.01.2006"
//   - Signal notification: Debugger is notified about OS interrupt signals
//   - log.logToFile: Log to file is initially disabled
//   - log.logDateFileFormat: Default log file date format set to "02.01.2006"
func NewDebugLogger(moduleName string, submodules ...string) *Debugger {
	p := &Debugger{
		moduleName:       moduleName,
		submodules:       submodules,
		logLevel:         InfoLevel,                        // Default log level is set to Info.
		log:              log{},                            // Initialize internal log instance.
		signalChannel:    make(chan os.Signal, 1),          // Create a channel for OS signals with buffer size 1.
		logFinishChannel: make(chan struct{}),              // Create a channel to signal log finishing.
		logFormatter:     PlainTextFormatter{"02.01.2006"}, // Default log formatter is set to Plain Text with a specific date format.
	}
	signal.Notify(p.signalChannel, os.Interrupt) // Notify the debugger about OS interrupt signals.
	p.log.isToFile = false                       // Log to file is initially disabled.
	p.log.DateFileFormat = "02.01.2006"          // Default log file date format.
	p.isConsoleOutput = true
	p.log.FileName = "log_file"
	p.log.FilePath = "logs"
	p.log.FileType = ".log"
	return p
}

// OPTIONS

// SetDebugMode enables or disables debug mode for the Debugger instance.
func (d *Debugger) SetDebugMode(mode bool) *Debugger {
	d.debugMode = mode
	return d
}

// SetDateFormat sets the date format for log messages.
func (d *Debugger) SetDateFormat(format string) *Debugger {
	d.dateFormat = format
	d.logFormatter.SetLogDateFormat(format)
	return d
}

// SetLogDate enables or disables the inclusion of the log date in file names.
func (d *Debugger) SetLogDate(isLogDate bool) *Debugger {
	d.log.isDateFile = isLogDate
	return d
}

// SetLogDateFormat sets the date format for log file names.
func (d *Debugger) SetLogDateFormat(format string) *Debugger {
	d.log.DateFileFormat = format
	return d
}

// SetConsoleOutput enables or disables console output for log messages.
func (d *Debugger) SetConsoleOutput(mode bool) *Debugger {
	d.isConsoleOutput = mode
	return d
}

// SetDetailedErrorOutput enables or disables detailed error output.
func (d *Debugger) SetDetailedErrorOutput(enabled bool) *Debugger {
	d.detailedErrorOutput = enabled
	return d
}

// GetSignalChannel returns the signal channel used for interrupt signals.
func (d *Debugger) GetSignalChannel() chan os.Signal {
	return d.signalChannel
}

// LOG FILES

// SetLogFile enables or disables logging to a file and sets the file path and name.
func (d *Debugger) SetIsLogFile(isLog bool) *Debugger {
	d.log.isToFile = isLog
	return d
}
func (d *Debugger) SetLogFileName(filename string) *Debugger {
	d.log.FileName = filename
	return d
}
func (d *Debugger) SetLogPath(filename string) *Debugger {
	d.log.FilePath = filename
	return d
}

// SetLogFileType sets the log file type.
func (d *Debugger) SetLogFileType(_type string) *Debugger {
	d.log.FileType = _type
	return d
}

// CreateLogFile creates the log file based on the configured settings.
func (d *Debugger) CreateLogFile() error {
	err := d.createLogFile()
	return err
}

func (d *Debugger) SetLogDateFileNameFormat(_type string) *Debugger {
	return d
}

func (d *Debugger) SetDefaultLogPath() *Debugger {
	d.log.FilePath = "logs"
	return d
}

//FORMATTERS

// SetLogFormatter sets the log formatter for the Debugger instance.
func (d *Debugger) SetLogFormatter(formatter LogFormatter) *Debugger {
	d.logFormatter = formatter
	return d
}

// SetUserDefinedFormatter sets a user-defined log formatter function and date format.
func (d *Debugger) SetUserDefinedFormatter(formatFunc UserDefinedFormatterFunc, dateFormat string) *Debugger {
	d.logFormatter = UserDefinedFormatter{formatFunc, dateFormat}
	return d
}

// SetCustomLogLevelNames sets custom log level names provided by the user.
func (d *Debugger) SetCustomLogLevelNames(customLogLevelNames map[LogLevel]string) *Debugger {
	d.customLogLevelNames = customLogLevelNames
	return d
}
