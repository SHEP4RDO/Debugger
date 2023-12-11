package mklog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// log represents the logging configuration and functionality.
type log struct {
	isToFile        bool   // Flag indicating whether to log to a file.
	isDateFile      bool   // Flag indicating whether to include the date in the log file name.
	FilePath        string // Path to the directory where log files are stored.
	File            *os.File
	DateFileFormat  string // Date format used in the log file name.
	FileName        string // Base name of the log file.
	CurrentFileName string // Current full name of the log file.
	FileType        string // Type of the log file (e.g., ".log").
}

// createLogFile initializes and opens the log file if logging to a file is enabled.
func (d *Debugger) createLogFile() error {
	if d.log.isToFile {
		if err := os.MkdirAll(d.log.FilePath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		fileName := filepath.Join(d.log.FilePath, d.log.FileName+d.log.FileType)
		if d.log.isDateFile {
			fileName = filepath.Join(d.log.FilePath, time.Now().Format(d.log.DateFileFormat)+"_"+d.log.FileName)
		}

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		d.log.File = file
		d.log.CurrentFileName = fileName
	}

	go func() {
		select {
		case <-d.signalChannel:
			d.closeLogFile()
		case <-d.logFinishChannel:
		}
	}()

	return nil
}

// closeLogFile closes the log file and signals the log finishing channel.
func (d *Debugger) closeLogFile() {
	if d.log.File != nil {
		close(d.logFinishChannel)
		d.log.File.Close()
		d.log.File = nil
	}
}

// writeLog writes the provided log message to the log file if logging to a file is enabled.
func (l *log) writeLog(msg string) error {
	if l.File != nil {
		_, err := l.File.WriteString(msg)
		return err
	}
	return fmt.Errorf("log file is not open")
}
