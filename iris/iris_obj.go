package iris

import (
	"fmt"
	"strconv"
)

func NewLine() string { return "\n" }

type RenderedColumn []string
type RenderedScreen []RenderedColumn // there are vertical columns next to each other.

func ScreenEmpty(width, height int) RenderedScreen {
	screen := RenderedScreen{}
	for x := 0; x < width; x++ { // build columns
		column := RenderedColumn{}
		for y := 0; y < height; y++ {
			column = append(column, "")
		}
		screen = append(screen, column)
	}
	return screen
}

/*
I store everything in strings.
 , is a list separator so never use it as a key or a value
 ' ' space is an id separator, don't use space in any win id name
 Windows id characters: [a-zA-Z0-9_-]
*/
type Window map[string]string
type Windows map[string]Window

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
func (win Window) RenderToScreenMatrix() RenderedScreen {
	width := Atoi(win[KeyXright]) - Atoi(win[KeyXleft]) + 1
	height := Atoi(win[KeyXright]) - Atoi(win[KeyXleft]) + 1
	screen := ScreenEmpty(width, height)
	return screen
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
