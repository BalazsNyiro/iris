package iris

import (
	"fmt"
	"strconv"
	"strings"
)

/*
I store everything in strings.
 , is a list separator so never use it as a key or a value
 ' ' space is an id separator, don't use space in any win id name
 Windows id characters: [a-zA-Z0-9_-]
*/

// ////////////////////////// WINDOWS ////////////////////////////////////////////////////////
type Window map[string]string
type Windows map[string]Window

// I use a separated structure to store char objects because 'Window' can be a free key/value map with this solution.
// If I insert the Chars into Window, I have to define a rigid struct, and I'd like to avoid that.
type WindowChars []CharObj
type WindowsChars map[string]WindowChars

// TESTED in Test_new_window
func (win Window) RenderToMatrixCharsOfWin(windowsChars WindowsChars, matrixFiller string) MatrixChars {

	// in matrixFiller we can pass one char (in the string) or more chars too (debug)
	// and the need of debug selection is the reason why matrixFiller is a string and not a rune
	matrixFillerChar := rune(matrixFiller[0])
	if matrixFiller == "debug" {
		matrixFillerChar = rune(win[KeyDebugWindowFillerChar][0])
	}
	winWidth := Str2Int(win[KeyXrightCalculated]) - Str2Int(win[KeyXleftCalculated]) + 1
	winHeight := Str2Int(win[KeyYbottomCalculated]) - Str2Int(win[KeyYtopCalculated]) + 1
	autoLineBreakAtWinEnd := true

	winName := win[KeyWinName] // read the id out from the window
	matrixChars := MatrixCharsEmptyOfWindows(winWidth, winHeight, matrixFillerChar, KeyWinName+":"+winName)
	matrixChars = MatrixCharsInsertContentOfWindows(matrixChars, winWidth, winHeight, windowsChars[winName], autoLineBreakAtWinEnd)

	return matrixChars
}

// TESTED
func WinSourceUpdate(windowsChars WindowsChars, winName, contentType, contentSrc string) WindowsChars {
	// possible types: simpleText, html
	// or in a different fun later: CharObjects can be loaded directly without text-> CharObj conversion
	if contentType == "simpleText" {
		for _, char := range contentSrc {
			charObj := CharObjNew(char)
			chars := windowsChars[winName] // before the update
			chars = append(chars, charObj)
			windowsChars[winName] = chars
		}
	}
	return windowsChars
}

// TESTED
func WinCoordsCalculateUpdate(windows Windows) Windows {
	for winName, _ := range WindowsKeepPublic(windows) {
		// fmt.Println("Calc winName", winName)
		windows[winName][KeyXleftCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyXleft], windows), "+", windows[winName][KeyXshift])
		windows[winName][KeyXrightCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyXright], windows), "+", windows[winName][KeyXshift])
		windows[winName][KeyYtopCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyYtop], windows), "+", windows[winName][KeyYshift])
		windows[winName][KeyYbottomCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyYbottom], windows), "+", windows[winName][KeyYshift])
	}
	return windows
}

//////////////////////////// WINDOWS ////////////////////////////////////////////////////////

// they can contain complex expressions that has to be calculated
var KeyXleft = "xLeft"
var KeyXright = "xRight"
var KeyYtop = "yTop"
var KeyYbottom = "yBottom"

// shift: fix simple MatrixCoord modifier
// the window has 1 shift value so it's a global for the 4 corner
var KeyXshift = "xShift"
var KeyYshift = "yShift"

// Calculated = (ExpressionResult + shift)
// example: KeyXleftCalculated  = KeyXleft + keyXshift
// example: KeyXrightCalculated  = KeyXright + keyXshift

var KeyWidthCalculated = "widthCalculated"
var KeyHeightCalculated = "heightCalculated"
var KeyXleftCalculated = "xLeftCalculated"
var KeyXrightCalculated = "xRightCalculated"
var KeyYtopCalculated = "yTopCalculated"
var KeyYbottomCalculated = "yBottomCalculated"
var KeyDebugWindowFillerChar = "debugWindowFillerChar"
var KeyWinName = "winName"
var KeyVisible = "visible"
var KeyLayerNum = "LayerRenderNum" // smaller is rendered first

