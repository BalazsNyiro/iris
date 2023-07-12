// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	"os"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 10
var TimeIntervalTerminalSizeDetectMillisec = 100

type Windows map[string]Window
type Window struct {
	id string

	// top-left coord: 0, 0 in the root terminal
	yTop              int
	yBottom           int
	xLeft             int
	xRight            int
	lines             []string
	backgroundDefault string
}

type ScreenLayers []ScreenLayer
type ScreenLayer struct {
	xLeft  int
	yTop   int
	matrix []ScreenColumn
}
type ScreenColumn []ScreenChar
type ScreenChar struct {
	txtValue string
}

// newLineSeparator: \n, \r\n - that you can find at the end of the lines - a line separator
func UserInterfaceStart(ch_data_input chan string, newlineSeparator string) {
	ui_init()
	ch_user_input := make(chan string)
	go channel_read_user_input(ch_user_input)

	ch_terminal_size_change_detect := make(chan [2]int)
	go channel_read_terminal_size_change_detect(ch_terminal_size_change_detect)

	// windows is a read-only variable everywhere,
	windows := Windows{} // modified/updated ONLY here:
	go data_input_interpret(ch_data_input, &windows, newlineSeparator)

	widthSysNow, heightSysNow := TerminalDimensionsWithSyscall()
	terminalSizeActual := [2]int{widthSysNow, heightSysNow}

	loopCounter := 0
	for {
		loopCounter++

		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			fmt.Println("Keys pressed:", stdin)
			action = action_of_user_input(stdin)

		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			terminalSizeActual = terminal_size_change
		default: //               or not coming
			_ = ""
		}

		if action == "quit" {
			UserInterfaceExit()
			break
		}
		layers := RenderAllWindowsIntoLayers(windows, terminalSizeActual)
		// fmt.Println("layers:", layers)
		DisplayAllLayers(layers, newlineSeparator, loopCounter)
		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////

func DisplayAllLayers(layers ScreenLayers, newlineSeparator string, loopCounter int) {
	// naive, TODO: display layers in order?
	// https://stackoverflow.com/questions/5367068/clear-a-terminal-screen-for-real/5367075#5367075

	// fmt.Print("\x1b")                        // clear all the screen
	fmt.Println("clear screen", loopCounter) // clear all the screen

	// size of biggest window + x/y positions
	widthMax, heightMax := 0, 0
	for _, layerStruct := range layers {
		matrix := layerStruct.matrix
		height := len(matrix[0]) + layerStruct.yTop // len of first column
		width := len(matrix) + layerStruct.xLeft
		heightMax = IntMax(heightMax, height)
		widthMax = IntMax(widthMax, width)
	}

	screenMerged := ScreenLayerCreate(
		0, 0,
		widthMax, heightMax, " ")

	for _, layerStruct := range layers {
		matrix := layerStruct.matrix
		height := len(matrix[0]) // len of first column
		width := len(matrix)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				xCalculated := layerStruct.xLeft + x
				yCalculated := layerStruct.yTop + y
				screenMerged.matrix[xCalculated][yCalculated].txtValue = layerStruct.matrix[x][y].txtValue
			}
		}
	}

	///////// DISPLAY merged layer //////////////////
	for y := 0; y < heightMax; y++ {
		for x := 0; x < widthMax; x++ {
			fmt.Print(screenMerged.matrix[x][y].txtValue)
		}
		fmt.Print(newlineSeparator)
	}
}

func ScreenLayerCreate(xLeft, yTop, width, height int, txtLayerDefault string) ScreenLayer {
	// fmt.Println("screen layer create:", xLeft, yTop, width, height)
	screenLayerNew := ScreenLayer{xLeft: xLeft, yTop: yTop}
	for x := 0; x < width; x++ {
		column := ScreenColumn{}
		for y := 0; y < height; y++ {
			column = append(column, ScreenChar{txtValue: txtLayerDefault})
		}
		screenLayerNew.matrix = append(screenLayerNew.matrix, column)
	}
	// fmt.Println("new layer:", screenLayerNew)
	return screenLayerNew
}

