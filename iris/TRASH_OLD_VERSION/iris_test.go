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

func Test_WinNamesSortByAttribute(t *testing.T) {
	terminalWidth := 8
	terminalHeight := 4
	windows, _ := WinNewStateIntoWindows(terminalWidth, terminalHeight)

	windows = WinCreateIntoWindows(windows, "nameC", "1", "6", Int2Str(1), Int2Str(1), "A")
	windows = WinCreateIntoWindows(windows, "nameD", "2", "5", Int2Str(2), Int2Str(2), "D")
	windows = WinCreateIntoWindows(windows, "nameB", "3", "4", Int2Str(5), Int2Str(2), "C")
	windows = WinCreateIntoWindows(windows, "nameA", "4", "3", Int2Str(3), Int2Str(2), "B")

	sortedByNames := WinNamesSortByAttribute(windows,
		[]string{"nameC", "nameD", "nameB", "nameA"},
		KeyWinName, "string")
	compare_str_lists("sortByNames", sortedByNames, []string{"nameA", "nameB", "nameC", "nameD"}, t)

	sortedByNames = WinNamesSortByAttribute(windows,
		[]string{"nameC", "nameD", "nameB", "nameA"},
		KeyDebugWindowFillerChar, "string")
	compare_str_lists("sortByFillChar", sortedByNames,
		[]string{"nameC", "nameA", "nameB", "nameD"}, t)

	sortedByNames = WinNamesSortByAttribute(windows, // 1, 2, 3, 4
		[]string{"nameC", "nameD", "nameB", "nameA"},
		KeyXleft, "number")
	compare_str_lists("sortByFillChar", sortedByNames,
		[]string{"nameC", "nameD", "nameB", "nameA"}, t)

	sortedByNames = WinNamesSortByAttribute(windows, // 3, 4, 5, 6
		[]string{"nameC", "nameD", "nameB", "nameA"},
		KeyYtop, "number")
	compare_str_lists("sortByFillChar", sortedByNames,
		[]string{"nameA", "nameB", "nameD", "nameC"}, t)
}

func Test_MatrixCharsInsertContentIntoOneWindow(t *testing.T) {
	terminalWidth := 8
	terminalHeight := 4

	winName := "TestWindows"
	winTestWidth := 5
	winTestHeight := 2

	windows, windowsChars := WinNewStateIntoWindows(terminalWidth, terminalHeight)
	windows = WinCreateIntoWindows(windows, winName, "1", "2", Int2Str(winTestWidth), Int2Str(winTestHeight), "T")
	windowsChars = WinTextUpdate(windowsChars, winName, "simpleText", "Dogs are BARKING...")
	windows["prgState"]["winActiveId"] = winName

	matrixChars := MatrixCharsEmptyOfWindows(winTestWidth, winTestHeight, 'M', "winName:"+winName)
	matrixChars = MatrixCharsInsertContentOfWindows(matrixChars, winTestWidth, winTestHeight, windowsChars[winName], true)
	compare_rune_pair("MatrixCharsInsertContentA", matrixChars.Rendered[MatrixCoord{0, 0}].CharVal, 'D', t)
	compare_rune_pair("MatrixCharsInsertContentB", matrixChars.Rendered[MatrixCoord{1, 0}].CharVal, 'o', t)
	compare_rune_pair("MatrixCharsInsertContentC", matrixChars.Rendered[MatrixCoord{2, 0}].CharVal, 'g', t)
	compare_rune_pair("MatrixCharsInsertContentD", matrixChars.Rendered[MatrixCoord{3, 0}].CharVal, 's', t)
	compare_rune_pair("MatrixCharsInsertContentE", matrixChars.Rendered[MatrixCoord{4, 0}].CharVal, ' ', t)
	compare_rune_pair("MatrixCharsInsertContentF", matrixChars.Rendered[MatrixCoord{0, 1}].CharVal, 'a', t)
	compare_rune_pair("MatrixCharsInsertContentG", matrixChars.Rendered[MatrixCoord{1, 1}].CharVal, 'r', t)
	compare_rune_pair("MatrixCharsInsertContentH", matrixChars.Rendered[MatrixCoord{2, 1}].CharVal, 'e', t)
}

