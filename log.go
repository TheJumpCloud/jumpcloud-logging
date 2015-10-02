package log

import (
	"fmt"
	"io"
	internal_logger "log"
	"os"
	"sync"
)

//This is a wrapper only we can switch out loggers at will. The golang logger ecosphere is still volatile

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
)

const (
	MAX_DEFAULT_LOG_SIZE int64 = 2000000
)

type Logger interface {
	Critical(interface{})
	Error(interface{})
	Warn(interface{})
	Info(interface{})
	Debug(interface{})
	Trace(interface{})
	SetOutput(string)

	Level(int)
}

type Log struct {
	level       int
	maxLogSize  int64
	logFileName string
	logWriter   *os.File
	mu          sync.RWMutex
}

func NewLogger(level int) *Log {
	return &Log{
		level:       level,
		maxLogSize:  MAX_DEFAULT_LOG_SIZE,
		logFileName: "",
	}
}

var std = NewLogger(INFO)

func (log *Log) SetMaxLogSize(logSize int64) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.maxLogSize = logSize
}

func SetMaxLogSize(logSize int64) {
	std.maxLogSize = logSize
}

func (log *Log) Panic(message ...interface{}) {
	log.rotateLog()

	log.Println(message)
	panic(message)
}

func Panic(message ...interface{}) {
	std.Panic(message)
}

func (log *Log) Critical(message ...interface{}) {
	log.rotateLog()

	log.Println(message)
}

func Critical(message ...interface{}) {
	std.Critical(message)
}

func (log *Log) Error(message ...interface{}) {
	log.rotateLog()

	log.mu.RLock()
	level := log.level
	log.mu.RUnlock()

	if level <= ERROR {
		log.Println(message)
	}
}

func Error(message ...interface{}) {
	std.Error(message)
}

func (log *Log) Warn(message ...interface{}) {
	log.rotateLog()

	log.mu.RLock()
	level := log.level
	log.mu.RUnlock()

	if level <= WARN {
		log.Println(message)
	}
}

func Warn(message ...interface{}) {
	std.Warn(message)
}

func (log *Log) Info(message ...interface{}) {
	log.rotateLog()

	log.mu.RLock()
	level := log.level
	log.mu.RUnlock()

	if level <= INFO {
		log.Println(message)
	}
}

func Info(message ...interface{}) {
	std.Info(message)
}

func (log *Log) Debug(message ...interface{}) {
	log.rotateLog()

	log.mu.RLock()
	level := log.level
	log.mu.RUnlock()

	if level <= DEBUG {
		log.Println(message)
	}
}

func Debug(message ...interface{}) {
	std.Debug(message)
}

func (log *Log) Trace(message ...interface{}) {
	log.rotateLog()

	log.mu.RLock()
	level := log.level
	log.mu.RUnlock()

	if level <= TRACE {
		log.Println(message)
	}
}

func Trace(message ...interface{}) {
	std.Trace(message)
}

func (log *Log) Println(message ...interface{}) {
	log.rotateLog()

	internal_logger.Println(message)
}

func Println(message ...interface{}) {
	std.Println(message)
}

func (log *Log) Level(level int) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.level = level
}

func Level(level int) {
	std.Level(level)
}

func (log *Log) SetOutput(path string) (err error) {
	log.mu.Lock()
	defer log.mu.Unlock()

	return log.setOutput(path)
}

// Not thread-safe
func (log *Log) setOutput(path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		err = fmt.Errorf("Unable to open log file - err='%s'", err.Error())
		return
	}

	log.logFileName = path
	log.logWriter = file
	internal_logger.SetOutput(file)
	return
}

func SetOutput(path string) (err error) {
	return std.SetOutput(path)
}

func (log *Log) CloseOutput() (err error) {
	log.mu.Lock()
	defer log.mu.Unlock()

	return log.closeOutput()
}

// Not thread-safe
func (log *Log) closeOutput() (err error) {
	return log.logWriter.Close()
}

func CloseOutput() (err error) {
	return std.CloseOutput()
}

func (log *Log) rotateLog() {
	// No need to rotate a log file that doesn't exist...
	if log.logFileName == "" {
		return
	}

	// Acquire mutex to rotate
	log.mu.Lock()
	defer log.mu.Unlock()

	fileInfo, err := os.Stat(log.logFileName)
	if err != nil {
		fmt.Printf("ERROR: Could not stat log file '%s' - err='%s'\n", log.logFileName, err.Error())

		// Attempt recovery to prevent filling the filesystem...
		log.closeOutput()
		os.Remove(log.logFileName)
		log.setOutput(log.logFileName)
		return
	}

	// Rotate the log file if it grows too large
	if fileInfo.Size() > log.maxLogSize {
		// Close our current file
		log.logWriter.Close()

		// Delete any existing rotated file
		rotateName := log.logFileName + ".prev"
		os.Remove(rotateName)

		// Copy the existing file over so we can truncate in place
		dst, err := os.OpenFile(rotateName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Printf("ERROR: Could not open target file for rotation. err='%s'\n", err.Error())
		} else {
			src, err := os.OpenFile(log.logFileName, os.O_RDONLY, 0600)
			if err != nil {
				fmt.Printf("ERROR: Could not open source file for rotation. err='%s'\n", err.Error())
			} else {
				io.Copy(dst, src)
			}
		}

		// Re-open the target file with truncation
		file, err := os.OpenFile(log.logFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err == nil {
			log.logWriter = file
			internal_logger.SetOutput(file)
			fmt.Println("Rotated log successfully!")
		} else {
			fmt.Printf("ERROR: Unable to open log file with truncation for rotation. err='%s'\n", err.Error())
		}
	}
}
