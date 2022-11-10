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
func Test_render(t *testing.T) {
}

func Test_value_string_to_number(t *testing.T) {

}

func Test_Empty(t *testing.T) {

}

// TEST FUNCTIONS ////////////////////////////////////////////////////////////
func compare_int_pair(received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR INT - received: %v\n  wanted: %v", received, wanted)
	}
}

func compare_str_pair(received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR STR - received: %v\n  wanted: %v", received, wanted)
	}
}
