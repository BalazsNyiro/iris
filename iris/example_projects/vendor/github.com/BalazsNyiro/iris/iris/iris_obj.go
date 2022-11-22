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

	screens := []RenderedScreen{}
	for _, winName := range winNames {
		screens = append(screens, windows[winName].RenderToScreen())
	}

	for _, screen := range screens {
		if screen.width > widthMax {
			widthMax = screen.width
		}
		if screen.height > heightMax {
			heightMax = screen.height
		}
	}

	composed := RenderedScreen{width: widthMax, height: heightMax, name: "composed", pixels: pixels{}}

	for _, screen := range screens {
		for y := 0; y < screen.height; y++ {
			for x := 0; x < screen.width; x++ {
				coordinate := coord{x, y}
				composed.pixels[coordinate] = screen.pixels[coordinate]
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

func (win Window) RenderToScreen() RenderedScreen {
	// TODO: use calculated width/height when they are ready!
	width := Atoi(win[KeyXright]) - Atoi(win[KeyXleft]) + 1
	height := Atoi(win[KeyYbottom]) - Atoi(win[KeyYtop]) + 1
	screen := ScreenEmpty(width, height, win[KeyDebugWindowFillerChar], KeyWinId+":"+win[KeyWinId])
	return screen
}

//////////////////////////// WINDOWS ////////////////////////////////////////////////////////

var KeyXleft = "xLeft"
var KeyXright = "xRight"
var KeyYtop = "yTop"
var KeyYbottom = "yBottom"

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
	return WinNew(Win, "Terminal", "0", "0",
		strconv.Itoa(terminalWidth-1),
		strconv.Itoa(terminalHeight-1),
		"T",
	)
}

func Atoi(txt string) int {
	num, error := strconv.Atoi(txt)
	if error == nil {
		return num
	}
	fmt.Println("Atoi error: ", error)
	return 0
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

		// here you can see calculated fix positions only, the actual positions
		KeyXleftCalculated:       "",
		KeyXrightCalculated:      "",
		KeyYtopCalculated:        "",
		KeyYbottomCalculated:     "",
		KeyWidthCalculated:       "",
		KeyHeightCalculated:      "",
		KeyDebugWindowFillerChar: debugWindowFiller,
	}
	return windows
}