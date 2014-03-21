package log

import (
	internal_logger "log"
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

type Logger interface {
	Critical(interface{})
	Error(interface{})
	Warn(interface{})
	Info(interface{})
	Debug(interface{})
	Trace(interface{})

	Level(int)
}

type Log struct {
	level int
}

func NewLogger(level int) *Log {
	return &Log{level: level}
}

var std = NewLogger(INFO)

func (log *Log) Panic(message interface{}) {
	log.Println(message)
	panic(message)
}

func Panic(message interface{}) {
	std.Panic(message)
}

func (log *Log) Critical(message interface{}) {
	log.Println(message)
}

func Critical(message interface{}) {
	std.Critical(message)
}

func (log *Log) Error(message interface{}) {
	if log.level < ERROR {
		log.Println(message)
	}
}

func Error(message interface{}) {
	std.Error(message)
}

func (log *Log) Warn(message interface{}) {
	if log.level < WARN {
		log.Println(message)
	}
}

func Warn(message interface{}) {
	std.Warn(message)
}

func (log *Log) Info(message interface{}) {
	if log.level < INFO {
		log.Println(message)
	}
}

func Info(message interface{}) {
	std.Info(message)
}

func (log *Log) Debug(message interface{}) {
	if log.level < DEBUG {
		log.Println(message)
	}
}

func Debug(message interface{}) {
	std.Debug(message)
}

func (log *Log) Trace(message interface{}) {
	if log.level < TRACE {
		log.Println(message)
	}
}

func Trace(message interface{}) {
	std.Trace(message)
}

func (log *Log) Println(message interface{}) {
	internal_logger.Println(message)
}

func Println(message interface{}) {
	std.Println(message)
}

func (log *Log) Level(level int) {
	log.level = level
}

func Level(level int) {
	std.Level(level)
}
