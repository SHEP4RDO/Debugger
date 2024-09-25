package mklog

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

// LogLevel represents the severity levels of log messages.
type LogLevel int

const (
	TraceLevel   LogLevel = iota // TraceLevel represents trace log messages.
	DebugLevel                   // DebugLevel represents debug log messages.
	InfoLevel                    // InfoLevel represents informational log messages.
	WarningLevel                 // WarningLevel represents warning log messages.
	ErrorLevel                   // ErrorLevel represents error log messages.
	FatalLevel                   // FatalLevel represents fatal log messages.
)

// UnmarshalYAML parses the log level from a YAML configuration file.
func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var levelStr string
	if err := unmarshal(&levelStr); err != nil {
		return err
	}

	switch strings.ToUpper(levelStr) {
	case "INFO":
		*l = InfoLevel
	case "DEBUG":
		*l = DebugLevel
	case "WARN", "WARNING":
		*l = WarningLevel
	case "ERROR":
		*l = ErrorLevel
	case "FATAL":
		*l = FatalLevel
	default:
		return fmt.Errorf("invalid log level: %s", levelStr)
	}
	return nil
}

var (
	// Default file paths for logging
	MKLOG_DirDefault      = "logs"     // Default directory for log files
	MKLOG_FileNameDefault = "log_file" // Default log file name
	MKLOG_FileTypeDefault = ".log"     // Default log file type

	// Default periods for file management
	MKLOG_FileFolderPeriodDefault = time.Hour * 24 // Default period for log folder management

	// Default time formats for logging
	MKLOG_TimeFolderFormatDefault = "2006-01-02"          // Default date format for folder names
	MKLOG_TimeFileFormatDefault   = "2006-01-02"          // Default date format for log files
	MKLOG_TimeLogFormatDefault    = "2006-01-02 15:04:05" // Default timestamp format for log entries

	// Default log formatter configuration
	MKLOG_FormatterDefault = PlainTextFormatter{"02.01.2006"}

	// Default buffer size for asynchronous logging
	MKLOG_BufferSizeDefault = 100 // Default size of the log buffer
)

// AsyncLog configures asynchronous logging settings.
type AsyncLog struct {
	Enable     bool `json:"Enable" yaml:"Enable"`           // Enable asynchronous logging
	BufferSize int  `json:"buffer_size" yaml:"buffer_size"` // Size of the log buffer
}

// LogRule defines the rules for logging levels and outputs.
type LogRule struct {
	MinLevel            LogLevel            `json:"min_level" yaml:"min_level"`                           // Minimum log level
	MaxLevel            LogLevel            `json:"max_level" yaml:"max_level"`                           // Maximum log level
	CurrentLevel        LogLevel            `json:"current_level" yaml:"current_level"`                   // Current log level
	FileName            string              `json:"file_name" yaml:"file_name"`                           // Name of the log file
	FileType            string              `json:"file_type" yaml:"file_type"`                           // Type of the log file
	IsDateFile          bool                `json:"is_date_file" yaml:"is_date_file"`                     // Flag for date-based file naming
	DateFileFormat      string              `json:"date_file_format" yaml:"date_file_format"`             // Format for date in file names
	LogFormatter        LogFormatter        `json:"log_formatter" yaml:"log_formatter"`                   // Formatter for log entries
	ModuleName          string              `json:"module_name" yaml:"module_name"`                       // Name of the module being logged
	Submodules          []string            `json:"submodules" yaml:"submodules"`                         // List of submodules for logging
	IsConsoleOutput     bool                `json:"is_console_output" yaml:"is_console_output"`           // Flag for console output of logs
	DebugMode           bool                `json:"debug_mode" yaml:"debug_mode"`                         // Flag for enabling debug mode
	DebugModeStatus     LogLevel            `json:"debug_mode_status" yaml:"debug_mode_status"`           // Current status of debug mode
	DateFormat          string              `json:"date_format" yaml:"date_format"`                       // Date format for log entries
	DetailedErrorOutput bool                `json:"detailed_error_output" yaml:"detailed_error_output"`   // Flag for detailed error output
	CustomLogLevelNames map[LogLevel]string `json:"custom_log_level_names" yaml:"custom_log_level_names"` // Custom names for log levels

	FileLog    FileLog    `json:"file_log" yaml:"file_log"`       // Configuration for file logging
	FileFolder FileFolder `json:"file_folder" yaml:"file_folder"` // Configuration for folder logging
	AsyncLog   AsyncLog   `json:"async_log" yaml:"async_log"`     // Configuration for asynchronous logging

	logFinishChannel chan struct{}  `json:"-" yaml:"-"` // Channel to signal completion of logging
	signalChannel    chan os.Signal `json:"-" yaml:"-"` // Channel for OS signal handling
	logChannel       chan string    `json:"-" yaml:"-"` // Channel for log message transmission
}

