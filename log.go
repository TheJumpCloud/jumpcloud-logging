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

func (log *Log) Panic(format string, message ...interface{}) {
	log.rotateLog()
	log.Println(format, message...)

	panic(fmt.Sprintf(format, message...))
}

func Panic(format string, message ...interface{}) {
	std.Panic(format, message...)
}

func (log *Log) Critical(format string, message ...interface{}) {
	log.rotateLog()

	log.Println(format, message...)
}

func Critical(format string, message ...interface{}) {
	std.Critical(format, message...)
}

func (log *Log) Error(format string, message ...interface{}) {
	log.rotateLog()

	if log.level <= ERROR {
		log.Println(format, message...)
	}
}

func Error(format string, message ...interface{}) {
	std.Error(format, message...)
}

func (log *Log) Warn(format string, message ...interface{}) {
	log.rotateLog()

	if log.level <= WARN {
		log.Println(format, message...)
	}
}

func Warn(format string, message ...interface{}) {
	std.Warn(format, message...)
}

func (log *Log) Info(format string, message ...interface{}) {
	log.rotateLog()

	if log.getLevel() <= INFO {
		log.Println(format, message...)
	}
}

func Info(format string, message ...interface{}) {
	std.Info(format, message...)
}

func (log *Log) Debug(format string, message ...interface{}) {
	log.rotateLog()

	if log.getLevel() <= DEBUG {
		log.Println(format, message...)
	}
}

func Debug(format string, message ...interface{}) {
	std.Debug(format, message...)
}

func (log *Log) Trace(format string, message ...interface{}) {
	log.rotateLog()

	if log.getLevel() <= TRACE {
		log.Println(format, message...)
	}
}

func Trace(format string, message ...interface{}) {
	std.Trace(format, message...)
}

func (log *Log) Println(format string, message ...interface{}) {
	outFmt := fmt.Sprintf("[%d] %s", os.Getpid(), format)

	internal_logger.Printf(outFmt, message...)
}

func Println(format string, message ...interface{}) {
	std.Println(format, message...)
}

func (log *Log) Level(level int) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.level = level
}

func Level(level int) {
	std.Level(level)
}

// thread-safe, use for accessing level in exported methods
func (log *Log) getLevel() int {
	log.mu.RLock()
	defer log.mu.RUnlock()

	return log.level
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
			//fmt.Println("Rotated log successfully!")
		} else {
			fmt.Printf("ERROR: Unable to open log file with truncation for rotation. err='%s'\n", err.Error())
		}
	}
}