// TESTED in Test_new_window
func WindowsNewState(terminalWidth, terminalHeight int) (Windows, WindowsChars) {
	windows := Windows{}
	// prgState contains all general data
	windows = WinCreateIntoWindows(windows, "prgState", "-1", "-1", "0", "0", "S")
	windows["prgState"]["winActiveId"] = ""

	windows = WinCreateIntoWindows(windows, "Terminal", "0", "0",
		strconv.Itoa(terminalWidth-1),
		strconv.Itoa(terminalHeight-1),
		"T",
	)

	return windows, WindowsChars{}
}

// TESTED
func CoordExpressionEval(exp string, windows Windows) string {
	// TODO: () handling
	// fmt.Println("======= simple expression eval:", exp, "==========")

	// minimum 1 space between expression elems
	// win:WindowsName:Attribute
	// fun:min ( )
	// example "win:Terminal:KeyXleftCalculated + win:OtherWin:KeyYtopCalculated / 2:

	tokens := strings.Split(exp, " ")
	tokens = TokenReplaceWinPlaceholders(windows, tokens)

	// if a token == "" then it is deleted
	// calculate all operator, and remove left/right values
	id, operator := TokenOperatorNext(tokens)
	for id > -1 {
		fmt.Println(">>> operator:", operator)
		if !ExprOperatorIsValid(operator) {
			fmt.Println("unknown operator:", operator)
			return "0" // if the expression has syntax error, return with 0
		}

		valueLeft, valueRight, idValueLeft, idValueRight, idError := TokenParametersCollect(tokens, id)
		if idError != "" {
			fmt.Println("operator parameter error:", idError)
			return "0" // if the param is missing, return with 0
		}

		tokens[idValueLeft] = "" // clean params, overwrite operator with result
		tokens[id] = StrMath(valueLeft, operator, valueRight)
		tokens[idValueRight] = ""

		tokens = StrListRemoveEmptyElems(tokens, true)
		id, operator = TokenOperatorNext(tokens)
	}

	// at this point there is no more operator in tokens
	// and all operator value is calculated.
	// so tokens has one value only
	return tokens[0]
}

// TESTED
func WinCreateIntoWindows(windows Windows, id, keyXleft, keyYtop, keyXright, keyYbottom, debugWindowFiller string) Windows {
	if id == "prgState" {
		// program state is not a real window,
		// it stores the current settings
		windows[id] = Window{
			KeyLayerNum: Int2Str(len(windows)),
		}
	} else {
		windows[id] = Window{

			// RULES:
			// the coords are 0 based. so (0, 0) represents the top-left MatrixCoord.

			// the not calculated values can be fix or relative values (20%) for example,
			// or later complex functions
			// the positions can be simple fix numbers: 20

			/////////// COMPLEX EXPRESSIONS ///////////////////////////////
			// complex expressions: "win:id_without_space * 0.8 - win:id2 / 2"
			// known operators: * / + -
			// minimum one space is mandatory between all elems as a separator
			// win:  a win id: [a-zA-Z0-9-_.]

			// IMPORTANT: All four MatrixCoord values can have different expressions!
			// it means you can use more parent windows as parents, or other values,
			// the time for example.

			// TODO: fun:function_name(params) - user can call functions to calculate
			// the position. example: move coordinates based on current seconds.

			KeyWinName: id,

			KeyXleft:   keyXleft,
			KeyXright:  keyXright,
			KeyYtop:    keyYtop,
			KeyYbottom: keyYbottom,

			KeyXshift: "0",
			KeyYshift: "0",

			// here you can see calculated fix positions only, the actual positions
			KeyXleftCalculated:       keyXleft,  // initially, before first calculation
			KeyXrightCalculated:      keyXright, // use these values
			KeyYtopCalculated:        keyYtop,
			KeyYbottomCalculated:     keyYbottom,
			KeyWidthCalculated:       Int2Str(Str2Int(keyXright) - Str2Int(keyXleft) + 1),
			KeyHeightCalculated:      Int2Str(Str2Int(keyYbottom) - Str2Int(keyYtop) + 1),
			KeyDebugWindowFillerChar: debugWindowFiller,

			KeyVisible: "true",

			// the default render layer num is follow the natural windows creation
			KeyLayerNum: Int2Str(len(windows)),
		}
	}
	return windows
}
