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
func Test_CalculateAllWindowCoords(t *testing.T) {
	windows := WindowsNewState(4, 2)
	windows = WinNew(windows, "Child", "0", "0", "1", "0", "C")
	// we have 2 windows here: "Terminal" (default) and "Child"

	compare_str_pair("CalcAll 1", windows["Child"][KeyXleftCalculated], "0", t)
	windows["Child"][KeyXshift] = StrMath(windows["Child"][KeyXshift], "+", "1")
	windows = CalculateAllWindowCoords(windows)
	compare_str_pair("CalcAll 1", windows["Child"][KeyXleftCalculated], "1", t)
	/*
		windows[winActiveId][KeyXshift] = StrMath(windows[winActiveId][KeyXshift], "+", "1")
		windows[winActiveId][KeyXshift] = StrMath(windows[winActiveId][KeyXshift], "-", "1")
		windows[winActiveId][KeyYshift] = StrMath(windows[winActiveId][KeyYshift], "+", "1")
		windows[winActiveId][KeyYshift] = StrMath(windows[winActiveId][KeyYshift], "-", "1")
	*/
}

func Test_StrMath(t *testing.T) {
	compare_str_pair("StrMath ", StrMath("1", "+", "2"), "3", t)
	compare_str_pair("StrMath ", StrMath("2", "-", "3"), "-1", t)
	compare_str_pair("StrMath ", StrMath("2", "*", "-3"), "-6", t)
	compare_str_pair("StrMath ", StrMath("6", "/", "-3"), "-2", t)

	// I don't want to stop at zero division so it's not an error for me
	compare_str_pair("StrMath ", StrMath("6", "/", "0"), "0", t)
	compare_str_pair("StrMath ", StrMath("0", "/", "6"), "0", t)
}

func Test_IsNumber(t *testing.T) {

	compare_bool_pair("IsNumber 1", IsNumber("-"), false, t)
	compare_bool_pair("IsNumber 2", IsNumber("+"), false, t)
	compare_bool_pair("IsNumber 3", IsNumber(""), false, t)
	compare_bool_pair("IsNumber 3", IsNumber(" "), false, t)
	compare_bool_pair("IsNumber 4", IsNumber(" +"), false, t)
	compare_bool_pair("IsNumber 5", IsNumber(" -"), false, t)
	compare_bool_pair("IsNumber 6", IsNumber(" ++"), false, t)
	compare_bool_pair("IsNumber 7", IsNumber(" --"), false, t)
	compare_bool_pair("IsNumber 8", IsNumber(" +-"), false, t)
	compare_bool_pair("IsNumber 9", IsNumber(" --"), false, t)

	received := IsNumber(" -1")
	compare_bool_pair("IsNumber 10", received, true, t)
	compare_bool_pair("IsNumber 11", IsNumber(" -1-"), false, t)
	compare_bool_pair("IsNumber 12", IsNumber("2a"), false, t)
}

func Test_new_window(t *testing.T) {
	windows := WindowsNewState(4, 2)
	// this windows fills the parent terminal,
	// so the right/bottom coords are equal with width/height
	winTerminalWidth := windows["Terminal"][KeyXright]
	winTerminalHeight := windows["Terminal"][KeyYbottom]
	compare_str_pair("new win 1", winTerminalWidth, "3", t)
	compare_str_pair("new win 2", winTerminalHeight, "1", t)

	winRenderedScreen := windows["Terminal"].RenderToScreenOfWin()
	wantedRendered := "" +
		"TTTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 3", winRenderedScreen.toString(), wantedRendered, t)

	////////////////////////// children ////////////////////////////////
	windows = WinNew(windows, "Child", "0", "0", "1", "0", "C")
	childRenderedScreen := windows["Child"].RenderToScreenOfWin()
	wantedChildRendered := "CC"
	compare_str_pair("new win 4", childRenderedScreen.toString(), wantedChildRendered, t)

	screenComposed := ScreensComposeToScreen(windows, []string{"Terminal", "Child"})
	wantedComposedRendered := "" +
		"CCTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 5", screenComposed.toString(), wantedComposedRendered, t)
}

func Test_ScreenNew(t *testing.T) {
	screen := ScreenEmpty(3, 2, "S", "Test_ScreenNew")
	compare_int_pair("ScreenNew 1", len(screen.pixels), 6, t) // 6 elems are in the screen
	compare_str_pair("ScreenNew 1", screen.pixels[coord{0, 0}], "S", t)

	wantedRendered := "" +
		"SSS" + NewLine() +
		"SSS"
	compare_str_pair("ScreenNew 2", screen.toString(), wantedRendered, t)

}
func Test_Empty(t *testing.T) {

}

// TEST FUNCTIONS ////////////////////////////////////////////////////////////
func compare_int_pair(source string, received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR INT (%v) - received: %v\n  wanted: %v",
			source, received, wanted)
	}
}

func compare_str_pair(source, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR STR (%v) - received: %v\n  wanted: %v",
			source, received, wanted)
	}
}
func compare_bool_pair(source string, received, wanted bool, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR BOOL (%v) - received: %v\n  wanted: %v",
			source, received, wanted)
	}
}
