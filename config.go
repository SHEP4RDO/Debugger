package mklog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type ConfigParser interface {
	ParseConfig(data []byte, config *Config) error
}

type JSONConfigParser struct{}

func (j *JSONConfigParser) ParseConfig(data []byte, config *Config) error {
	return json.Unmarshal(data, config)
}

type YAMLConfigParser struct{}

func (y *YAMLConfigParser) ParseConfig(data []byte, config *Config) error {
	return yaml.Unmarshal(data, config)
}

type LogFormatterConfig struct {
	Type       string `yaml:"type" json:"type"`
	DateFormat string `yaml:"date_format" json:"date_format"`
}

type LogConfigManager struct {
	parsers               map[string]ConfigParser
	userDefinedFormatters map[string]UserDefinedFormatterFunc
}

type AsyncLogConf struct {
	Enable     bool `yaml:"enable" json:"enable"`
	BufferSize int  `yaml:"buffer_size" json:"buffer_size"`
}

type FolderFileConf struct {
	Enable           bool          `yaml:"enable" json:"enable"`                         // Flag indicating whether to create folders based on time.
	FileFolderPeriod time.Duration `yaml:"file_folder_period" json:"file_folder_period"` // Period for creating folder for log files.
	TimeFolderFormat string        `yaml:"time_folder_format" json:"time_folder_format"` // Format for time folders.
}

type LogFileConf struct {
	DailyLog          bool   `yaml:"daily_log_enable" json:"daily_log_enable"`
	Enable            bool   `yaml:"enable" json:"enable"`                             // Flag indicating whether to log to a file.
	IsLimitedFileSize bool   `yaml:"is_limited_file_size" json:"is_limited_file_size"` // Flag indicating whether to limit file size.
	MaxFileSize       int64  `yaml:"max_file_size" json:"max_file_size"`               // Maximum size of the log file.
	FilePath          string `yaml:"file_path" json:"file_path"`                       // Path to the directory where log files are stored.
	FileName          string `yaml:"file_name" json:"file_name"`                       // Base name of the log file.
	FileType          string `yaml:"file_type" json:"file_type"`                       // Type of the log file (e.g., ".log").
	DateFileFormat    string `yaml:"date_file_format" json:"date_file_format"`
	DetailedError     bool   `yaml:"detailed_error" json:"detailed_error"`
}

type LogRulesConf struct {
	MinLevel         LogLevel           `yaml:"min_level" json:"min_level"`
	MaxLevel         LogLevel           `yaml:"max_level" json:"max_level"`
	CurrentLevel     LogLevel           `yaml:"current_level" json:"current_level"`
	DateFormat       string             `yaml:"date_format" json:"date_format"`
	LogFormatterType LogFormatterConfig `yaml:"log_formatter" json:"log_formatter"`
	ModuleName       string             `yaml:"module_name" json:"module_name"`
	Submodules       []string           `yaml:"submodules" json:"submodules"`
	ConsoleEnable    bool               `yaml:"console_enable" json:"console_enable"`
	IsDebugMod       bool               `yaml:"is_debug_mod" json:"is_debug_mod"`
	DebugModeStatus  LogLevel           `yaml:"debug_mode_status" json:"debug_mode_status"`
	LogFile          LogFileConf        `yaml:"file_log" json:"file_log"`
	FolderFIle       FolderFileConf     `yaml:"folder_file" json:"folder_file"`
	AsyncLog         AsyncLogConf       `yaml:"async_log" json:"async_log"`
}

type Config struct {
	LogRules map[string][]LogRulesConf `yaml:"log_rules" json:"log_rules"`
}

func NewLogConfigManager() *LogConfigManager {
	manager := &LogConfigManager{
		parsers: map[string]ConfigParser{
			".json": &JSONConfigParser{},
			".yaml": &YAMLConfigParser{},
			".yml":  &YAMLConfigParser{},
		},
		userDefinedFormatters: make(map[string]UserDefinedFormatterFunc),
	}
	return manager
}

