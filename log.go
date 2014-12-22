package log

import (
	"fmt"
	"io"
	internal_logger "log"
	"os"
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
	log.maxLogSize = logSize
}

func SetMaxLogSize(logSize int64) {
	std.maxLogSize = logSize
}

func (log *Log) Panic(message ...interface{}) {
	log.Println(message)
	panic(message)
}

func Panic(message ...interface{}) {
	rotateLog()
	std.Panic(message)
}

func (log *Log) Critical(message ...interface{}) {
	log.Println(message)
}

func Critical(message ...interface{}) {
	rotateLog()
	std.Critical(message)
}

func (log *Log) Error(message ...interface{}) {
	if log.level <= ERROR {
		log.Println(message)
	}
}

func Error(message ...interface{}) {
	rotateLog()
	std.Error(message)
}

func (log *Log) Warn(message ...interface{}) {
	if log.level <= WARN {
		log.Println(message)
	}
}

func Warn(message ...interface{}) {
	rotateLog()
	std.Warn(message)
}

func (log *Log) Info(message ...interface{}) {
	if log.level <= INFO {
		log.Println(message)
	}
}

func Info(message ...interface{}) {
	rotateLog()
	std.Info(message)
}

func (log *Log) Debug(message ...interface{}) {
	if log.level <= DEBUG {
		log.Println(message)
	}
}

func Debug(message ...interface{}) {
	rotateLog()
	std.Debug(message)
}

func (log *Log) Trace(message ...interface{}) {
	if log.level <= TRACE {
		log.Println(message)
	}
}

func Trace(message ...interface{}) {
	rotateLog()
	std.Trace(message)
}

func (log *Log) Println(message ...interface{}) {
	internal_logger.Println(message)
}

func Println(message ...interface{}) {
	rotateLog()
	std.Println(message)
}

func (log *Log) Level(level int) {
	log.level = level
}

func Level(level int) {
	std.Level(level)
}

func (log *Log) SetOutput(writer io.Writer) {
	internal_logger.SetOutput(writer)
}

func SetOutput(path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		err = fmt.Errorf("Unable to open log file - err='%s'", err.Error())
		return
	}

	std.logFileName = path
	std.logWriter = file

	std.SetOutput(file)

	return
}

func CloseOutput() (err error) {
	err = std.logWriter.Close()

	return
}

func rotateLog() {

	// No need to rotate a log file that doesn't exist...
	if std.logFileName == "" {
		return
	}

	fileInfo, err := os.Stat(std.logFileName)
	if err != nil {
		fmt.Printf("ERROR: Could not stat log file '%s' - err='%s'\n", std.logFileName, err.Error())

		// Attempt recovery to prevent filling the filesystem...
		CloseOutput()
		os.Remove(std.logFileName)
		SetOutput(std.logFileName)
		return
	}

	// Rotate the log file if it grows too large
	if fileInfo.Size() > std.maxLogSize {

		// Close the current file so we can rename it...
		CloseOutput()

		// If the prev file already exists, just blast it...
		os.Remove(std.logFileName + ".prev")

		err := os.Rename(std.logFileName, std.logFileName+".prev")
		if err != nil {
			fmt.Printf("ERROR: Could not rename log file for rotation, attempting to remove log file instead. err='%s'\n", err.Error())

			os.Remove(std.logFileName)
		}

		// Re-open the file...
		SetOutput(std.logFileName)
	}
}
