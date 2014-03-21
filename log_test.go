package log

import( 
	"testing"
)

func TestLevel(t *testing.T){
	log := NewLogger(CRITICAL)

	if log.level != CRITICAL{
		t.Error("Level not set for logger")
	}
}
func TestPanic(t *testing.T) {
	log := NewLogger(CRITICAL)

	log.Error("Shouldn't see this")

	defer func(){
		if r := recover(); r != nil {

		}else{
			t.Error("Expected a panic")
		}
	}()

	log.Panic("this is some text")
}

func TestStdLog(t *testing.T) {
	Critical("Fatal non-fatal message")
	Level(DEBUG)

	defer func(){
		if r := recover(); r != nil {

		}else{
			t.Error("Expected a panic")
		}
	}()
	Panic("Panic message")
}
