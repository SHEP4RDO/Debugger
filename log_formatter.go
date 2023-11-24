package mklog

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// LogFormatter is an interface that defines the methods required to format log messages.
type LogFormatter interface {
	// Format formats the log message using the specified parameters and returns the formatted log string.
	Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string

	// SetLogDateFormat sets the date format used in the log messages.
	SetLogDateFormat(format string)

	// GetLogDateFormat returns the current date format used in the log messages.
	GetLogDateFormat() string
}

// PlainTextFormatter is a LogFormatter implementation that formats log messages in plain text.
type PlainTextFormatter struct {
	dateFormat string
}

// Format formats the log message in plain text.
func (f PlainTextFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string {
	if f.dateFormat == "" {
		f.dateFormat = "2006-01-02 15:04:05.000"
	}
	if len(submodules) > 0 {
		return fmt.Sprintf("%s | %s | [%s] - %v: %s\n",
			timestamp.Format(f.dateFormat),
			logLevel,
			moduleName,
			submodules,
			logMessage,
		)
	} else {
		return fmt.Sprintf("%s | %s | [%s] : %s\n",
			timestamp.Format(f.dateFormat),
			logLevel,
			moduleName,
			logMessage,
		)
	}
}

// SetLogDateFormat sets the date format for plain text log messages.
func (f PlainTextFormatter) SetLogDateFormat(format string) {
	f.dateFormat = format
}

// GetLogDateFormat returns the current date format for plain text log messages.
func (f PlainTextFormatter) GetLogDateFormat() string {
	return f.dateFormat
}

// JSONFormatter is a LogFormatter implementation that formats log messages in JSON.
type JSONFormatter struct {
	dateFormat string
}

// Format formats the log message in JSON.
func (f JSONFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string {
	if f.dateFormat == "" {
		f.dateFormat = "2006-01-02 15:04:05.000"
	}
	logData := make(map[string]interface{})
	if len(submodules) > 0 {
		logData = map[string]interface{}{
			"timestamp":  timestamp.Format(f.dateFormat),
			"logLevel":   logLevel,
			"moduleName": moduleName,
			"submodules": submodules,
			"logMessage": logMessage,
		}
	} else {
		logData = map[string]interface{}{
			"timestamp":  timestamp.Format(f.dateFormat),
			"logLevel":   logLevel,
			"moduleName": moduleName,
			"logMessage": logMessage,
		}
	}

	logJSON, _ := json.Marshal(logData)
	return string(logJSON) + "\n"
}

// SetLogDateFormat sets the date format for JSON log messages.
func (f JSONFormatter) SetLogDateFormat(format string) {
	f.dateFormat = format
}

// GetLogDateFormat returns the current date format for JSON log messages.
func (f JSONFormatter) GetLogDateFormat() string {
	return f.dateFormat
}

// UserDefinedFormatterFunc is a function type for user-defined log message formatting.
type UserDefinedFormatterFunc func(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string

// UserDefinedFormatter is a LogFormatter implementation that allows users to define their log message formatting.
type UserDefinedFormatter struct {
	formatFunc UserDefinedFormatterFunc
	dateFormat string
}

// Format formats the log message using a user-defined function.
func (f UserDefinedFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string {
	if f.dateFormat == "" {
		f.dateFormat = "2006-01-02 15:04:05.000"
	}
	return f.formatFunc(logMessage, logLevel, moduleName, submodules, timestamp.Format(f.dateFormat))
}

// SetLogDateFormat sets the date format for user-defined log messages.
func (f UserDefinedFormatter) SetLogDateFormat(format string) {
	f.dateFormat = format
}

// GetLogDateFormat returns the current date format for user-defined log messages.
func (f UserDefinedFormatter) GetLogDateFormat() string {
	return f.dateFormat
}

// XMLFormatter is a LogFormatter implementation that formats log messages in XML.
type XMLFormatter struct {
	dateFormat string
}

// Format formats the log message in XML.
func (f XMLFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string {
	if f.dateFormat == "" {
		f.dateFormat = "2006-01-02 15:04:05.000"
	}

	if len(submodules) > 0 {
		return fmt.Sprintf("<LogEntry>\n"+
			"    <Timestamp>%s</Timestamp>\n"+
			"    <LogLevel>%s</LogLevel>\n"+
			"    <ModuleName>%s</ModuleName>\n"+
			"    <Submodules>%v</Submodules>\n"+
			"    <Message>%s</Message>\n"+
			"</LogEntry>\n",
			timestamp.Format(f.dateFormat),
			logLevel,
			moduleName,
			submodules,
			logMessage,
		)
	} else {
		return fmt.Sprintf("<LogEntry>\n"+
			"    <Timestamp>%s</Timestamp>\n"+
			"    <LogLevel>%s</LogLevel>\n"+
			"    <ModuleName>%s</ModuleName>\n"+
			"    <Message>%s</Message>\n"+
			"</LogEntry>\n",
			timestamp.Format(f.dateFormat),
			logLevel,
			moduleName,
			logMessage,
		)
	}
}

// SetLogDateFormat sets the date format for XML log messages.
func (f XMLFormatter) SetLogDateFormat(format string) {
	f.dateFormat = format
}

// GetLogDateFormat returns the current date format for XML log messages.
func (f XMLFormatter) GetLogDateFormat() string {
	return f.dateFormat
}

// YAMLFormatter is a LogFormatter implementation that formats log messages in YAML.
type YAMLFormatter struct {
	dateFormat string
}

// Format formats the log message in YAML.
func (f YAMLFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp time.Time) string {
	if f.dateFormat == "" {
		f.dateFormat = "2006-01-02 15:04:05.000"
	}

	logData := make(map[string]interface{})
	logData["timestamp"] = timestamp.Format(f.dateFormat)
	logData["logLevel"] = logLevel
	logData["moduleName"] = moduleName

	if len(submodules) > 0 {
		logData["submodules"] = submodules
	}

	logData["logMessage"] = logMessage

	logYAML, _ := yaml.Marshal(logData)
	return string(logYAML) + "\n"
}

// SetLogDateFormat sets the date format for YAML log messages.
func (f YAMLFormatter) SetLogDateFormat(format string) {
	f.dateFormat = format
}

// GetLogDateFormat returns the current date format for YAML log messages.
func (f YAMLFormatter) GetLogDateFormat() string {
	return f.dateFormat
}
