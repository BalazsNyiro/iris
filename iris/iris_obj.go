package iris

import (
	"fmt"
	"strconv"
	"strings"
)

func NewLine() string { return "\n" }

// Attr     map[string]string
type coord [2]int
type pixels map[coord]string
type RenderedScreen struct {
	name   string
	width  int
	height int
	pixels pixels
}

func (screen RenderedScreen) toString() string {
	out := []string{}
	for y := 0; y < screen.height; y++ {
		if len(out) > 0 {
			out = append(out, NewLine())
		}
		for x := 0; x < screen.width; x++ {
			coordinate := coord{x, y}
			out = append(out, screen.pixels[coordinate])
		}
	}
	return strings.Join(out, "")
}

////////////////////////////////////////////////////////////////////////////////////

func ScreenEmpty(width, height int, defaultScreenFiller, name string) RenderedScreen {
	screen := RenderedScreen{width: width, height: height, name: name, pixels: pixels{}}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			coordinate := coord{x, y}
			screen.pixels[coordinate] = defaultScreenFiller
		}
	}
	return screen
}

func ScreensComposeToScreen(windows Windows, winNames []string) RenderedScreen {
	widthMax, heightMax := 0, 0

	// This part is to find the max width/height only. //////////////
	screensOfWindows := []RenderedScreen{}
	for _, winName := range winNames {
		screensOfWindows = append(screensOfWindows, windows[winName].RenderToScreenOfWin())
	}

	for _, screen := range screensOfWindows {
		if screen.width > widthMax {
			widthMax = screen.width
		}
		if screen.height > heightMax {
			heightMax = screen.height
		}
	}
	// This part is to find the max width/height only. //////////////

	composed := RenderedScreen{width: widthMax, height: heightMax, name: "composed", pixels: pixels{}}

	for _, winName := range winNames {
		screen := windows[winName].RenderToScreenOfWin()
		for yInWin := 0; yInWin < screen.height; yInWin++ {
			for xInWin := 0; xInWin < screen.width; xInWin++ {
				coordInWinLocal := coord{xInWin, yInWin}
				coordInRootTerminal := coord{
					Atoi(windows[winName][KeyXleftCalculated]) + xInWin,
					Atoi(windows[winName][KeyYtopCalculated]) + yInWin}
				composed.pixels[coordInRootTerminal] = screen.pixels[coordInWinLocal]
			}
		}

	}
	return composed
	// winScreenLocal is a small screen that represents only the window
	// winScreenLocal := windows[winName].RenderToScreenOfWin()
	// screenTerminalSized := ScreenEmpty(width, height, win[KeyDebugWindowFillerChar], KeyWinId+":"+win[KeyWinId])
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

func (win Window) RenderToScreenOfWin() RenderedScreen {
	// TODO: use calculated width/height when they are ready!
	width := Atoi(win[KeyXrightCalculated]) - Atoi(win[KeyXleftCalculated]) + 1
	height := Atoi(win[KeyYbottomCalculated]) - Atoi(win[KeyYtopCalculated]) + 1
	screen := ScreenEmpty(width, height, win[KeyDebugWindowFillerChar], KeyWinId+":"+win[KeyWinId])
	return screen
}

func CalculateAllWindowCoords(windows Windows) Windows {
	for winName, _ := range windows {
		fmt.Println("Calc winName", winName)
		windows[winName][KeyXleftCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyXleft], windows), "+", windows[winName][KeyXshift])
		windows[winName][KeyXrightCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyXright], windows), "+", windows[winName][KeyXshift])
		windows[winName][KeyYtopCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyYtop], windows), "+", windows[winName][KeyYshift])
		windows[winName][KeyYbottomCalculated] = StrMath(CoordExpressionEval(windows[winName][KeyYbottom], windows), "+", windows[winName][KeyYshift])
	}
	return windows
}

/*
func CalculateCoords(win Window) Window {

}

*/

//////////////////////////// WINDOWS ////////////////////////////////////////////////////////

// they can contain complex expressions that has to be calculated
var KeyXleft = "xLeft"
var KeyXright = "xRight"
var KeyYtop = "yTop"
var KeyYbottom = "yBottom"

// shift: fix simple coord modifier
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
func OperatorNext(tokens []string) int {
	// math operator precedence: * / are the first
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("*,/", token) {
			return id
		}
	}
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("+,-", token) {
			return id
		}
	}
	return -1
}

// FIXME: use smaller sub-functions, not a huge monolitic fun
func CoordExpressionEval(exp string, windows Windows) string {
	// FIXME: () handling
	fmt.Println("======= simple expression eval =======")

	// minimum 1 space between expression elems
	// win:WindowsName:Attribute
	// fun:min ( )
	// example "( win:Terminal:KeyXleftCalculated + win:OtherWin:KeyYtopCalculated ) / 2:
	exp = strings.TrimSpace(exp)
	exp = StrDoubleSpacesRemove(exp)

	tokens := strings.Split(exp, " ")

	// replace all windows elem into fix values
	for id, token := range tokens {
		token = strings.TrimSpace(token)
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

	// if a token == "" then it is deleted
	// calculate all operator, and remove left/right values
	id := OperatorNext(tokens)
	for id > -1 {
		operator := tokens[id]
		fmt.Println(">>> operator:", operator)
		// if tokens has the next param for the operator:
		if strings.Contains("+,-,*,/", operator) {
			if len(tokens) > (id + 1) {
				valueLeft := tokens[id-1]
				valueRight := tokens[id+1]
				fmt.Println("  ", valueLeft, operator, valueRight)

				tokens[id-1] = ""
				tokens[id] = StrMath(valueLeft, operator, valueRight)
				tokens[id+1] = ""
			} else {
				fmt.Println("missing operator parameter:", operator)
				return "0" // if the param is missing, return with 0
			}

		} else { // unknown operator
			fmt.Println("unknown operator:", operator)
			return "0" // if the expression has syntax error, return with 0
		}

		tokens = StrListRemoveEmptyElems(tokens, true)

		id = OperatorNext(tokens)
	}

	// at this point there is no more operator in tokens
	// and all operator value is calculated.
	// so tokens has one value only
	return tokens[0]
}

func WinNew(windows Windows, id, keyXleft, keyYtop, keyXright, keyYbottom, debugWindowFiller string) Windows {
	windows[id] = Window{

		// RULES:
		// the coords are 0 based. so (0, 0) represents the top-left coord.

		// the not calculated values can be fix or relative values (20%) for example,
		// or later complex functions
		// the positions can be simple fix numbers: 20

		/////////// COMPLEX EXPRESSIONS ///////////////////////////////
		// complex expressions: "win:id_without_space * 0.8 - win:id2 / 2"
		// known operators: * / + -
		// minimum one space is mandatory between all elems as a separator
		// win:  a win id: [a-zA-Z0-9-_.]

		// IMPORTANT: All four coord values can have different expressions!
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
	}
	return windows
}
