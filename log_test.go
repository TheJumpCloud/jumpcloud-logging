package log

import( 
	"testing"
)

func TestLevel(t *testing.T){
	log := NewLogger(FATAL)

	log.Fatal("this is some text")
}

func TestPanic(t *testing.T) {
	log := NewLogger(FATAL)
	defer func(){
		if r := recover(); r != nil {

		}else{
			t.Error("Expected a panic")
		}
	}()

	log.Panic("this is some text")
}

func TestStdLog(t *testing.T) {
	Fatal("Fatal non-fatal message")
	Level(DEBUG)

	defer func(){
		if r := recover(); r != nil {

		}else{
			t.Error("Expected a panic")
		}
	}()
	Panic("Panic message")
}
