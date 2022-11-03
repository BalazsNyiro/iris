package iris

import (
	"testing"
)

// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

/*
func Test_terminal_detect(t *testing.T) {
	widthStty, heightStty := TerminalDimensionsSttySize()
	fmt.Println("test stty:", widthStty, heightStty)
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("test syscall:", widthSys, heightSys)
}
func Test_document_create(t *testing.T) {
	rootObj := DocumentCreate("0", "50%", "50%", "40", "20")

}

*/
func Test_value_string_to_number(t *testing.T) {
	//rootObj := DocumentCreate("0", "50%", "50%", "60", "20")

	// no parent object usage, not a relative % value:
	value := valueStringToNumber("20", nil, "")
	compare_int_pair(value, 20, t)

	// % values need to know the parent measures to calculate relations
	rootObj := DocumentCreate("0", "60", "20", "", "0")
	valuePercent := valueStringToNumber("20%", &rootObj, "width")
	compare_int_pair(valuePercent, 12, t)

	rootObjHalfTerminal := DocumentCreate("0", "50%", "50%", "50", "20")
	valuePercentInHalfRoot := valueStringToNumber("20%", &rootObjHalfTerminal, "width")
	compare_int_pair(valuePercentInHalfRoot, 5, t)
}

func compare_int_pair(received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
	}
}
