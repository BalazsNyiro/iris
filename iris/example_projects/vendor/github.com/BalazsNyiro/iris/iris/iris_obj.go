package iris

import (
	"fmt"
	"strconv"
	"strings"
)

func NewLine() string { return "\n" }

// Attr     map[string]string
type Coord [2]int

type MatrixCharsRenderedWithFgBgSettings map[Coord]CharObj

type MatrixChars struct {
	name     string
	width    int
	height   int
	Rendered MatrixCharsRenderedWithFgBgSettings
}

func (matrixChars MatrixChars) toString() string {
	out := []string{}
	for y := 0; y < matrixChars.height; y++ {
		if len(out) > 0 {
			out = append(out, NewLine())
		}
		for x := 0; x < matrixChars.width; x++ {
			coordinate := Coord{x, y}
			out = append(out, matrixChars.Rendered[coordinate].render())
		}
	}
	return strings.Join(out, "")
}

// //////////////////////////////////////////////////////////////////////////////////
// internal function, it uses integer width, height values because of calculations (avoid conversion)
// TESTED
func MatrixCharsEmptyOfWindows(width, height int, matrixFiller rune, winName string) MatrixChars {
	matrixChars := MatrixChars{width: width, height: height, name: winName, Rendered: MatrixCharsRenderedWithFgBgSettings{}}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			coordinate := Coord{x, y}
			matrixChars.Rendered[coordinate] = CharObjNew(matrixFiller)
		}
	}
	return matrixChars
}

// //////////////////////////////////////////////////////////////////////////////////////
// internal function, it uses integer width, height values because of calculations (avoid conversion)
// TESTED
func MatrixCharsInsertContentOfWindows(matrixChars MatrixChars, winWidth, winHeight int, windowChars WindowChars, lineBreakIfTxtTooLong bool) MatrixChars {
	x, y := 0, 0 // x starts with 0. If x == winWidth it means that x is not inside the windows area because of the 0 based
	idNext := 0  // create the variable only once
	newLine := NewLine()
	newLineLen := len(newLine)
	newLineRune := rune(newLine[0])

	for id := 0; id < len(windowChars); id++ { // counting. so if x == winWidth then you have to move into the next line.
		idNext = id + 1
		charObj := windowChars[id]

		if true { // ##################################### NEWLINE HANDLING ###############
			isNewLine := false

			if charObj.CharVal == '\r' && newLine == "\r\n" && idNext < len(windowChars) {
				if windowChars[idNext].CharVal == '\n' {
					isNewLine = true
					id++ // skip the next \r
				}
			}

			if newLineLen == 1 && charObj.CharVal == newLineRune {
				isNewLine = true // work with \r, \n newlines too
			}

			if isNewLine { // and the newline char objects are not added into the matrixChars, because
				x = 0 //      they are represented with the increased y value.
				y += 1
				continue
			}
		} // ###################################### NEWLINE HANDLING #############################

		if lineBreakIfTxtTooLong && (x == winWidth) && (y < winHeight-1) {
			x = 0
			y += 1
		}
		if y < winHeight && x < winWidth { // with 'x <', 'y <' it copies the visible part only
			coordinate := Coord{x, y}
			matrixChars.Rendered[coordinate] = charObj
			x += 1
		}
	}
	return matrixChars
} ////////////////////////////////////////////////////////////////////////////////////////

func MatrixCharsCompose(windows Windows, windowsChars WindowsChars, winNamesToComposite []string, matrixFiller string) MatrixChars {
	widthMax, heightMax := 0, 0

	winNamesToComposite = WinNamesKeepPublic(winNamesToComposite, false)

	composed := MatrixChars{width: widthMax, height: heightMax, name: "composed", Rendered: MatrixCharsRenderedWithFgBgSettings{}}

	// lower Layer number is rendered earlier
	for _, winName := range WinNamesSort(windows, winNamesToComposite, KeyLayerNum, "number") {
		matrixActualWin := windows[winName].RenderToMatrixCharsOfWin(windowsChars, matrixFiller)

		winLocalXLeftCalculated := Str2Int(windows[winName][KeyXleftCalculated]) // read the values only once,
		winLocalYTopCalculated := Str2Int(windows[winName][KeyYtopCalculated])   // avoid to be changed in the for loop

		for yInWin := 0; yInWin < matrixActualWin.height; yInWin++ {
			for xInWin := 0; xInWin < matrixActualWin.width; xInWin++ {
				coordInWinLocal := Coord{xInWin, yInWin}
				coordInRootTerminal := Coord{winLocalXLeftCalculated + xInWin, winLocalYTopCalculated + yInWin}
				composed.Rendered[coordInRootTerminal] = matrixActualWin.Rendered[coordInWinLocal]
			}
		}
		// +1: because the coords are 0 based numbers, which means `if x == 9` then width = 10
		// -1: because the matrix's first position is similar with the Calculated's last position
		// so one character is double calculated
		composed.width = IntMax(composed.width, +1+winLocalXLeftCalculated+matrixActualWin.width-1)
		composed.height = IntMax(composed.height, +1+winLocalYTopCalculated+matrixActualWin.height-1)
	}
	return composed
}

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

