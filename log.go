package log

import (
	"fmt"
	internal_logger "log"
	"os"
	"sync"
)

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
	level            int
	maxLogSize       int64
	logFileName      string
	logWriter        *os.File
	mu               sync.RWMutex
	currentByteCount int64
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
	outFmt := fmt.Sprintf("[%s] %s", "PANIC", format)
	log.Println(outFmt, message...)

	panic(fmt.Sprintf(outFmt, message...))
}

func Panic(format string, message ...interface{}) {
	std.Panic(format, message...)
}

func (log *Log) Critical(format string, message ...interface{}) {
	outFmt := fmt.Sprintf("[%s] %s", "CRITICAL", format)
	log.Println(outFmt, message...)
}

func Critical(format string, message ...interface{}) {
	std.Critical(format, message...)
}

func (log *Log) Error(format string, message ...interface{}) {
	if log.level <= ERROR {
		outFmt := fmt.Sprintf("[%s] %s", "ERROR", format)
		log.Println(outFmt, message...)
	}
}

func Error(format string, message ...interface{}) {
	std.Error(format, message...)
}

func (log *Log) Warn(format string, message ...interface{}) {
	if log.level <= WARN {
		outFmt := fmt.Sprintf("[%s] %s", "WARN", format)
		log.Println(outFmt, message...)
	}
}

func Warn(format string, message ...interface{}) {
	std.Warn(format, message...)
}

func (log *Log) Info(format string, message ...interface{}) {
	if log.getLevel() <= INFO {
		outFmt := fmt.Sprintf("[%s] %s", "INFO", format)
		log.Println(outFmt, message...)
	}
}

func Info(format string, message ...interface{}) {
	std.Info(format, message...)
}

func (log *Log) Debug(format string, message ...interface{}) {
	if log.getLevel() <= DEBUG {
		outFmt := fmt.Sprintf("[%s] %s", "DEBUG", format)
		log.Println(outFmt, message...)
	}
}

func Debug(format string, message ...interface{}) {
	std.Debug(format, message...)
}

func (log *Log) Trace(format string, message ...interface{}) {
	if log.getLevel() <= TRACE {
		outFmt := fmt.Sprintf("[%s] %s", "TRACE", format)
		log.Println(outFmt, message...)
	}
}

func Trace(format string, message ...interface{}) {
	std.Trace(format, message...)
}

func (log *Log) Println(format string, message ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()

	outFmt := fmt.Sprintf("[%d] %s", os.Getpid(), format)
	outString := fmt.Sprintf(outFmt, message...)
	internal_logger.Print(outString)

	log.currentByteCount += int64(len(outString) + 21) //Adding the date, newline.
	if log.needsRotating() {
		log.rotateLog()
		internal_logger.Print("Log rotated")
		log.currentByteCount += 32
	}
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

func (log *Log) getLevel() int {
	log.mu.RLock()
	defer log.mu.RUnlock()

	return log.level
}

func (log *Log) SetOutput(path string) (err error) {
	log.mu.Lock()
	defer log.mu.Unlock()

	//Open the file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		err = fmt.Errorf("Unable to open log file - err='%s'", err.Error())
		return
	}

	log.logFileName = path
	log.logWriter = file
	internal_logger.SetOutput(file)

	//Reset the line count as appropriate
	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo != nil {
		log.currentByteCount = fileInfo.Size()
	} else {
		log.currentByteCount = 0
	}

	return
}

func SetOutput(path string) (err error) {
	return std.SetOutput(path)
}

func (log *Log) CloseOutput() (err error) {
	log.mu.Lock()
	defer log.mu.Unlock()

	return log.logWriter.Close()
}

func CloseOutput() (err error) {
	return std.CloseOutput()
}

func (log *Log) needsRotating() bool {
	return log.logFileName != "" && log.currentByteCount > log.maxLogSize
}

func (log *Log) rotateLog() {
	// Close our current file
	log.logWriter.Close()

	// Delete any existing rotated file
	rotateName := log.logFileName + ".prev"
	os.Remove(rotateName)

	// Rename the exiting file
	os.Rename(log.logFileName, rotateName)

	// Re-open the target file with truncation
	file, err := os.OpenFile(log.logFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err == nil {
		log.logWriter = file
		internal_logger.SetOutput(file)
		log.currentByteCount = 0
	} else {
		fmt.Printf("ERROR: Unable to open log file with truncation for rotation. err='%s'\n", err.Error())
	}
}
