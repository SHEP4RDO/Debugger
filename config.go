package mklog

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfigProvider is an interface that defines methods to retrieve configuration parameters for logging.
type ConfigProvider interface {
	IsDebugMode() bool         // IsDebugMode returns whether debug mode is enabled.
	DateFormat() string        // DateFormat returns the format for log timestamps.
	IsConsoleOutput() bool     // IsConsoleOutput returns whether console output is enabled.
	DetailedErrorOutput() bool // DetailedErrorOutput returns whether detailed error output is enabled.
	LogDateFile() bool         // LogDateFile returns whether log date file is enabled.
	LogDateFileFormat() string // LogDateFileFormat returns the format for the log date file timestamps.
	LogFilePath() string       // LogFilePath returns the path for log files.
	LogFileName() string       // LogFileName returns the name for log files.
	LogFileType() string       // LogFileType returns the type (extension) for log files.
	LogFormat() string         // LogFormat returns the desired log format (e.g., PlainText, JSON, XML, UserDefined).
	UserFormat() string        // UserFormat returns the custom log format if LogFormat is UserDefined.

}

// SetByConfigProvider initializes a Debugger instance using the provided ConfigProvider.
func SetByConfigProvider(provider ConfigProvider, moduleName string, submodules ...string) (*Debugger, error) {
	debugger := NewDebugLogger(moduleName, submodules...)
	if provider.IsDebugMode() {
		debugger.SetDebugMode(provider.IsDebugMode())
	}
	if dateFormat := provider.DateFormat(); dateFormat != "" {
		debugger.SetDateFormat(dateFormat)
	}
	if provider.IsConsoleOutput() {
		debugger.SetConsoleOutput(provider.IsConsoleOutput())
	}
	if provider.DetailedErrorOutput() {
		debugger.SetDetailedErrorOutput(provider.DetailedErrorOutput())
	}

	if logFormat := provider.LogFormat(); logFormat != "" {
		err := setFormatterProvider(debugger, provider.DateFormat(), logFormat, provider.UserFormat())
		if err != nil {
			return nil, err
		}
	}
	if provider.LogDateFile() {
		debugger.SetLogDate(provider.LogDateFile())
		if logFileName := provider.LogFileName(); logFileName == "" || logFileName == " " {
			return nil, errors.New("logFileName is required when logDateFile is true")
		}
		if logFilePath := provider.LogFilePath(); logFilePath == "" {
			debugger.SetDefaultLogPath()
		} else {
			debugger.SetLogPath(logFilePath)
			debugger.SetLogFileName(provider.LogFileName())
		}
		if logFileType := provider.LogFileType(); logFileType != "" {
			debugger.SetLogFileType(logFileType)
		} else {
			debugger.SetLogFileType(".log")
		}
	} else if provider.IsConsoleOutput() {
		debugger.Error("no set file name and path to log file. you need at least set a file name log.", nil)
	}
	if err := debugger.createLogFile(); err != nil {
		return nil, err
	}

	return debugger, nil
}

// setFormatterProvider sets the log formatter based on the specified log format.
// It takes the file path, file name, module name, and submodules (optional) as parameters.
func setFormatterProvider(debugger *Debugger, dateFormat, format, userFormat string) error {
	lowercaseFormat := strings.ToLower(format)
	switch lowercaseFormat {
	case "plaintext":
		debugger.SetLogFormatter(PlainTextFormatter{dateFormat: format})
	case "json":
		debugger.SetLogFormatter(JSONFormatter{dateFormat: format})
	case "xml":
		debugger.SetLogFormatter(XMLFormatter{dateFormat: format})
	case "userdefined", "custom":
		{
			if userFormat != "" {
				debugger.logFormatter = UserDefinedFormatter{
					formatFunc: func(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
						return fmt.Sprintf(userFormat,
							timestamp, logLevel, moduleName, submodules, logMessage,
						)
					},
				}
				return nil
			} else {
				return errors.New("userFormat must be provided for userdefined log format")
			}
		}
	default:
		return errors.New("unsupported log format")
	}
	return nil
}

// Config is a struct that represents the configuration settings for logging.
type Config struct {
	isDebugMode         bool   `yaml:"isDebugMode" json:"isDebugMode" xml:"isDebugMode"`
	DateFormat          string `yaml:"dateFormat" json:"dateFormat" xml:"dateFormat"`
	IsConsoleOutput     bool   `yaml:"isConsoleOutput" json:"isConsoleOutput" xml:"isConsoleOutput"`
	DetailedErrorOutput bool   `yaml:"detailedErrorOutput" json:"detailedErrorOutput" xml:"detailedErrorOutput"`

	LogDateFile       bool   `yaml:"logDateFile" json:"logDateFile" xml:"logDateFile"`
	LogDateFileFormat string `yaml:"logDateFileFormat" json:"logDateFileFormat" xml:"logDateFileFormat"`
	LogFilePath       string `yaml:"logFilePath" json:"logFilePath" xml:"logFilePath"`
	LogFileName       string `yaml:"logFileName" json:"logFileName" xml:"logFileName"`
	LogFileType       string `yaml:"logFileType" json:"logFileType" xml:"logFileType"`

	LogFormat  string `json:"logFormat" yaml:"logFormat" xml:"logFormat"`
	UserFormat string `json:"userFormat" yaml:"userFormat" xml:"userFormat"`
}