// shift: fix simple Coord modifier
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
	Win := Windows{}
	// prgState contains all general data
	Win2 := WinNew(Win, "prgState", "-1", "-1", "0", "0", "S")
	Win2["prgState"]["winActiveId"] = ""

	return WinNew(Win2, "Terminal", "0", "0",
		strconv.Itoa(terminalWidth-1),
		strconv.Itoa(terminalHeight-1),
		"T",
	), WindowsChars{}
}

// if the next operator is "": there is no more operator
// TESTED
func TokenOperatorNext(tokens []string) (int, string) {
	operatorNext := "unknown"
	tokens = StrListRemoveEmptyElems(tokens, true)
	// math operator precedence: * / are the first
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("*,/", token) {
			operatorNext = token
			return id, operatorNext
		}
	}
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("+,-", token) {
			operatorNext = token
			return id, operatorNext
		}
	}
	return -1, operatorNext
}

func TokenReplaceWinPlaceholders(windows Windows, tokens []string) []string {
	tokens = StrListRemoveEmptyElems(tokens, true)
	for id, token := range tokens {
		if len(token) > 4 && token[0:4] == "win:" { // win:Terminal:xRightCalculated
			splitted := strings.Split(token, ":")
			winName := splitted[1]
			attrib := splitted[2]

			tokens[id] = "0" // set the normal value if key/attrib exists:
			if winObj, keyInMap := windows[winName]; keyInMap {
				if valueAttrib, attribInMap := winObj[attrib]; attribInMap {
					tokens[id] = valueAttrib
				}
			}
		}
	}
	return tokens
}

// TESTED
func ParametersCollect(tokens []string, tokenId int) (string, string, int, int, string) {
	errMsg := ""
	valueLeft := ""
	valueRight := ""

	idValueLeft, idValueRight := tokenId-1, tokenId+1
	if idValueLeft < 0 {
		errMsg = errMsg + "express param left id < 0:" + Int2Str(idValueLeft) + ";"
	}
	if idValueRight < 0 {
		errMsg = errMsg + "express param right id < 0:" + Int2Str(idValueRight) + ";"
	}
	idMax := len(tokens) - 1
	idMaxStr := Int2Str(idMax)
	if idValueLeft > idMax {
		errMsg = errMsg + "express param left id > len(tokens)-1:" + Int2Str(idValueLeft) + " len tokens: " + idMaxStr + ";"
	}
	if idValueRight > idMax {
		errMsg = errMsg + "express param right id > len(tokens)-1:" + Int2Str(idValueRight) + " len tokens: " + idMaxStr + ";"
	}
	if errMsg == "" {
		valueLeft = tokens[idValueLeft]
		valueRight = tokens[idValueRight]
	}

	return valueLeft, valueRight, idValueLeft, idValueRight, errMsg
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

		valueLeft, valueRight, idValueLeft, idValueRight, idError := ParametersCollect(tokens, id)
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
func WinNew(windows Windows, id, keyXleft, keyYtop, keyXright, keyYbottom, debugWindowFiller string) Windows {
	if id == "prgState" {
		// program state is not a real window,
		// it stores the current settings
		windows[id] = Window{
			KeyLayerNum: Int2Str(len(windows)),
		}
	} else {
		windows[id] = Window{

			// RULES:
			// the coords are 0 based. so (0, 0) represents the top-left Coord.

			// the not calculated values can be fix or relative values (20%) for example,
			// or later complex functions
			// the positions can be simple fix numbers: 20

			/////////// COMPLEX EXPRESSIONS ///////////////////////////////
			// complex expressions: "win:id_without_space * 0.8 - win:id2 / 2"
			// known operators: * / + -
			// minimum one space is mandatory between all elems as a separator
			// win:  a win id: [a-zA-Z0-9-_.]

			// IMPORTANT: All four Coord values can have different expressions!
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