// Debugger is a logging utility that provides various configuration options for logging.
type Debugger struct {
	LogRules map[string][]*LogRule `yaml:"log_rules"` // Map of logging rules categorized by module names
}

// NewDebugLogger initializes a new Debugger instance with default logging rules for a module.
func NewDebugLogger(moduleName string, submodules ...string) *Debugger {
	initRule := &LogRule{
		MinLevel:            DebugLevel,                                 // Set minimum log level to Debug
		MaxLevel:            FatalLevel,                                 // Set maximum log level to Fatal
		CurrentLevel:        InfoLevel,                                  // Set current log level to Info
		LogFormatter:        PlainTextFormatter{"02.01.2006"},           // Set log formatting
		ModuleName:          moduleName,                                 // Set the module name
		IsConsoleOutput:     true,                                       // Enable console output
		DebugMode:           true,                                       // Enable debug mode
		DebugModeStatus:     TraceLevel,                                 // Set debug mode status
		DateFormat:          "02.01.2006",                               // Set date format for logs
		DetailedErrorOutput: false,                                      // Disable detailed error output by default
		logFinishChannel:    make(chan struct{}),                        // Channel for signaling log completion
		signalChannel:       make(chan os.Signal, 1),                    // Channel for handling OS signals
		logChannel:          make(chan string, MKLOG_BufferSizeDefault), // Channel for log message transmission
		FileLog: FileLog{
			Enable:     false, // Disable file logging by default
			IsDateFile: false, // Disable date-based file naming by default
		},
	}

	p := &Debugger{
		LogRules: make(map[string][]*LogRule), // Initialize log rules map
	}

	p.LogRules[moduleName] = append(p.LogRules[moduleName], initRule) // Add initial logging rule for the module

	signal.Notify(p.LogRules[moduleName][0].signalChannel, os.Interrupt) // Setup OS interrupt notification
	return p                                                             // Return the initialized Debugger instance
}

// AddRule adds a new logging rule to the Debugger instance for a specified module.
// If the module does not exist, it initializes a new slice for log rules.
func (d *Debugger) AddRule(moduleName string, rule LogRule) *Debugger {
	if _, exists := d.LogRules[moduleName]; !exists {
		d.LogRules[moduleName] = []*LogRule{}
	}
	d.LogRules[moduleName] = append(d.LogRules[moduleName], &rule)
	return d
}

