// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	"os"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 1000 // 10 is the prod value
var TimeIntervalTerminalSizeDetectMillisec = 100        // 100 is the prod value

// Char is the smallest object.
// Future: a complex obj with foreground/bg colors, display attributes
type Char struct {
	runeVal rune
}

func (c Char) display() rune {
	return c.runeVal
}

/////////////////////////////////////////////////////////////////

type Line []Char

func (line Line) LineToStr() string {
	out := []rune{}
	for _, Char := range line {
		out = append(out, Char.runeVal)
	}
	return string(out)
}
func LineFromStr(txt string) Line {
	Line := Line{}
	for _, runeVal := range txt {
		Line = append(Line, Char{runeVal: runeVal})
	}
	return Line
}

/////////////////////////////////////////////////////////////////

type MessageAndCharacters struct {
	msg     string
	addLine Line
}

// ///////////////////////////////////////////////////////////////
type Windows map[string]Window

func (wins Windows) printAll() {
	for _, win := range wins {
		win.print()
	}
}

// A window is a logical unit.
// It has settings, and content, but doesn't now
// anything how it will be rendered
type Window struct {
	id string

	// top-left coord: 0, 0 in the root terminal
	yTop              int
	xLeft             int
	width             int
	height            int
	lines             []Line
	backgroundDefault string
	winId             string
}

func (w Window) print() {
	fmt.Println("winId:  ", w.winId)
	for _, line := range w.lines {
		fmt.Println("winLine:", line.LineToStr())
	}
}

/////////////////////////////////////////////////////////////////

type ScreenLayers []ScreenLayer_CharMatrix

type ScreenLayer_CharMatrix struct {
	xLeft        int
	yTop         int
	matrix       []ScreenColumn
	creationinfo string
}

func (sl ScreenLayer_CharMatrix) print() {
	fmt.Println("\nlayerId:", sl.creationinfo)
	matrixHeight := len(sl.matrix[0])
	for y := 0; y < matrixHeight; y++ {
		for _, column := range sl.matrix {
			fmt.Print(column[y])
		}
		fmt.Println()
	}
}

func (layer ScreenLayer_CharMatrix) layerToTxt(lineSep string) string {
	yMax := len(layer.matrix[0])
	columns := layer.matrix

	outputRunes := []rune{}
	for y := 0; y < yMax; y++ {
		for _, column := range columns {
			outputRunes = append(outputRunes, column[y].display())
		}
		if y < yMax-1 {
			for _, r := range []rune(lineSep) { // theoretically the line separator
				outputRunes = append(outputRunes, r) // can be  \r\n, too, more than one char
			}
		}
	}
	return string(outputRunes)
}

type ScreenColumn []Char

/////////////////////////////////////////////////////////////////

