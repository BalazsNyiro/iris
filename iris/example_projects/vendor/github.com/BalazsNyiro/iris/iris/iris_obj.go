package iris

import (
	"fmt"
	"strconv"
	"strings"
)

func NewLine() string { return "\n" }

// Attr     map[string]string
type Coord [2]int

// a rendered char can have ANSI settings so a single displayed char
// with foreground and background settings can be a long string
type CharRenderedWithFgBgSettings struct {
	colorFgRGB string
	colorBgRGB string
	character  string
}

func (c CharRenderedWithFgBgSettings) toString() string {
	return c.character
}

type MatrixCharsRenderedWithFgBgSettings map[Coord]CharRenderedWithFgBgSettings

type RenderedScreen struct {
	name                string
	width               int
	height              int
	matrixCharsRendered MatrixCharsRenderedWithFgBgSettings
}

func (screen RenderedScreen) toString() string {
	out := []string{}
	for y := 0; y < screen.height; y++ {
		if len(out) > 0 {
			out = append(out, NewLine())
		}
		for x := 0; x < screen.width; x++ {
			coordinate := Coord{x, y}
			out = append(out, screen.matrixCharsRendered[coordinate].toString())
		}
	}
	return strings.Join(out, "")
}

////////////////////////////////////////////////////////////////////////////////////

func ScreenEmpty(width, height int, defaultScreenFiller, name string) RenderedScreen {
	screen := RenderedScreen{width: width, height: height, name: name, matrixCharsRendered: MatrixCharsRenderedWithFgBgSettings{}}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			coordinate := Coord{x, y}
			screen.matrixCharsRendered[coordinate] = CharRenderedWithFgBgSettings{character: defaultScreenFiller}
		}
	}
	return screen
}

// TODO: TEST IT
func ScreenSrcLoad(screen RenderedScreen, width, height int, src, srcType string, lineBreakIfTxtTooLong bool) RenderedScreen {
	// based on srcType and src the screen is modified here
	if srcType == "simpleText" {
		y := 0
		x := 0
		for _, runeNow := range src {
			if lineBreakIfTxtTooLong && x > width && y < height-1 {
				x = 0
				y = y + 1
			}
			if y < height && x < width {
				coordinate := Coord{x, y}
				screen.matrixCharsRendered[coordinate] = CharRenderedWithFgBgSettings{character: string(runeNow)}
				x = x + 1
			}
		}
	}
	return screen
}

func ScreensCompose(windows Windows, winNamesToComposite []string, screenFiller string) RenderedScreen {
	widthMax, heightMax := 0, 0

	// FIXME: windows rendering is based on Name list order.
	// it means if you change the order, a win can overlap another one.
	// maybe it would be better to give a LayerNum into the windows and render them based on that value??

	// keep the original order because the later rendered win overlaps the previous ones
	winNamesToComposite = win_names_keep_publics(winNamesToComposite, false)

	if true { // This part is to find the max width/height only. //////////////
		screensOfWindows := []RenderedScreen{}
		for _, winName := range winNamesToComposite { // default filler: we want to detect the width/height only
			screensOfWindows = append(screensOfWindows, windows[winName].RenderToScreenOfWin("default"))
		}
		for _, screen := range screensOfWindows {
			if screen.width > widthMax {
				widthMax = screen.width
			}
			if screen.height > heightMax {
				heightMax = screen.height
			}
		}
	} // This part is to find the max width/height only. //////////////

	composed := RenderedScreen{width: widthMax, height: heightMax, name: "composed", matrixCharsRendered: MatrixCharsRenderedWithFgBgSettings{}}

	for _, winName := range winNamesToComposite {
		screen := windows[winName].RenderToScreenOfWin(screenFiller)
		for yInWin := 0; yInWin < screen.height; yInWin++ {
			for xInWin := 0; xInWin < screen.width; xInWin++ {
				coordInWinLocal := Coord{xInWin, yInWin}
				coordInRootTerminal := Coord{
					Atoi(windows[winName][KeyXleftCalculated]) + xInWin,
					Atoi(windows[winName][KeyYtopCalculated]) + yInWin}
				composed.matrixCharsRendered[coordInRootTerminal] = screen.matrixCharsRendered[coordInWinLocal]
			}
		}

	}
	return composed
}

/*
I store everything in strings.
 , is a list separator so never use it as a key or a value
 ' ' space is an id separator, don't use space in any win id name
 Windows id characters: [a-zA-Z0-9_-]
*/

//////////////////////////// WINDOWS ////////////////////////////////////////////////////////
type Window map[string]string
type Windows map[string]Window

// TESTED in Test_new_window
func (win Window) RenderToScreenOfWin(screenFillerChar string) RenderedScreen {
	if screenFillerChar == "debug" {
		screenFillerChar = win[KeyDebugWindowFillerChar]
	}
	width := Atoi(win[KeyXrightCalculated]) - Atoi(win[KeyXleftCalculated]) + 1
	height := Atoi(win[KeyYbottomCalculated]) - Atoi(win[KeyYtopCalculated]) + 1
	screen := ScreenEmpty(width, height, screenFillerChar, KeyWinId+":"+win[KeyWinId])
	autoLineBreakAtWinEnd := true
	screen = ScreenSrcLoad(screen, width, height, win[KeyWinContentSrc], win[KeyWinContentType], autoLineBreakAtWinEnd)
	return screen
}

// TESTED
func WinCoordsCalculate(windows Windows) Windows {
	for winName, _ := range windows_keep_publics(windows) {
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
var KeyWinId = "winId"
var KeyWinContentSrc = "winContentSrc"
var KeyWinContentType = "winContentType"

// TESTED in Test_new_window
func WindowsNewState(terminalWidth, terminalHeight int) Windows {
	Win := Windows{}
	// prgState contains all general data
	Win2 := WinNew(Win, "prgState", "-1", "-1", "0", "0", "S")
	Win2["prgState"]["winActiveId"] = ""
	return WinNew(Win2, "Terminal", "0", "0",
		strconv.Itoa(terminalWidth-1),
		strconv.Itoa(terminalHeight-1),
		"T",
	)
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
		errMsg = errMsg + "express param left id < 0:" + Itoa(idValueLeft) + ";"
	}
	if idValueRight < 0 {
		errMsg = errMsg + "express param right id < 0:" + Itoa(idValueRight) + ";"
	}
	idMax := len(tokens) - 1
	idMaxStr := Itoa(idMax)
	if idValueLeft > idMax {
		errMsg = errMsg + "express param left id > len(tokens)-1:" + Itoa(idValueLeft) + " len tokens: " + idMaxStr + ";"
	}
	if idValueRight > idMax {
		errMsg = errMsg + "express param right id > len(tokens)-1:" + Itoa(idValueRight) + " len tokens: " + idMaxStr + ";"
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
	fmt.Println("======= simple expression eval:", exp, "==========")

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
		windows[id] = Window{}
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

			KeyWinId: id,

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
			KeyWidthCalculated:       Itoa(Atoi(keyXright) - Atoi(keyXleft) + 1),
			KeyHeightCalculated:      Itoa(Atoi(keyYbottom) - Atoi(keyYtop) + 1),
			KeyDebugWindowFillerChar: debugWindowFiller,
			KeyWinContentSrc:         "",
			KeyWinContentType:        "simpleText",
		}
	}
	return windows
}