// NewLogRule creates a new logging rule with default configuration for a given module name.
// It accepts optional configuration functions to customize the log rule.
func (d *Debugger) NewLogRule(moduleName string, opts ...Option) *Debugger {
	// Create a base configuration with default values.
	lr := &LogRule{
		MinLevel:        InfoLevel,             // Minimum log level for this rule.
		MaxLevel:        ErrorLevel,            // Maximum log level for this rule.
		CurrentLevel:    InfoLevel,             // Current log level for logging.
		ModuleName:      moduleName,            // Name of the module associated with this rule.
		IsConsoleOutput: false,                 // Disable console output by default.
		DateFormat:      "02-01-2006 15:04:05", // Default date format for logs.
		FileLog: FileLog{
			Enable:     false,      // Disable file logging by default.
			FilePath:   "logs",     // Default directory for log files.
			FileName:   "log_file", // Default log file name.
			FileType:   ".log",     // Default file extension for log files.
			IsDateFile: false,      // Disable date in file name by default.
		},
		signalChannel:    make(chan os.Signal, 1), // Channel to handle OS signals.
		logFinishChannel: make(chan struct{}),     // Channel to signal the end of logging.
	}

	// Apply any provided options to customize the log rule.
	for _, opt := range opts {
		opt(lr)
	}

	// Set default log formatter if not specified.
	if lr.LogFormatter == nil {
		lr.LogFormatter = &PlainTextFormatter{"02.01.2006"}
		fmt.Println("[Warning] LogFormatter not set, using PlainTextFormatter.")
	}

	// Add the new log rule to the array of rules for the module.
	d.LogRules[moduleName] = append(d.LogRules[moduleName], lr)
	signal.Notify(lr.signalChannel, os.Interrupt) // Notify on OS interrupts.

	// Create the log file if file logging is enabled.
	if lr.FileLog.Enable {
		if err := lr.createLogFile(); err != nil {
			fmt.Println("[mklog] error while creating log file ", lr.ModuleName, ": ", err)
		}
	}

	// Start asynchronous logging if enabled.
	if lr.AsyncLog.Enable {
		lr.logChannel = make(chan string, lr.AsyncLog.BufferSize)
		go lr.StartAsyncLogging()
	}

	return d
}

// CloseAsyncLogging closes all log channels for asynchronous logging in the Debugger instance.
func (d *Debugger) CloseAsyncLogging() {
	for _, rules := range d.LogRules {
		for _, v := range rules {
			if v.AsyncLog.Enable {
				close(v.logChannel) // Close the log channel to stop logging.
			}
		}
	}
}

// SetDebugMode enables or disables debug mode for the log rule.
func (d *LogRule) SetDebugMode(mode bool) *LogRule {
	d.DebugMode = mode
	return d
}

// SetDebugLevel sets the debug level for the log rule.
func (d *LogRule) SetDebugLevel(level LogLevel) *LogRule {
	d.DebugModeStatus = level
	return d
}

// SetDateFormat sets the date format for log messages in the log rule.
func (d *LogRule) SetDateFormat(format string) *LogRule {
	d.DateFormat = format
	return d
}

// SetLogDate enables or disables the inclusion of the log date in file names for the log rule.
func (d *LogRule) SetLogDate(isLogDate bool) *LogRule {
	d.FileLog.IsDateFile = isLogDate
	return d
}

// SetLogDateFormat sets the date format for log file names in the log rule.
func (d *LogRule) SetLogDateFormat(format string) *LogRule {
	d.FileLog.DateFileFormat = format
	return d
}

// SetConsoleOutput enables or disables console output for log messages in the log rule.
func (d *LogRule) SetConsoleOutput(mode bool) *LogRule {
	d.IsConsoleOutput = mode
	return d
}

// SetDetailedErrorOutput enables or disables detailed error output.
func (d *LogRule) SetDetailedErrorOutput(enabled bool) *LogRule {
	d.DetailedErrorOutput = enabled
	return d
}

// StartAsyncLogging starts a goroutine to handle asynchronous logging.
// It listens for log messages and writes them to the log file and/or console.
func (lr *LogRule) StartAsyncLogging() {
	if !lr.AsyncLog.Enable {
		fmt.Println("[mklog] async logging is disabled")
		return
	}
	go func() {
		for logMessage := range lr.logChannel {
			if lr.FileLog.Enable {
				if err := lr.writeLog(logMessage); err != nil {
					fmt.Println("[mklog] error while writing to log file ", lr.ModuleName, " : ", err)
				}
			}
			if lr.IsConsoleOutput {
				fmt.Println(logMessage) // Output log message to the console.
			}
		}
	}()
}

//#region File

// CreateLogFile initializes the log file for the current LogRule.
// It calls the private createLogFile method and returns any error encountered.
func (d *LogRule) CreateLogFile() error {
	err := d.createLogFile()
	return err
}

// SetIsLogFile enables or disables logging to a file based on the provided boolean value.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetIsLogFile(isLog bool) *LogRule {
	d.FileLog.Enable = isLog
	return d
}

// SetLogFileName sets the name of the log file.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetLogFileName(filename string) *LogRule {
	d.FileLog.FileName = filename
	return d
}

