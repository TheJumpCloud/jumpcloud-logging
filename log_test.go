package log

import( 
	"testing"
)

func TestLevel(t *testing.T){
	log := NewLogger(FATAL)

	log.Fatal("this is some text")
}
