package mklog

import "time"

type Option func(*LogRule)

// WithMinLevel sets the minimum logging level.
func WithMinLevel(minLevel LogLevel) Option {
	return func(lr *LogRule) {
		lr.MinLevel = minLevel
	}
}

// WithMaxLevel sets the maximum logging level.
func WithMaxLevel(maxLevel LogLevel) Option {
	return func(lr *LogRule) {
		lr.MaxLevel = maxLevel
	}
}

// WithCurrentLevel sets the current logging level.
func WithCurrentLevel(currentLevel LogLevel) Option {
	return func(lr *LogRule) {
		lr.CurrentLevel = currentLevel
	}
}

// WithFileLogging enables file logging and sets file parameters.
func WithFileLogging(filePath, fileName, fileType string) Option {
	return func(lr *LogRule) {
		lr.FileLog.Enable = true
		lr.FileLog.FilePath = filePath
		lr.FileLog.FileName = fileName
		lr.FileLog.FileType = fileType
	}
}

// WithFileLoggingDateFormat enables file logging with date formatting options.
func WithFileLoggingDateFormat(filePath, fileName, fileType, dateFormat string, isDaily bool) Option {
	return func(lr *LogRule) {
		lr.FileLog.Enable = true
		lr.FileLog.FilePath = filePath
		lr.FileLog.FileName = fileName
		lr.FileLog.FileType = fileType
		lr.FileLog.IsDateFile = isDaily
		lr.FileLog.DateFileFormat = dateFormat
	}
}

// WithTimeFolder enables folder organization by time period.
func WithTimeFolder(timeFolderFormat string, folderPeriod time.Duration, isFolderTime bool) Option {
	return func(lr *LogRule) {
		lr.FileFolder.Enable = isFolderTime
		lr.FileFolder.TimeFolderFormat = timeFolderFormat
		lr.FileFolder.FileFolderPeriod = folderPeriod
	}
}

// WithDateFormat sets the date format for logs.
func WithDateFormat(format string) Option {
	return func(lr *LogRule) {
		lr.DateFormat = format
	}
}

// WithDebugMode enables debug mode with a specified debug level.
func WithDebugMode(debugMode bool, debugLevel LogLevel) Option {
	return func(lr *LogRule) {
		lr.DebugMode = debugMode
		lr.DebugModeStatus = debugLevel
	}
}

// WithDetailedErrorOutput enables detailed error output, including stack traces.
func WithDetailedErrorOutput(enable bool) Option {
	return func(rule *LogRule) {
		rule.DetailedErrorOutput = enable
	}
}

// WithLogFormatter sets a custom log formatter.
func WithLogFormatter(formatter LogFormatter) Option {
	return func(lr *LogRule) {
		lr.LogFormatter = formatter
	}
}

// WithMaxFileSize sets the maximum file size for log files.
func WithMaxFileSize(maxFileSize int64) Option {
	return func(lr *LogRule) {
		lr.FileLog.MaxFileSize = maxFileSize
		lr.FileLog.IsLimitedFileSize = true
	}
}

// WithConsoleOutput enables or disables console output for logs.
func WithConsoleOutput(consoleOutput bool) Option {
	return func(lr *LogRule) {
		lr.IsConsoleOutput = consoleOutput
	}
}

// WithFormatter sets a log formatter (duplicate of WithLogFormatter).
func WithForrmatter(formatter LogFormatter) Option {
	return func(lr *LogRule) {
		lr.LogFormatter = formatter
	}
}

// WithAsyncLog enables asynchronous logging with a specified buffer size.
func WithAsyncLog(enable bool, bufferSize int) Option {
	return func(lr *LogRule) {
		lr.AsyncLog.Enable = enable
		lr.AsyncLog.BufferSize = bufferSize
	}
}