func Test_parameterCollect(t *testing.T) {
	tokens := []string{"1", "*", "2"}
	valueLeft, valueRight, idValueLeft, idValueRight, idError := TokenParametersCollect(tokens, 1)
	compare_str_pair("paramCollect1", valueLeft, "1", t)
	compare_str_pair("paramCollect1", valueRight, "2", t)
	compare_int_pair("paramCollect1", idValueLeft, 0, t)
	compare_int_pair("paramCollect1", idValueRight, 2, t)
	compare_str_pair("paramCollect1", idError, "", t)

	tokens = []string{"1", "*"} // missing right param
	valueLeft, valueRight, idValueLeft, idValueRight, idError = TokenParametersCollect(tokens, 1)
	compare_str_pair("paramCollect2", valueLeft, "", t)
	compare_str_pair("paramCollect2", valueRight, "", t)
	compare_bool_pair("paramCollect2", len(idError) > 0, true, t)

	tokens = []string{"*", "2"} // missing left param
	valueLeft, valueRight, idValueLeft, idValueRight, idError = TokenParametersCollect(tokens, 0)
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
	windows, _ := WinNewStateIntoWindows(4, 2)

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
	windows, _ := WinNewStateIntoWindows(4, 2)
	windows = WinCreateIntoWindows(windows, "Child", "0", "0", "1", "0", "C")
	// we have 2 windows here: "Terminal" (default) and "Child"

	compare_str_pair("CalcAll 1", windows["Child"][KeyXleftCalculated], "0", t)
	windows["Child"][KeyXshift] = StrMath(windows["Child"][KeyXshift], "+", "1")
	windows = WinCoordsCalculateUpdate(windows)
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

func Test_RenderToMatrixOfWin(t *testing.T) {
	windows, windowsChars := WinNewStateIntoWindows(5, 5)
	windows = WinCreateIntoWindows(windows, "Child", "1", "1", "3", "3", "C")
	windowsChars = WinTextUpdate(windowsChars, "Child", "simpleText", "abcdefghijklmnopq")       // the input test is long but the displayed
	childRenderedMatrixChars := windows["Child"].RenderToMatrixCharsOfWin(windowsChars, "debug") // area is only 3x3 char
	wantedChildRendered := "" +
		"abc" + NewLine() +
		"def" + NewLine() +
		"ghi"
	compare_str_pair("RenderToMatrixCharsOfWin", childRenderedMatrixChars.toString(), wantedChildRendered, t)
}

func Test_new_window(t *testing.T) {
	windows, windowsChars := WinNewStateIntoWindows(4, 2)
	// this windows fills the parent terminal,
	// so the right/bottom coords are equal with width/height
	winTerminalWidth := windows["Terminal"][KeyXright]
	winTerminalHeight := windows["Terminal"][KeyYbottom]
	compare_str_pair("new win 1", winTerminalWidth, "3", t)
	compare_str_pair("new win 2", winTerminalHeight, "1", t)

	winRenderedMatrixChars := windows["Terminal"].RenderToMatrixCharsOfWin(windowsChars, "debug")
	wantedRendered := "" +
		"TTTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 3", winRenderedMatrixChars.toString(), wantedRendered, t)

	////////////////////////// children ////////////////////////////////
	windows = WinCreateIntoWindows(windows, "Child", "2", "1", "3", "1", "C")
	childRenderedMatrixChars := windows["Child"].RenderToMatrixCharsOfWin(windowsChars, "debug")
	wantedChildRendered := "CC"
	compare_str_pair("new win 4", childRenderedMatrixChars.toString(), wantedChildRendered, t)

	matrixCharsComposed := MatrixCharsCompose(windows, windowsChars, []string{"Terminal", "Child"}, "debug")
	wantedComposedRendered := "" +
		"TTTT" + NewLine() +
		"TTCC" // Child has a shifted position (2,1) so Compose has to handle the shifting, too
	compare_str_pair("new win 5", matrixCharsComposed.toString(), wantedComposedRendered, t)

	// replace the layer orders
	windows["Terminal"][KeyLayerNum] = "10" // Child Layer Num is around 1 or 2, not 10
	matrixCharsComposed = MatrixCharsCompose(windows, windowsChars, []string{"Terminal", "Child"}, "debug")
	wantedComposedRendered = "" +
		"TTTT" + NewLine() +
		"TTTT"
	compare_str_pair("new win 6 - layers re-ordered", matrixCharsComposed.toString(), wantedComposedRendered, t)

}

func Test_MatrixNew(t *testing.T) {
	matrixChars := MatrixCharsEmptyOfWindows(3, 2, 'S', "Test_MatrixNew")
	compare_int_pair("MatrixNew 1", len(matrixChars.Rendered), 6, t) // 6 elems are in the matrixChars
	compare_rune_pair("MatrixNew 1", matrixChars.Rendered[MatrixCoord{0, 0}].CharVal, 'S', t)

	wantedRendered := "" +
		"SSS" + NewLine() +
		"SSS"
	compare_str_pair("MatrixNew 2", matrixChars.toString(), wantedRendered, t)

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

func compare_rune_pair(source string, received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR RUNE (%v) - received: %v\n  wanted: %v",
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
