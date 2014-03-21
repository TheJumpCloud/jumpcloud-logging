package log

import (
	"github.com/jcelliott/lumber"
)

//This is a wrapper only we can switch out loggers at will. The golang logger ecosphere is still volatile

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	Fatal(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
	Trace(string)

	Level(int)
}

type Log struct{
	internal_logger lumber.Logger	
}

func NewLogger(o int) *Log {
	return &Log{internal_logger:lumber.NewConsoleLogger(o)}
}

func (log *Log) Fatal(message string){
	log.internal_logger.Fatal(message)
}

func (log *Log) Error(message string){
	log.internal_logger.Error(message)
}

func (log *Log) Warn(message string){
	log.internal_logger.Warn(message)
}

func (log *Log) Info(message string){
	log.internal_logger.Info(message)
}

func (log *Log) Debug(message string){
	log.internal_logger.Debug(message)
}

func (log *Log) Trace(message string){
	log.internal_logger.Trace(message)
}