// SetByConfig initializes a Debugger instance using configuration settings from a file.
// It takes the file path, file name, module name, and submodules (optional) as parameters.
// The configuration file format is determined by the file extension.
// Supported formats: YAML (.yaml, .yml), JSON (.json), XML (.xml).
// Returns a Debugger instance and an error if any.
func SetByConfig(path, fileName, moduleName string, submodules ...string) (*Debugger, error) {
	conf, err := readConfig(path, fileName)
	if err != nil {
		return nil, err
	}

	debugger := NewDebugLogger(moduleName, submodules...)

	if conf.isDebugMode {
		debugger.SetDebugMode(conf.isDebugMode)
	}

	if conf.DateFormat != "" {
		debugger.SetDateFormat(conf.DateFormat)
	}

	if conf.IsConsoleOutput {
		debugger.SetConsoleOutput(conf.IsConsoleOutput)
	}

	if conf.DetailedErrorOutput {
		debugger.SetDetailedErrorOutput(conf.DetailedErrorOutput)
	}

	if conf.LogFormat != "" {
		err := conf.setFormatter(debugger)
		if err != nil {
			return nil, err
		}
	}
	if conf.LogDateFile {
		debugger.SetLogDate(conf.LogDateFile)
		if conf.LogFileName == "" || conf.LogFileName == " " {
			return nil, errors.New("logFileName is required when logDateFile is true")
		}
		if conf.LogFilePath == "" {
			debugger.SetDefaultLogPath()
		} else {
			debugger.SetLogPath(conf.LogFilePath)
			debugger.SetLogFileName(conf.LogFileName)
		}
		if conf.LogFileType != "" {
			debugger.SetLogFileType(conf.LogFileType)
		} else {
			debugger.SetLogFileType(".log")
		}
	} else if conf.IsConsoleOutput {
		debugger.Error("no set file name and path to log file. you need at least set a file name log.", nil)
	}
	if err := debugger.createLogFile(); err != nil {
		return nil, err
	}
	return debugger, nil
}

// readConfig reads configuration settings from the specified file and returns a Config instance.
// The file format is determined by the file extension.
// Supported formats: YAML (.yaml, .yml), JSON (.json), XML (.xml).
// Returns a Config instance and an error if any.
func readConfig(path, fileName string) (*Config, error) {
	if path == "" {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return nil, fmt.Errorf("failed to read current directory: %w", err)
		}

		for _, file := range files {
			if strings.EqualFold(file.Name(), fileName+".yaml") || strings.EqualFold(file.Name(), fileName+".yml") || strings.EqualFold(file.Name(), fileName+".json") || strings.EqualFold(file.Name(), fileName+".xml") {
				path = file.Name()
				break
			}
		}
	}

	if path == "" {
		return nil, errors.New("no config file specified")
	}

	ext := filepath.Ext(path)

	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".xml":
		err = xml.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}

// setFormatter sets the log formatter for the Debugger based on the specified log format in the configuration.
// Returns an error if the log format is unsupported or if user-defined format is selected without providing userFormat.
func (conf *Config) setFormatter(debugger *Debugger) error {
	switch conf.LogFormat {
	case "json", "JSON":
		debugger.SetLogFormatter(JSONFormatter{dateFormat: conf.DateFormat})
		return nil
	case "xml", "XML":
		debugger.SetLogFormatter(XMLFormatter{dateFormat: conf.DateFormat})
		return nil
	case "plaintext", "Plaintext", "PlainText":
		debugger.SetLogFormatter(PlainTextFormatter{dateFormat: conf.DateFormat})
		return nil
	case "userdefined", "Userdefined", "custom", "Custom":

		if conf.UserFormat != "" {

			debugger.logFormatter = UserDefinedFormatter{
				formatFunc: func(logMessage string, logLevel string, moduleName string, submodules []string, timestamp string) string {
					return fmt.Sprintf(conf.UserFormat,
						timestamp, logLevel, moduleName, submodules, logMessage,
					)
				},
			}
			return nil
		} else {
			return errors.New("userFormat must be provided for userdefined log format")
		}
	default:
		return fmt.Errorf("unsupported log format: %s", conf.LogFormat)
	}
}