func RenderAllWindowsIntoLayers(windowsRO Windows, terminalSize [2]int) ScreenLayers {
	// fmt.Println("terminal size:", terminalSize)
	screenBackground := ScreenLayerCreate(
		0, 0,
		terminalSize[0],
		terminalSize[1], ".")

	layers := ScreenLayers{screenBackground}

	for _, windows := range windowsRO {
		fmt.Println("render: winId", windows)
		screenNow := ScreenLayerCreate(
			3, 2,
			8,
			6, "a")
		layers = append(layers, screenNow)
	}
	return layers
}

func data_input_interpret(ch_data_input chan string, windows *Windows, newlineSeparator string) {
	for {
		select {
		case dataInput, _ := <-ch_data_input:
			// fmt.Println("data input:", dataInput)
			if strings.HasPrefix(dataInput, "select:win") {
				winUpdated := select_win(dataInput, windows, newlineSeparator)
				if winUpdated != "" {
					// fmt.Println("after select:win, addSimpleText", (*windows)[winUpdated].lines)
				}
			}
		default:
			_ = ""
		}
	}
}

/*
'add:simpleText:' is always the last added elem, everything after it is added automatically
into the lines
*/
func select_win(dataInput string, windows *Windows, newlineSeparator string) string {
	winId := ""
	addSimpleTextDetectedLine := -1

	for lineNum, lineOrig := range strings.Split(dataInput, newlineSeparator) {
		line := strings.TrimSpace(lineOrig)
		// fmt.Println("select_win, line:", line)
		elems := strings.Split(line, ":")

		if len(elems) == 3 {

			// select:win:nameOfWin
			if elems[0] == "select" && elems[1] == "win" {
				winId = strings.TrimSpace(elems[2])
				if _, exist := (*windows)[winId]; !exist {
					(*windows)[winId] = Window{}
				}
			}

			win := (*windows)[winId]
			if elems[0] == "set" && elems[1] == "backgroundDefault" {
				win.backgroundDefault = elems[2]
			}
			if elems[0] == "set" && elems[1] == "xLeft" {
				win.xLeft = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "xRight" {
				win.xRight = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "yTop" {
				win.yTop = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "yBottom" {
				win.yBottom = Str2Int(elems[2])
			}
			(*windows)[winId] = win

		}

		if winId == "" {
			continue
		}

		// process only the first line here, then later add all other lines, too
		if elems[0] == "add" && elems[1] == "simpleText" {
			addSimpleTextDetectedLine = lineNum
			win := (*windows)[winId]
			win.lines = append(win.lines, strings.SplitN(lineOrig, "add:simpleText:", -1)[1])
			(*windows)[winId] = win
			break
		}

	}
	if addSimpleTextDetectedLine > -1 {
		for lineNum, lineOrig := range strings.Split(dataInput, newlineSeparator) {
			if lineNum > addSimpleTextDetectedLine {
				win := (*windows)[winId]
				win.lines = append(win.lines, lineOrig)
				(*windows)[winId] = win
			}
		}
	}
	return winId
}

func action_of_user_input(stdin string) string {
	action := ""
	if stdin == "q" {
		action = "quit"
	}

	if stdin == "l" {
	}
	if stdin == "h" {
	}
	if stdin == "j" {
	}
	if stdin == "k" {
	}
	return action
}

func ui_init() {
	terminal_console_clear()
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func UserInterfaceExit() {
	terminal_console_character_show()
	terminal_console_input_buffering_enable()
}

// /////////////////////////////////////////////////
// keypress detection is based on this example:
// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
// thank you.
func channel_read_user_input(ch chan string) {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		ch <- string(b)
	}
} ///////////////////////////////////////////////////

func channel_read_terminal_size_change_detect(ch chan [2]int) {
	widthSys, heightSys := 0, 0
	for {
		widthSysNow, heightSysNow := TerminalDimensionsWithSyscall()
		if widthSysNow != widthSys || heightSysNow != heightSys {
			widthSys = widthSysNow
			heightSys = heightSysNow
			ch <- [2]int{widthSys, heightSys}
		}
		TimeSleep(TimeIntervalTerminalSizeDetectMillisec)
	}
}