// SetFilePath specifies the path for the log file.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetFilePath(filepath string) *LogRule {
	d.FileLog.FilePath = filepath
	return d
}

// SetLogFileType defines the type (extension) of the log file.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetLogFileType(_type string) *LogRule {
	d.FileLog.FileType = _type
	return d
}

// SetLogDateFileFormat sets the format for date in the log file name.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetLogDateFileFormat(_type string) *LogRule {
	d.FileLog.DateFileFormat = _type
	return d
}

// SetUseTimeFolder enables or disables the use of a time-based folder structure for logs.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetUseTimeFolder(isUseTimeFolder bool) *LogRule {
	d.FileFolder.Enable = isUseTimeFolder
	return d
}

// SetTimeFolderFormat sets the format for the time-based folder structure.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetTimeFolderFormat(format string) *LogRule {
	d.FileFolder.TimeFolderFormat = format
	return d
}

// SetFileFolderPeriod defines the duration for how often the folder structure is updated.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetFileFolderPeriod(period time.Duration) *LogRule {
	d.FileFolder.FileFolderPeriod = period
	return d
}

// SetLimitedFileSize enables or disables the limitation on the log file size.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetLimitedFileSize(isLimited bool) *LogRule {
	d.FileLog.IsLimitedFileSize = isLimited
	return d
}

// SetMaxFileSize specifies the maximum size limit for the log file.
// It returns the updated LogRule instance to allow method chaining.
func (d *LogRule) SetMaxFileSize(size int64) *LogRule {
	d.FileLog.MaxFileSize = size
	return d
}

// InitFiles initializes all log files defined in the Debugger's log rules.
// It iterates over each rule and calls the createLogFile method for each.
func (d *Debugger) InitFiles() *Debugger {
	for _, rule := range d.LogRules {
		for _, v := range rule {
			v.createLogFile()
		}
	}

	return d
}

//#endregion

// #region Channels

// GetSignalChannel returns the signal channel used for interrupt signals.
func (d *LogRule) GetSignalChannel() chan os.Signal {
	return d.signalChannel
}

// GetLogFinishChannel returns the channel used to signal log finishing.
func (d *LogRule) GetLogFinishChannel() chan struct{} {
	return d.logFinishChannel
}

// SetLogChannel set the channel that using for log messages.
func (d *LogRule) SetLogChannel(channel chan string) *LogRule {
	d.logChannel = channel
	return d
}

// GetLogChannel returns the channel used for log messages.
func (d *LogRule) GetLogChannel() chan string {
	return d.logChannel
}

//#endregion

//#region Defaults

// DefaultConsoleLogging creates a Debugger instance with default settings for console logging.
// It initializes log rules with a minimum log level of InfoLevel and a maximum log level of FatalLevel,
// and enables console output and debug mode.
func DefaultConsoleLogging(moduleName string) *Debugger {
	d := &Debugger{
		LogRules: make(map[string][]*LogRule),
	}

	d.NewLogRule(
		moduleName,
		WithMinLevel(InfoLevel),
		WithMaxLevel(FatalLevel),
		WithConsoleOutput(true),
		WithDebugMode(false, InfoLevel),
		WithDateFormat(MKLOG_TimeLogFormatDefault),
		WithForrmatter(MKLOG_FormatterDefault),
	)
	return d
}

// DefaultLogFileSettings creates a Debugger instance with default settings for file logging.
// It initializes log rules similar to DefaultConsoleLogging, but includes file logging settings.
func DefaultLogFileSettings(moduleName string) *Debugger {
	d := &Debugger{
		LogRules: make(map[string][]*LogRule),
	}
	d.NewLogRule(moduleName,
		WithMinLevel(InfoLevel),
		WithMaxLevel(FatalLevel),
		WithConsoleOutput(true),
		WithDebugMode(false, InfoLevel),
		WithDateFormat(MKLOG_TimeLogFormatDefault),
		WithFileLoggingDateFormat(MKLOG_DirDefault, MKLOG_FileNameDefault, MKLOG_FileTypeDefault, MKLOG_TimeFileFormatDefault, true),
	)
	return d
}

