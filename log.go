package mklog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// log represents the logging configuration and functionality.
type FileLog struct {
	// bools
	Enable            bool `json:"enable" yaml:"enable"`                             // Flag indicating whether to log to a file.
	IsDateFile        bool `json:"is_date_file" yaml:"is_date_file"`                 // Flag indicating whether to include the date in the log file name.
	IsLimitedFileSize bool `json:"is_limited_file_size" yaml:"is_limited_file_size"` // Flag indicating whether to limit the file size.

	// files
	File            *os.File `json:"-" yaml:"-"`                                 // Pointer to the log file (ignored in configuration).
	MaxFileSize     int64    `json:"max_file_size" yaml:"max_file_size"`         // Maximum size of the log file.
	FileName        string   `json:"file_name" yaml:"file_name"`                 // Base name of the log file.
	FilePath        string   `json:"file_path" yaml:"file_path"`                 // Path to the directory where log files are stored.
	CurrentFileName string   `json:"current_file_name" yaml:"current_file_name"` // Current full name of the log file.
	FileType        string   `json:"file_type" yaml:"file_type"`                 // Type of the log file (e.g., ".log").
	DateFileFormat  string   `json:"date_file_format" yaml:"date_file_format"`   // Date format used in the log file name.
}

type FileFolder struct {
	Enable           bool          `json:"enable" yaml:"enable"`                         // Flag indicating whether to enable folder logging.
	TimeFolderFormat string        `json:"time_folder_format" yaml:"time_folder_format"` // Time format for log folders.
	FileFolderPeriod time.Duration `json:"file_folder_period" yaml:"file_folder_period"` // Period for creating new folders for log files.
}

// createLogFile initializes and opens the log file if logging to a file is enabled.
func (d *LogRule) createLogFile() error {
	if d.FileLog.Enable {
		// Create the log directory if it does not exist.
		if err := os.MkdirAll(d.FileLog.FilePath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		var logFolder string
		// Determine whether to create a time-based folder for log files.
		if d.FileFolder.Enable {
			currentTime := time.Now()
			var folderName string

			// Format folder name based on the specified time period.
			if d.FileFolder.FileFolderPeriod < time.Hour {
				folderName = currentTime.Format(d.FileFolder.TimeFolderFormat)
			} else {
				folderName = currentTime.Truncate(d.FileFolder.FileFolderPeriod).Format(d.FileFolder.TimeFolderFormat)
			}

			// Create the time-based folder for log files.
			logFolder = filepath.Join(d.FileLog.FilePath, folderName)
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create time folder: %w", err)
			}
		} else {
			logFolder = d.FileLog.FilePath // Use the main log directory.
		}

		var fileName string
		// Determine the log file name based on the date settings.
		if d.FileLog.IsDateFile {
			dateStr := time.Now().Format(d.FileLog.DateFileFormat)
			fileName = filepath.Join(logFolder, fmt.Sprintf("%s_%s%s", dateStr, d.FileLog.FileName, d.FileLog.FileType))
		} else {
			fileName = filepath.Join(logFolder, fmt.Sprintf("%s%s", d.FileLog.FileName, d.FileLog.FileType))
		}

		// Open the log file for writing.
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		d.FileLog.File = file
		d.FileLog.CurrentFileName = fileName
	}
	return nil
}

// CloseLogFile closes the log file and signals the log finishing channel.
func (d *LogRule) CloseLogFile() {
	if d.FileLog.File != nil {
		close(d.logFinishChannel) // Signal that logging has finished.
		d.FileLog.File.Close()    // Close the log file.
		d.FileLog.File = nil      // Clear the file pointer.
	}
}

// writeLog writes the provided log message to the log file if logging to a file is enabled.
func (d *LogRule) writeLog(msg string) error {
	if d.FileLog.File != nil {
		now := time.Now()
		var fileName string

		// Determine the appropriate log file name.
		if d.FileLog.IsDateFile {
			fileDate := now.Format("02.01.2006")
			fileName = fmt.Sprintf("%s_%s%s", fileDate, d.FileLog.FileName, d.FileLog.FileType)
		} else {
			fileName = fmt.Sprintf("%s%s", d.FileLog.FileName, d.FileLog.FileType)
		}

		// Check if the log file name has changed and create a new log file if necessary.
		if fileName != d.FileLog.CurrentFileName {
			d.FileLog.CurrentFileName = fileName
			if err := d.createLogFile(); err != nil {
				return err
			}
		}

		// Check if the log file size limit is enabled and trim if necessary.
		if d.FileLog.IsLimitedFileSize {
			fileInfo, err := d.FileLog.File.Stat()
			if err != nil {
				return fmt.Errorf("failed to get file info: %w", err)
			}

			newMsgSize := int64(len(msg))

			// If the log file exceeds the maximum size, trim it.
			if fileInfo.Size()+newMsgSize > d.FileLog.MaxFileSize {
				overSize := (fileInfo.Size() + newMsgSize) - d.FileLog.MaxFileSize
				if err := d.trimLogFile(overSize); err != nil {
					return fmt.Errorf("failed to trim log file: %w", err)
				}
			}
		}

		// Write the log message to the file.
		_, err := d.FileLog.File.WriteString(msg)
		return err
	}
	return fmt.Errorf("log file is not open") // Return an error if the log file is not open.
}

// trimLogFile trims the beginning of the log file by the specified size to fit the new log message.
func (d *LogRule) trimLogFile(overSize int64) error {
	// Open the existing log file for reading and writing.
	oldFile, err := os.OpenFile(d.FileLog.CurrentFileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open old log file: %w", err)
	}
	defer oldFile.Close() // Ensure the file is closed after this function returns.

	// Get the information about the old log file.
	oldFileInfo, err := oldFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	bytesToKeep := oldFileInfo.Size() - overSize // Calculate how many bytes to keep.

	if bytesToKeep <= 0 {
		// If there are no bytes to keep, truncate the file.
		if err := oldFile.Truncate(0); err != nil {
			return fmt.Errorf("failed to truncate log file: %w", err)
		}
		return nil
	}

	buffer := make([]byte, bytesToKeep) // Create a buffer for the remaining log data.

	// Read the remaining log data into the buffer.
	if _, err := oldFile.ReadAt(buffer, overSize); err != nil {
		return fmt.Errorf("failed to read remaining log data: %w", err)
	}

	// Write the remaining log data back to the start of the file.
	if _, err := oldFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek in log file: %w", err)
	}

	if _, err := oldFile.Write(buffer); err != nil {
		return fmt.Errorf("failed to write remaining log data: %w", err)
	}

	// Truncate the file to the new size.
	if err := oldFile.Truncate(bytesToKeep); err != nil {
		return fmt.Errorf("failed to truncate log file: %w", err)
	}

	return nil
}
