package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLevel(t *testing.T) {
	log := NewLogger(CRITICAL)

	if log.level != CRITICAL {
		t.Error("Level not set for logger")
	}
}
func TestPanic(t *testing.T) {
	log := NewLogger(CRITICAL)

	log.Error("Shouldn't see this")

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("Expected a panic")
		}
	}()

	log.Panic("this is some text")
}

func TestStdLog(t *testing.T) {
	Critical("Fatal non-fatal message")
	Level(DEBUG)

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("Expected a panic")
		}
	}()
	Panic("Panic message")
}

func TestSetOutput(t *testing.T) {
	os.MkdirAll("tmp", 0777)
	path := "./tmp/nonexistantfile.txt"
	err := SetOutput(path)

	if err != nil {
		t.Error("There was an error setting output: " + err.Error())
	}

	logMsg := "This message should appear in the logs"
	Info(logMsg)

	logFile, _ := ioutil.ReadFile(path)
	logContents := string(logFile)

	if !strings.Contains(logContents, logMsg) {
		t.Error("Log didn't contain what we expected: " + logContents)
	}

	os.RemoveAll("tmp")
}