// DefaultLogFileAndFolderSettings creates a Debugger instance with default settings for both file and folder logging.
// It initializes log rules with console output, file logging, and a time-based folder structure.
func DefaultLogFileAndFolderSettings(moduleName string) *Debugger {
	d := &Debugger{
		LogRules: make(map[string][]*LogRule),
	}

	d.NewLogRule(
		moduleName,
		WithMinLevel(InfoLevel),
		WithMaxLevel(FatalLevel),
		WithConsoleOutput(true),
		WithDebugMode(false, InfoLevel),
		WithForrmatter(MKLOG_FormatterDefault),
		WithDateFormat(MKLOG_TimeLogFormatDefault),
		WithFileLoggingDateFormat(MKLOG_DirDefault, MKLOG_FileNameDefault, MKLOG_FileTypeDefault, MKLOG_TimeFileFormatDefault, true),
		WithTimeFolder(MKLOG_TimeFolderFormatDefault, MKLOG_FileFolderPeriodDefault, true),
	)
	return d
}

// DefaultSeparateLogAndError creates a Debugger instance with separate logging settings for standard and error logs.
// It sets up two log rules: one for InfoLevel to ErrorLevel and another for ErrorLevel to FatalLevel.
func DefaultSeparateLogAndError(moduleName string) *Debugger {
	d := &Debugger{
		LogRules: make(map[string][]*LogRule),
	}

	d.NewLogRule(
		moduleName,
		WithMinLevel(InfoLevel),
		WithMaxLevel(ErrorLevel),
		WithFileLoggingDateFormat(MKLOG_DirDefault, MKLOG_FileNameDefault, MKLOG_FileTypeDefault, MKLOG_TimeFileFormatDefault, true),
		WithTimeFolder(MKLOG_TimeFolderFormatDefault, MKLOG_FileFolderPeriodDefault, true),
		WithConsoleOutput(true),
		WithDebugMode(false, InfoLevel),
		WithForrmatter(MKLOG_FormatterDefault),
	)

	d.NewLogRule(
		moduleName,
		WithMinLevel(ErrorLevel),
		WithMaxLevel(FatalLevel),
		WithFileLoggingDateFormat(MKLOG_DirDefault, "err", ".err", MKLOG_TimeFileFormatDefault, true),
		WithTimeFolder(MKLOG_TimeFolderFormatDefault, MKLOG_FileFolderPeriodDefault, true),
		WithDateFormat(MKLOG_TimeLogFormatDefault),
		WithDebugMode(false, InfoLevel),
		WithDetailedErrorOutput(true),
		WithForrmatter(MKLOG_FormatterDefault),
	)

	return d
}

//#endregion

//#region Formatters

// SetLogFormatter sets the log formatter for the Debugger instance.
func (d *LogRule) SetLogFormatter(formatter LogFormatter) *LogRule {
	d.LogFormatter = formatter
	return d
}

// SetUserDefinedFormatter sets a user-defined log formatter function and date format.
func (d *LogRule) SetUserDefinedFormatter(formatFunc UserDefinedFormatterFunc) *LogRule {
	d.LogFormatter = UserDefinedFormatter{formatFunc}
	return d
}

// SetCustomLogLevelNames sets custom log level names provided by the user.
func (d *LogRule) SetCustomLogLevelNames(customLogLevelNames map[LogLevel]string) *LogRule {
	d.CustomLogLevelNames = customLogLevelNames
	return d
}

//#endregion

// StringToLogLevel maps a string to a LogLevel.
func StringToLogLevel(level string) (LogLevel, error) {
	level = strings.ToLower(level)
	switch level {
	case "debug", "d":
		return DebugLevel, nil
	case "info", "i":
		return InfoLevel, nil
	case "warning", "w":
		return WarningLevel, nil
	case "error", "e":
		return ErrorLevel, nil
	case "fatal", "f":
		return FatalLevel, nil
	case "trace", "t":
		return TraceLevel, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", level)
	}
}

type Loggable interface {
	GetLogLevelName() string
}

func (l LogLevel) GetLogLevelName() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (lr *LogRule) GetLogLevelName(logLevel LogLevel) string {
	if name, exists := lr.CustomLogLevelNames[logLevel]; exists {
		return name
	}
	return logLevel.GetLogLevelName()
}