func (m *LogConfigManager) RegisterParser(fileType string, parser ConfigParser) {
	m.parsers[fileType] = parser
}

func (m *LogConfigManager) RegisterUserDefinedFormatter(name string, formatFunc UserDefinedFormatterFunc) {
	m.userDefinedFormatters[name] = formatFunc
}

func (m *LogConfigManager) LoadConfig(filePath string) (*Debugger, error) {
	ext := filepath.Ext(filePath)
	parser, ok := m.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("[mklog] unsupported config format: %s", ext)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("[mklog] failed to read config file: %w", err)
	}

	var config Config
	if err := parser.ParseConfig(data, &config); err != nil {
		return nil, fmt.Errorf("[mklog] failed to parse config file: %w", err)
	}

	debugger := &Debugger{
		LogRules: make(map[string][]*LogRule),
	}

	for ruleName, rules := range config.LogRules {
		for _, rule := range rules {

			formatter, err := rule.getFormatter(m.userDefinedFormatters)
			if err != nil {
				return nil, fmt.Errorf("[mklog] failed to get formatter: %w", err)
			}

			if rule.AsyncLog.Enable && rule.AsyncLog.BufferSize <= 0 {
				fmt.Println("[mklog] Buffersize set to default value")
				rule.AsyncLog.BufferSize = MKLOG_BufferSizeDefault
			}

			if rule.FolderFIle.Enable && rule.LogFile.Enable {
				if err := rule.checkFilePath(); err != nil {
					return nil, fmt.Errorf("[mklog] failed to check file path: %w", err)
				}

				if err := rule.checkFolderSettings(); err != nil {
					return nil, fmt.Errorf("[mklog] failed to check folder settings: %w", err)
				}

				debugger.NewLogRule(ruleName,
					WithMinLevel(rule.MinLevel),
					WithMaxLevel(rule.MaxLevel),
					WithConsoleOutput(rule.ConsoleEnable),
					WithDebugMode(rule.IsDebugMod, rule.DebugModeStatus),
					WithDetailedErrorOutput(rule.LogFile.DetailedError),
					WithDateFormat(rule.DateFormat),
					WithForrmatter(formatter),
					WithTimeFolder(rule.FolderFIle.TimeFolderFormat, rule.FolderFIle.FileFolderPeriod, rule.FolderFIle.Enable),
					WithFileLoggingDateFormat(rule.LogFile.FilePath, rule.LogFile.FileName, rule.LogFile.FileType, rule.LogFile.DateFileFormat, rule.LogFile.DailyLog),
					WithAsyncLog(rule.AsyncLog.Enable, rule.AsyncLog.BufferSize),
				)
			} else if rule.LogFile.Enable {

				if err := rule.checkFilePath(); err != nil {
					return nil, fmt.Errorf("[mklog] failed to check file path: %w", err)
				}

				debugger.NewLogRule(ruleName,
					WithMinLevel(rule.MinLevel),
					WithMaxLevel(rule.MaxLevel),
					WithConsoleOutput(rule.ConsoleEnable),
					WithDebugMode(rule.IsDebugMod, rule.DebugModeStatus),
					WithDetailedErrorOutput(rule.LogFile.DetailedError),
					WithDateFormat(rule.DateFormat),
					WithForrmatter(formatter),
					WithFileLoggingDateFormat(rule.LogFile.FilePath, rule.LogFile.FileName, rule.LogFile.FileType, rule.LogFile.DateFileFormat, rule.LogFile.DailyLog),
					WithAsyncLog(rule.AsyncLog.Enable, rule.AsyncLog.BufferSize),
				)
			} else if rule.ConsoleEnable {
				debugger.NewLogRule(ruleName,
					WithMinLevel(rule.MinLevel),
					WithMaxLevel(rule.MaxLevel),
					WithConsoleOutput(rule.ConsoleEnable),
					WithDetailedErrorOutput(rule.LogFile.DetailedError),
					WithDebugMode(rule.IsDebugMod, rule.DebugModeStatus),
					WithDateFormat(rule.DateFormat),
					WithForrmatter(formatter),
					WithAsyncLog(rule.AsyncLog.Enable, rule.AsyncLog.BufferSize),
				)
			}
		}
	}

	return debugger, nil
}

