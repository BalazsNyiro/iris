package iris

import "strconv"

func NewLine() string { return "\n" }

type RenderedColumn []string
type RenderedScreen []RenderedColumn // there are vertical columns next to each other.

/*
I store everything in strings.
 , is a list separator so never use it as a key or a value
 ' ' space is an id separator, don't use space in any win id name
 Windows id characters: [a-zA-Z0-9_-]
*/
type Window map[string]string
type Windows map[string]Window

var KeyWidth = "width"
var KeyHeight = "height"
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

func WindowsNewState(terminalWidth, terminalHeight string) Windows {
	Win := Windows{}
	terminal_width, _ := strconv.Atoi(terminalWidth)
	terminal_height, _ := strconv.Atoi(terminalHeight)
	return WinNew(Win, "root", "0", "0",
		strconv.Itoa(terminal_width-1),
		strconv.Itoa(terminal_height-1),
	)
}

func WinNew(windows Windows, id, keyXleft, keyYtop, keyXright, keyYbottom string) Windows {
	windows[id] = Window{

		// RULES:

		// the not calculated values can be fix or relative values (20%) for example,
		// or later complex functions
		KeyXleft:   keyXleft,
		KeyXright:  keyXright,
		KeyYtop:    keyYtop,
		KeyYbottom: keyYbottom,

		// here you can see calculated fix positions only, the actual positions
		KeyXleftCalculated:   "",
		KeyXrightCalculated:  "",
		KeyYtopCalculated:    "",
		KeyYbottomCalculated: "",
		KeyWidthCalculated:   "",
		KeyHeightCalculated:  "",
	}
	return windows
}
