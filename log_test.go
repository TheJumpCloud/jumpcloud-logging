package log

import (
	"fmt"
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

	// Close the output file so we can remove the temp dir
	CloseOutput()

	os.RemoveAll("tmp")
}

func TestErrorWithPrintf(t *testing.T) {
	const (
		testStr  = "test-string"
		testStr2 = "another test string"
		testStr3 = "test789"
		testStr4 = "test233"
	)

	os.MkdirAll("tmp", 0777)
	path := "./tmp/nonexistentfile.txt"
	err := SetOutput(path)
	defer func() {
		CloseOutput()
		os.Remove(path)
	}()

	Error("my output should contain '%s'", testStr)
	Error(testStr2)
	Error(50)
	Error("another", testStr3, testStr4)

	myData, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Could not read log file '%s' for test, err='%s'", path, err)
	}

	data := string(myData)

	t.Logf("data='%s'", data)

	if strings.Index(data, testStr) == -1 {
		t.Fatalf("Could not find the expected test constant '%s' in the log file, log contained '%s'", testStr, string(myData))
	}

	if strings.Index(data, testStr2) == -1 {
		t.Fatalf("Could not find the expected test constant '%s' in the log file, log contained '%s'", testStr2, string(myData))
	}

	if strings.Index(data, "Unsupported argument type") == -1 {
		t.Fatalf("Could not find output for incorrect arg type in the log file, log contained '%s'", string(myData))
	}

	if strings.Index(data, testStr4) == -1 {
		t.Fatalf("Could not find output for lists of strings '%s', log contained '%s'", testStr4, string(myData))
	}
}

const (
	TEST_LOG_SIZE = 10000
)

func TestRotation(t *testing.T) {
	os.MkdirAll("tmp", 0777)
	path := "./tmp/nonexistentfile.txt"
	err := SetOutput(path)

	if err != nil {
		t.Error("There was an error setting output: " + err.Error())
	}

	SetMaxLogSize(TEST_LOG_SIZE)

	for i := 0; i < 250; i++ {
		Info(fmt.Sprintf("%03d-01234567890123456789012345678901234567890123456789", i))
	}

	//
	// Rotate the log file if it grows too large
	//
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Could not stat log file '%s', err='%s'", path, err)
	}

	if fileInfo.Size() > TEST_LOG_SIZE {
		t.Fatalf("Log file is %d bytes, larger than %d after rotation.", fileInfo.Size(), TEST_LOG_SIZE)
	}

	rotatedFile := path + ".prev"

	fileInfo, err = os.Stat(rotatedFile)
	if err != nil {
		t.Fatalf("Could not stat rotated log file '%s', err='%s'", rotatedFile, err)
	}

	if fileInfo.Size() > TEST_LOG_SIZE+100 {
		t.Fatalf("Rotated log file is bigger than expected (%d bytes), should be no more than %d", fileInfo.Size(), TEST_LOG_SIZE+100)
	}

	// Close the output file so we can remove the temp dir
	CloseOutput()

	os.RemoveAll("tmp")
}
