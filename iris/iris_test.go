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
	windows := WindowsNewState(4, 2)
	// this windows fills the parent terminal,
	// so the right/bottom coords are equal with width/height
	winTerminalWidth := windows["Terminal"][KeyXright]
	winTerminalHeight := windows["Terminal"][KeyYbottom]
	compare_str_pair(winTerminalWidth, "3", t)
	compare_str_pair(winTerminalHeight, "1", t)
	_ = winTerminalWidth
}

func Test_ScreenNew(t *testing.T) {
	screen := ScreenEmpty(3, 2, "S", "Test_ScreenNew")
	compare_int_pair(len(screen.pixels), 6, t) // 6 elems are in the screen
	compare_str_pair(screen.pixels[coord{0, 0}], "S", t)

	wantedRendered := "" +
		"SSS" + NewLine() +
		"SSS"
	compare_str_pair(screen.toString(), wantedRendered, t)

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
