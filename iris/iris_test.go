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

func Test_new_window(t *testing.T) {
	windowsState := WindowsNewState(4, 2)
	winTerminalWidth := windowsState["Terminal"][KeyWidth]
	// compare_str_pair(winTerminalWidth, "5", t)
	_ = winTerminalWidth
}

func Test_ScreenNew(t *testing.T) {
	screen := ScreenEmpty(3, 2)
	compare_int_pair(len(screen), 3, t)    // it has 3 columns
	compare_int_pair(len(screen[0]), 2, t) // it has 2 rows in a column
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
