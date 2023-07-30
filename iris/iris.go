/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package iris

import (
	"errors"
	"fmt"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 1000 // 10 is the prod value
var TimeIntervalTerminalSizeDetectMillisec = 100        // 100 is the prod value

// Char is the smallest object.
// Future: a complex obj with foreground/bg colors, display attributes
type Char struct {
	runeVal rune
}

func (c Char) display() string {
	return string(c.runeVal)
}

/////////////////////////////////////////////////////////////////

type LineChars []Char

func (line LineChars) LineToStr() string {
	out := []string{}
	for _, Char := range line {
		out = append(out, Char.display())
	}
	return strings.Join(out, "")
}
func LineCharsFromStr(txt string) LineChars {
	Line := LineChars{}
	for _, runeVal := range txt {
		Line = append(Line, Char{runeVal: runeVal})
	}
	return Line
}

/////////////////////////////////////////////////////////////////

type MessageAndCharactersForWindowsUpdate struct {
	msg     string
	addLine LineChars
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
	lines             []LineChars
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

func (layers ScreenLayers) getLayer(layerIdWanted string) (ScreenLayer_CharMatrix, error) {
	for _, layer := range layers {
		if layerIdWanted == layer.layerId {
			return layer, nil
		}
	}
	return LayerEmtpyIfWeHaveErrors(), errors.New("unknown layerId")
}

type ScreenLayer_CharMatrix struct {
	xLeft   int
	yTop    int
	matrix  []ScreenColumn
	layerId string
}

func (sl ScreenLayer_CharMatrix) print(dataInputLineSeparator string) {
	fmt.Println("\nlayerId:", sl.layerId)
	matrixHeight := len(sl.matrix[0])
	for y := 0; y < matrixHeight; y++ {
		for _, column := range sl.matrix {
			fmt.Print(column[y].display())
		}
		fmt.Print(dataInputLineSeparator)
	}
}

func (layer ScreenLayer_CharMatrix) layerToTxt(lineSep string) string {
	yMax := len(layer.matrix[0])
	columns := layer.matrix

	output := []string{}
	for y := 0; y < yMax; y++ {
		for _, column := range columns {
			output = append(output, column[y].display())
		}
		if y < yMax-1 {
			output = append(output, lineSep)
		}
	}
	return strings.Join(output, "")
}

type ScreenColumn []Char

/////////////////////////////////////////////////////////////////

func UserInterfaceStart(ch_data_input chan MessageAndCharactersForWindowsUpdate, dataInputLineSeparator string) {
	userInterfaceInit()
	ch_user_input := make(chan string)
	go channel_read_user_input(ch_user_input)

	ch_terminal_size_change_detect := make(chan [2]int)
	go channel_read_terminal_size_change_detect(ch_terminal_size_change_detect)

	// windows is a read-only variable everywhere,
	windows := Windows{} // modified/updated ONLY here:
	go channel_read_dataInputInterpret(ch_data_input, &windows, dataInputLineSeparator)

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

func userInterfaceInit() {
	terminal_console_clear()
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func UserInterfaceExit() {
	terminal_console_character_show()
	terminal_console_input_buffering_enable()
}

////////////////////////////////////////////////////////////////////////////////////////////////

// TESTED
func LayerEmtpyIfWeHaveErrors() ScreenLayer_CharMatrix {
	return LayerCreate(0, 0, 1, 1, "e", "emptyLayer")
}

func LayerCreate(xLeft, yTop, width, height int, txtLayerDefault string, layerId string) ScreenLayer_CharMatrix {
	// fmt.Println("screen layer create:", xLeft, yTop, width, height)
	screenLayerNew := ScreenLayer_CharMatrix{xLeft: xLeft, yTop: yTop, layerId: layerId}
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
		textBlockVisible := []LineChars{LineChars{}}
		lineNumFirstVisible := len(win.lines) - win.height

		// this solution handles the too long lines, and flow the text into the next line
		for lineNumInWin, lineReceived := range win.lines {
			if lineNumInWin >= lineNumFirstVisible {

				if len(textBlockVisible) > 0 {
					// if we have a previous line, add a new empty one,
					// and we can fill it up with characters in the next for loop
					// this adding happens only from the 2nd received line.
					textBlockVisible = append(textBlockVisible, LineChars{})
				}

				textBlockLastLineId := len(textBlockVisible) - 1
				LineActual := textBlockVisible[textBlockLastLineId]

				for _, charNow := range lineReceived {
					LineActual = append(LineActual, charNow)
					// if LineActual is filled with chars maximally
					if len(LineActual) == len(screenNow.matrix) {
						textBlockVisible[textBlockLastLineId] = LineActual
						textBlockVisible = append(textBlockVisible, LineChars{})
						textBlockLastLineId = len(textBlockVisible) - 1
						LineActual = textBlockVisible[textBlockLastLineId]
					}
				}
				textBlockVisible[textBlockLastLineId] = LineActual
			}
		}

		for _, lineCharsVisible := range textBlockVisible {
			fmt.Println("DEBUG LINE:", lineCharsVisible)
		}
		// Load the theoretically visible text into the windows' column structure
		len_textBlockVisible := len(textBlockVisible)
		lineNumFirstVisibleInWindows := 0
		if len_textBlockVisible > win.height {

			// win.height +1 explanation:
			/* if we have 5 lines, and the window height is 4:
			1
			2 VISIBLE
			3 VISIBLE
			4 VISIBLE
			5 VISIBLE
			so if we have 5 lines, then -win.height+1 -> 5 - 4 + 1 = 2 the correct first visible value
			*/

			// last -1:  because the numbering is 0 based, not 1
			lineNumFirstVisibleInWindows = (len_textBlockVisible - win.height + 1) - 1
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

func channel_read_dataInputInterpret(ch_data_input chan MessageAndCharactersForWindowsUpdate, windows *Windows, dataInputLineSeparator string) {
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

func dataInputProcessLineByLine(dataInput MessageAndCharactersForWindowsUpdate, windows *Windows, dataInputLineSeparator string) string {
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
