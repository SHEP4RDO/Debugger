package mklog

import (
	"fmt"
	"time"
)

// CustomTrace logs a message at the specified log level and handles error extraction.
// It checks all log rules to determine if the message should be logged based on the rules' conditions.
func (d *Debugger) CustomTrace(logLevel LogLevel, msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, logLevel)

		for _, v := range loggableRules {
			if v.DebugMode && v.DebugModeStatus == TraceLevel {
				v.CurrentLevel = logLevel

				finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
				if v.AsyncLog.Enable {
					v.logChannel <- finalMessage
				} else {
					v.print(finalMessage)
				}
			}
		}
	}
}

// CustomDebug logs a message at the specified log level, similarly to CustomTrace.
// It checks if the log should be output based on the rules defined in LogRules.
func (d *Debugger) CustomDebug(logLevel LogLevel, msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, logLevel)

		for _, v := range loggableRules {
			if v.DebugMode {
				v.CurrentLevel = logLevel

				finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
				if v.AsyncLog.Enable {
					v.logChannel <- finalMessage
				} else {
					v.print(finalMessage)
				}
			}
		}
	}
}

// Custom logs a message at a specified log level, checking the appropriate rules.
// This method is more general and does not have specific conditions like debug mode.
func (d *Debugger) Custom(logLevel LogLevel, msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, logLevel)

		for _, v := range loggableRules {
			v.CurrentLevel = logLevel

			finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
			if v.AsyncLog.Enable {
				v.logChannel <- finalMessage
			} else {
				v.print(finalMessage)
			}
		}
	}
}

// Debug logs a message at the Debug level and checks if it should be output based on the defined rules.
func (d *Debugger) Debug(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, DebugLevel)

		for _, v := range loggableRules {
			if v.DebugMode {
				v.CurrentLevel = DebugLevel

				finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
				if v.AsyncLog.Enable {
					v.logChannel <- finalMessage
				} else {
					v.print(finalMessage)
				}
			}
		}
	}
}

// Trace logs a message at the Trace level, outputting it based on the console and file settings.
func (d *Debugger) Trace(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, TraceLevel)

		for _, v := range loggableRules {
			if v.DebugMode && v.DebugModeStatus == TraceLevel {
				v.CurrentLevel = TraceLevel

				finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
				if v.AsyncLog.Enable {
					v.logChannel <- finalMessage
				} else {
					v.print(finalMessage)
				}
			}
		}
	}
}

// Info logs a message at the Info level, similar to other log methods, checking for applicable rules.
func (d *Debugger) Info(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, InfoLevel)

		for _, v := range loggableRules {
			v.CurrentLevel = InfoLevel

			finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)

			if v.AsyncLog.Enable {
				v.logChannel <- finalMessage
			} else {
				v.print(finalMessage)
			}
		}
	}
}

// Warning logs a message at the Warning level, checking if it should be printed based on the rules.
func (d *Debugger) Warning(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, WarningLevel)

		for _, v := range loggableRules {
			v.CurrentLevel = WarningLevel

			finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
			if v.AsyncLog.Enable {
				v.logChannel <- finalMessage
			} else {
				v.print(finalMessage)
			}
		}
	}
}

// Error logs a message at the Error level, outputting it based on the defined logging rules.
func (d *Debugger) Error(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, ErrorLevel)

		for _, v := range loggableRules {
			v.CurrentLevel = ErrorLevel

			finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
			if v.AsyncLog.Enable {
				v.logChannel <- finalMessage
			} else {
				v.print(finalMessage)
			}
		}
	}
}

// Fatal logs a message at the Fatal level, handling output based on rules set in LogRules.
func (d *Debugger) Fatal(msg string, args ...interface{}) {
	logMessage := fmt.Sprintf(msg, args...)
	err := d.extractError(args...)

	for _, rules := range d.LogRules {
		loggableRules := filterLoggableRules(rules, FatalLevel)

		for _, v := range loggableRules {
			v.CurrentLevel = FatalLevel

			finalMessage := v.prepareMessage(logMessage, v.CurrentLevel, err)
			if v.AsyncLog.Enable {
				v.logChannel <- finalMessage
			} else {
				v.print(finalMessage)
			}
		}
	}
}

// print outputs the final log message to the console and to the log file if enabled.
func (lr *LogRule) print(finalMessage string) {
	if lr.IsConsoleOutput {
		fmt.Println(finalMessage)
	}

	if lr.FileLog.Enable {
		if err := lr.writeLog(finalMessage + "\n"); err != nil {
			fmt.Println("[mklog] Error while writing to log file ", lr.ModuleName, " : ", err)
		}
	}
}

// extractError checks the arguments for any errors and returns the first found error.
func (d *Debugger) extractError(args ...interface{}) error {
	if len(args) > 0 {
		for _, v := range args {
			if e, ok := v.(error); ok {
				return e
			}
		}
	}
	return nil
}

// prepareMessage formats the log message with relevant details including timestamp and log level.
func (lr *LogRule) prepareMessage(logMessage string, logLevel LogLevel, optionalArgs ...interface{}) string {
	logLevelName := lr.GetLogLevelName(logLevel)
	finalMessage := lr.LogFormatter.Format(logMessage, logLevelName, lr.ModuleName, lr.Submodules, time.Now().Format(lr.DateFormat))

	for _, arg := range optionalArgs {
		if detailedErr, ok := arg.(DetailedError); ok {
			finalMessage += detailedErr.ErrorStack()
			break
		}
	}
	return finalMessage
}

// filterLoggableRules filters the log rules to find those that are applicable based on the log level.
func filterLoggableRules(rules []*LogRule, level LogLevel) []*LogRule {
	var loggableRules []*LogRule
	for _, rule := range rules {
		if rule.shouldLog(level) {
			loggableRules = append(loggableRules, rule)
		}
	}
	return loggableRules
}

// shouldLog determines if the log level falls within the rule's specified min and max levels.
func (lr *LogRule) shouldLog(logLevel LogLevel) bool {
	return lr.MinLevel <= logLevel && logLevel <= lr.MaxLevel
}
