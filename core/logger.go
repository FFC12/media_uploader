package core

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

// initializeLogger initializes the global logger with desired configurations.
func InitializeLogger() {
	if log != nil {
		return
	}

	log = logrus.New()

	// Set log level to Info by default
	log.SetLevel(logrus.DebugLevel)

	// Set the log format to JSON
	log.SetFormatter(&logrus.JSONFormatter{})

	// Set the output to both console and file with log rotation
	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of backups
		MaxAge:     1,  // days
		Compress:   true,
	})
}

// LogError logs an error message
// LogError logs an error message with file, function, and line information
func LogError(message string, err error) {
	log.WithFields(logrus.Fields{
		"file":    fileInfo(),
		"caller":  callerInfo(),
		"error":   err,
		"message": message,
	}).Error("Error occurred")
}

// LogWarning logs a warning message with file, function, and line information
func LogWarning(message string) {
	log.WithFields(logrus.Fields{
		"file":    fileInfo(),
		"caller":  callerInfo(),
		"message": message,
	}).Warn("Warning occurred")
}

// LogInfo logs an info message with file, function, and line information
func LogInfo(message string) {
	log.WithFields(logrus.Fields{
		"file":    fileInfo(),
		"caller":  callerInfo(),
		"message": message,
	}).Info("Info message")
} 

func LogFatal(message string, err error) {
	log.WithFields(logrus.Fields{
		"file":    fileInfo(),
		"caller":  callerInfo(),
		"error":   err,
		"message": message,
	}).Fatal("Fatal error occurred")
}

// fileInfo returns file information for logging
func fileInfo() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// callerInfo returns the name of the calling function for logging
func callerInfo() string {
	pc, _, _, _ := runtime.Caller(2)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}

// closeLogFile closes the log file
func closeLogFile() {
	// Since lumberjack.Logger handles closing the file, there's no need to explicitly close it here
}