func (rule *LogRulesConf) checkFilePath() error {
	if rule.LogFile.Enable {
		if rule.LogFile.FilePath != "" {
			if rule.LogFile.FileName == "" || rule.LogFile.FileType == "" {
				dir := filepath.Dir(rule.LogFile.FilePath)

				if rule.LogFile.FileName == "" {
					fileNameWithExt := filepath.Base(rule.LogFile.FilePath)
					rule.LogFile.FileName = fileNameWithoutExt(fileNameWithExt)
				}
				if rule.LogFile.FileType == "" {
					rule.LogFile.FileType = filepath.Ext(rule.LogFile.FilePath)
				}

				if rule.LogFile.FileName == "" || rule.LogFile.FileName == "." {
					rule.LogFile.FileName = MKLOG_FileNameDefault
				}

				if rule.LogFile.FileType == "" {
					rule.LogFile.FileType = MKLOG_FileTypeDefault
				}

				rule.LogFile.FilePath = dir
			}
		} else {
			if rule.LogFile.FileName == "" {
				rule.LogFile.FileName = MKLOG_FileNameDefault
			}
			if rule.LogFile.FileType == "" {
				rule.LogFile.FileType = MKLOG_FileTypeDefault
			}
			if rule.LogFile.FilePath == "" {
				rule.LogFile.FilePath = "./" + MKLOG_DirDefault
			}
		}

		if rule.LogFile.FileName == "" || rule.LogFile.FileType == "" || rule.LogFile.FilePath == "" {
			return fmt.Errorf("failed to set file details: file name, file type, and path must be specified")
		}
	}

	return nil
}

func (rule *LogRulesConf) checkFolderSettings() error {
	if rule.FolderFIle.Enable {
		if rule.FolderFIle.TimeFolderFormat == "" {
			rule.FolderFIle.TimeFolderFormat = MKLOG_TimeFolderFormatDefault
			fmt.Println("[mklog] TimeFolderFormat is not specified. Using default:", MKLOG_TimeFolderFormatDefault)
		}
		if rule.FolderFIle.FileFolderPeriod == 0 {
			rule.FolderFIle.FileFolderPeriod = MKLOG_FileFolderPeriodDefault
			fmt.Println("[mklog] FileFolderPeriod is not specified. Using default:", MKLOG_FileFolderPeriodDefault)
		}
	}

	if rule.FolderFIle.TimeFolderFormat == "" || rule.FolderFIle.FileFolderPeriod == 0 {
		return fmt.Errorf("[mklog] failed to set folder settings: TimeFolderFormat and FileFolderPeriod must be specified when Enable is true")
	}

	return nil
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (rule *LogRulesConf) getFormatter(userDefinedFormatters map[string]UserDefinedFormatterFunc) (LogFormatter, error) {
	formatterType := strings.ToLower(rule.LogFormatterType.Type)
	var formatter LogFormatter

	switch formatterType {
	case "plaintextformatter", "plaintext", "plain", "text", "simple":
		formatter = PlainTextFormatter{dateFormat: rule.DateFormat}
		return formatter, nil
	case "jsonformatter", "json":
		formatter = JSONFormatter{dateFormat: rule.DateFormat}
		return formatter, nil
	case "yamlformatter", "yaml", "yml":
		formatter = YAMLFormatter{dateFormat: rule.DateFormat}
		return formatter, nil
	case "xmlformatter", "xml":
		formatter = XMLFormatter{dateFormat: rule.DateFormat}
		return formatter, nil
	default:
		if formatFunc, exists := userDefinedFormatters[formatterType]; exists {
			formatter = UserDefinedFormatter{formatFunc: formatFunc}
		} else {
			return nil, fmt.Errorf("[mklog] unsupported log formatter type: %s", rule.LogFormatterType)
		}
	}

	return nil, fmt.Errorf("unsupported log formatter type: %s", rule.LogFormatterType)
}
