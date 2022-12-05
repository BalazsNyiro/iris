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
func Test_parameterCollect(t *testing.T) {
	tokens := []string{"1", "*", "2"}
	valueLeft, valueRight, idValueLeft, idValueRight, idError := ParametersCollect(tokens, 1)
	compare_str_pair("paramCollect1", valueLeft, "1", t)
	compare_str_pair("paramCollect1", valueRight, "2", t)
	compare_int_pair("paramCollect1", idValueLeft, 0, t)
	compare_int_pair("paramCollect1", idValueRight, 2, t)
	compare_str_pair("paramCollect1", idError, "", t)

	tokens = []string{"1", "*"} // missing right param
	valueLeft, valueRight, idValueLeft, idValueRight, idError = ParametersCollect(tokens, 1)
	compare_str_pair("paramCollect2", valueLeft, "", t)
	compare_str_pair("paramCollect2", valueRight, "", t)
	compare_bool_pair("paramCollect2", len(idError) > 0, true, t)

	tokens = []string{"*", "2"} // missing left param
	valueLeft, valueRight, idValueLeft, idValueRight, idError = ParametersCollect(tokens, 0)
	compare_str_pair("paramCollect2", valueLeft, "", t)
	compare_str_pair("paramCollect2", valueRight, "", t)
	compare_bool_pair("paramCollect2", len(idError) > 0, true, t)
}
func Test_ExprOperatorIsValid(t *testing.T) {
	compare_bool_pair("expr operator is valid 1a", ExprOperatorIsValid(""), false, t)
	compare_bool_pair("expr operator is valid 1a", ExprOperatorIsValid("?"), false, t)
	compare_bool_pair("expr operator is valid 1b", ExprOperatorIsValid("++"), false, t)
	compare_bool_pair("expr operator is valid 2", ExprOperatorIsValid("+"), true, t)
	compare_bool_pair("expr operator is valid 2", ExprOperatorIsValid("-"), true, t)
	compare_bool_pair("expr operator is valid 2", ExprOperatorIsValid("*"), true, t)
	compare_bool_pair("expr operator is valid 2", ExprOperatorIsValid("/"), true, t)
}

func Test_CleanStringLists(t *testing.T) {
	elems := []string{"a", "b", "", "c", ""}
	cleaned := StrListRemoveEmptyElems(elems, true)
	compare_str_lists("cleanStringList1", cleaned, []string{"a", "b", "c"}, t)
}

func Test_DoubleSpaceRemove(t *testing.T) {
	compare_str_pair(
		"double space remove",
		StrDoubleSpacesRemove("a   b     c"),
		"a b c",
		t)
}

func Test_OperatorPrecedence(t *testing.T) {
	operatorNextId, operatorNext := TokenOperatorNext([]string{"1", "+", "2"})
	compare_int_pair("operatorPrecedence1i", operatorNextId, 1, t)
	compare_str_pair("operatorPrecedence1s", operatorNext, "+", t)

	operatorNextId, operatorNext = TokenOperatorNext([]string{"1", "+", "2", "*", "3"})
	compare_int_pair("operatorPrecedence2i", operatorNextId, 3, t)
	compare_str_pair("operatorPrecedence2s", operatorNext, "*", t)

	operatorNextId, operatorNext = TokenOperatorNext([]string{"1", "+", "5", "-", "3", "/", "3"})
	compare_int_pair("operatorPrecedence3i", operatorNextId, 5, t)
	compare_str_pair("operatorPrecedence3s", operatorNext, "/", t)
}
func Test_CoordExpressionEval(t *testing.T) {
	windows := WindowsNewState(4, 2)

	expression := "1"
	result := CoordExpressionEval(expression, windows)
	compare_str_pair("eval1", result, "1", t)

	expression = "win:Terminal:" + KeyXrightCalculated
	result = CoordExpressionEval(expression, windows)
	compare_str_pair("eval2", result, "3", t)

	expression = "1 + 2"
	result = CoordExpressionEval(expression, windows)
	compare_str_pair("eval3", result, "3", t)

	expression = "win:Terminal:" + KeyXrightCalculated + " - 2 * 2"
	result = CoordExpressionEval(expression, windows)
	compare_str_pair("eval4", result, "-1", t)
}

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

	winRenderedScreen := windows["Terminal"].RenderToScreenOfWin("debug")
	wantedRendered := "" +
		"TTTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 3", winRenderedScreen.toString(), wantedRendered, t)

	////////////////////////// children ////////////////////////////////
	windows = WinNew(windows, "Child", "0", "0", "1", "0", "C")
	childRenderedScreen := windows["Child"].RenderToScreenOfWin("debug")
	wantedChildRendered := "CC"
	compare_str_pair("new win 4", childRenderedScreen.toString(), wantedChildRendered, t)

	screenComposed := ScreensComposeToScreen(windows, []string{"Terminal", "Child"}, "debug")
	wantedComposedRendered := "" +
		"CCTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 5", screenComposed.toString(), wantedComposedRendered, t)
}

func Test_ScreenNew(t *testing.T) {
	screen := ScreenEmpty(3, 2, "S", "Test_ScreenNew")
	compare_int_pair("ScreenNew 1", len(screen.matrixCharsRendered), 6, t) // 6 elems are in the screen
	compare_str_pair("ScreenNew 1", screen.matrixCharsRendered[Coord{0, 0}], "S", t)

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

func compare_str_lists(source string, received, wanted []string, t *testing.T) {
	if len(received) != len(wanted) {
		t.Fatalf("\nERROR STR LIST, different lengths - (%v) - received: %v\n  wanted: %v", source, received, wanted)
	}
	for id, _ := range received {
		compare_str_pair(source, received[id], wanted[id], t)
	}
}

func compare_bool_pair(source string, received, wanted bool, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR BOOL (%v) - received: %v\n  wanted: %v",
			source, received, wanted)
	}
}