func UserInterfaceStart(ch_data_input chan MessageAndCharacters, dataInputLineSeparator string) {
	ui_init()
	ch_user_input := make(chan string)
	go channel_read_user_input(ch_user_input)

	ch_terminal_size_change_detect := make(chan [2]int)
	go channel_read_terminal_size_change_detect(ch_terminal_size_change_detect)

	// windows is a read-only variable everywhere,
	windows := Windows{} // modified/updated ONLY here:
	go dataInputInterpret(ch_data_input, &windows, dataInputLineSeparator)

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
		fmt.Println("windows: ", windows)
		layers := LayersRenderFromWindows(windows, terminalSizeActual)
		// fmt.Println("layers:", layers)
		layersDisplayAll(layers, dataInputLineSeparator, loopCounter)
		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////

// TESTED
func LayerCreate(xLeft, yTop, width, height int, txtLayerDefault string, creationInfo string) ScreenLayer_CharMatrix {
	// fmt.Println("screen layer create:", xLeft, yTop, width, height)
	screenLayerNew := ScreenLayer_CharMatrix{xLeft: xLeft, yTop: yTop, creationinfo: creationInfo}
	defaultRune := 'r'
	if len(txtLayerDefault) > 0 {
		defaultRune = rune(txtLayerDefault[0])
	}
	for x := 0; x < width; x++ {
		column := ScreenColumn{}
		for y := 0; y < height; y++ {
			column = append(column, Char{runeVal: defaultRune})
		}
		screenLayerNew.matrix = append(screenLayerNew.matrix, column)
	}
	// fmt.Println("new layer:", screenLayerNew)
	return screenLayerNew
}

func LayersRenderFromWindows(windowsRO Windows, terminalSize [2]int) ScreenLayers {
	fmt.Println("terminal size:", terminalSize)
	layerBackground := LayerCreate(
		0, 0,
		terminalSize[0],
		terminalSize[1], ".", "layerBackground")

	layers := ScreenLayers{layerBackground}

	for _, win := range windowsRO {
		fmt.Println("render: winId >", win.winId, "< xLeft:", win.xLeft, "yTop:", win.yTop, "width:", win.width, "height:", win.height)
		screenNow := LayerCreate(
			win.xLeft, win.yTop,
			win.width, win.height, win.backgroundDefault, win.winId)

		// structure the character input into one COLUMN, (a visible structure)
		textBlockVisible := []Line{Line{}}
		lineNumFirstVisible := len(win.lines) - win.height

		// this solution handles the too long lines, and flow the text into the next line
		for lineNumInWin, lineReceived := range win.lines {
			if lineNumInWin >= lineNumFirstVisible {
				textBlockLastLineId := len(textBlockVisible) - 1
				LineActual := textBlockVisible[textBlockLastLineId]

				for _, charNow := range lineReceived {
					LineActual = append(LineActual, charNow)
					// if we are in the last available column now:
					if len(LineActual) == len(screenNow.matrix) {
						textBlockVisible[textBlockLastLineId] = LineActual
						textBlockVisible = append(textBlockVisible, Line{})
						textBlockLastLineId = len(textBlockVisible) - 1
						LineActual = textBlockVisible[textBlockLastLineId]
					}

				}
				textBlockVisible[textBlockLastLineId] = LineActual
				textBlockVisible = append(textBlockVisible, Line{})
			}
		}

		// Load the theoretically visible text into the windows' column structure
		len_textBlockVisible := len(textBlockVisible)
		lineNumFirstVisibleInWindows := 0
		if len_textBlockVisible > win.height {
			lineNumFirstVisibleInWindows = len_textBlockVisible - win.height
		}

		// fmt.Println("len textBlockVisible:", len_textBlockVisible)
		// fmt.Println("          win.height:", win.height)
		// fmt.Println(" lineNumFirstVisible:", lineNumFirstVisibleInWindows)

		for lineNum, lineVisible := range textBlockVisible {
			if lineNum >= lineNumFirstVisibleInWindows {
				for charPos, charNow := range lineVisible {
					column := screenNow.matrix[charPos]
					lineNumInWindows := lineNum - lineNumFirstVisibleInWindows
					column[lineNumInWindows] = charNow
					screenNow.matrix[charPos] = column
				}
			}

		}
		layers = append(layers, screenNow)
	}
	return layers
}

func layersDisplayAll(layers ScreenLayers, newlineSeparator string, loopCounter int) {
	// naive, TODO: display layers in order?
	// https://stackoverflow.com/questions/5367068/clear-a-terminal-screen-for-real/5367075#5367075

	// fmt.Print("\x1b")                        // clear all the screen
	fmt.Println("clear screen", loopCounter) // clear all the screen

	// size of biggest window + x/y positions
	widthMax, heightMax := 0, 0
	for _, layerStruct := range layers {
		matrix := layerStruct.matrix
		if len(matrix) < 1 {
			continue
		}
		width := len(matrix) + layerStruct.xLeft
		height := len(matrix[0]) + layerStruct.yTop // len of first column
		heightMax = IntMax(heightMax, height)
		widthMax = IntMax(widthMax, width)
	}

	screenMerged := LayerCreate(
		0, 0,
		widthMax, heightMax, " ", "screenMerged")

	for _, layerStruct := range layers {
		matrix := layerStruct.matrix
		width := len(matrix)
		if width < 1 {
			continue
		}
		height := len(matrix[0]) // len of first column
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				xCalculated := layerStruct.xLeft + x
				yCalculated := layerStruct.yTop + y
				screenMerged.matrix[xCalculated][yCalculated].runeVal = layerStruct.matrix[x][y].runeVal
			}
		}
	}

	///////// DISPLAY merged layer //////////////////
	for y := 0; y < heightMax; y++ {
		for x := 0; x < widthMax; x++ {
			fmt.Print(screenMerged.matrix[x][y].display())
		}
		fmt.Print(newlineSeparator)
	}
}

func dataInputInterpret(ch_data_input chan MessageAndCharacters, windows *Windows, dataInputLineSeparator string) {
	for {
		select {
		case dataInput, _ := <-ch_data_input:
			fmt.Println("\ndata input:", dataInput)
			if strings.HasPrefix(strings.TrimSpace(dataInput.msg), "select:win") {
				dataInputProcessLineByLine(dataInput, windows, dataInputLineSeparator)
			}
		default:
			_ = ""
		}
	}
}

func dataInputProcessLineByLine(dataInput MessageAndCharacters, windows *Windows, dataInputLineSeparator string) string {
	winId := ""

	for _, lineOrig := range strings.Split(dataInput.msg, dataInputLineSeparator) {
		line := strings.TrimSpace(lineOrig)
		elems := strings.Split(line, ":")

		if len(elems) == 3 {

			// select:win:nameOfWin
			if elems[0] == "select" && elems[1] == "win" {
				winId = strings.TrimSpace(elems[2])
				if _, exist := (*windows)[winId]; !exist {
					(*windows)[winId] = Window{winId: winId}
				}

				// process only the first line here, then later add all other lines, too
				if len(dataInput.addLine) > 0 {
					win := (*windows)[winId]
					win.lines = append(win.lines, dataInput.addLine)
					(*windows)[winId] = win
				}

			}

			win := (*windows)[winId]
			if elems[0] == "set" && elems[1] == "backgroundDefault" {
				win.backgroundDefault = elems[2]
			}
			if elems[0] == "set" && elems[1] == "xLeft" {
				win.xLeft = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "width" {
				win.width = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "yTop" {
				win.yTop = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "height" {
				win.height = Str2Int(elems[2])
			}
			(*windows)[winId] = win

		}

		if winId == "" {
			continue
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
