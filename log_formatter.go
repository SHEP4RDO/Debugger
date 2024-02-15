package mklog

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// LogFormatter is an interface that defines the methods required to format log messages.
type LogFormatter interface {
	// Format formats the log message using the specified parameters and returns the formatted log string.
	Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string
}

// PlainTextFormatter is a LogFormatter implementation that formats log messages in plain text.
type PlainTextFormatter struct {
	dateFormat string
}

// Format formats the log message in plain text.
func (f PlainTextFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
	if len(submodules) > 0 {
		return fmt.Sprintf("%s | %s | [%s] - %v: %s\n",
			timestamp,
			logLevel,
			moduleName,
			submodules,
			logMessage,
		)
	} else {
		return fmt.Sprintf("%s | %s | [%s] : %s\n",
			timestamp,
			logLevel,
			moduleName,
			logMessage,
		)
	}
}

// JSONFormatter is a LogFormatter implementation that formats log messages in JSON.
type JSONFormatter struct {
	dateFormat string
}

// Format formats the log message in JSON.
func (f JSONFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {

	logData := make(map[string]interface{})
	if len(submodules) > 0 {
		logData = map[string]interface{}{
			"timestamp":  timestamp,
			"logLevel":   logLevel,
			"moduleName": moduleName,
			"submodules": submodules,
			"logMessage": logMessage,
		}
	} else {
		logData = map[string]interface{}{
			"timestamp":  timestamp,
			"logLevel":   logLevel,
			"moduleName": moduleName,
			"logMessage": logMessage,
		}
	}

	logJSON, _ := json.Marshal(logData)
	return string(logJSON) + "\n"
}

// UserDefinedFormatterFunc is a function type for user-defined log message formatting.
type UserDefinedFormatterFunc func(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string

// UserDefinedFormatter is a LogFormatter implementation that allows users to define their log message formatting.
type UserDefinedFormatter struct {
	formatFunc UserDefinedFormatterFunc
}

// Format formats the log message using a user-defined function.
func (f UserDefinedFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
	return f.formatFunc(logMessage, logLevel, moduleName, submodules, timestamp)
}

// XMLFormatter is a LogFormatter implementation that formats log messages in XML.
type XMLFormatter struct {
	dateFormat string
}

// Format formats the log message in XML.
func (f XMLFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
	if len(submodules) > 0 {
		return fmt.Sprintf("<LogEntry>\n"+
			"    <Timestamp>%s</Timestamp>\n"+
			"    <LogLevel>%s</LogLevel>\n"+
			"    <ModuleName>%s</ModuleName>\n"+
			"    <Submodules>%v</Submodules>\n"+
			"    <Message>%s</Message>\n"+
			"</LogEntry>\n",
			timestamp,
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
			timestamp,
			logLevel,
			moduleName,
			logMessage,
		)
	}
}

// YAMLFormatter is a LogFormatter implementation that formats log messages in YAML.
type YAMLFormatter struct {
	dateFormat string
}

// Format formats the log message in YAML.
func (f YAMLFormatter) Format(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
	logData := make(map[string]interface{})
	logData["timestamp"] = timestamp
	logData["logLevel"] = logLevel
	logData["moduleName"] = moduleName

	if len(submodules) > 0 {
		logData["submodules"] = submodules
	}

	logData["logMessage"] = logMessage

	logYAML, _ := yaml.Marshal(logData)
	return string(logYAML) + "\n"
}
