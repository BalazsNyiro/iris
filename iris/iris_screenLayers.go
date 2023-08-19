/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package iris

import "fmt"

// TESTED
func ScreenLayerEmtpyIfWeHaveErrors() ScreenLayer_CharMatrix {
	return ScreenLayerCreate(0, 0, 1, 1, "e", "emptyLayer")
}

// a screenLayer is a matrix of characters. A displayed unit.
// I always tried to pass a simple window instead of parameters here,
// but don't do that. you need to create ScreenLayers without a window
func ScreenLayerCreate(xLeft, yTop, width, height int, txtLayerDefault string, layerId string) ScreenLayer_CharMatrix {
	// fmt.Println("screen layer create:", xLeft, yTop, width, height)
	defaultRune := 'r'
	if len(txtLayerDefault) > 0 {
		defaultRune = rune(txtLayerDefault[0])
	}
	screenLayerNew := ScreenLayer_CharMatrix{xLeft: xLeft, yTop: yTop, layerId: layerId, defaultRune: defaultRune}
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

func ScreenLayersRenderFromWindows(windowsRO Windows, terminalSize [2]int) ScreenLayers {
	fmt.Println("terminal size:", terminalSize)
	layerBackground := ScreenLayerCreate(
		0, 0,
		terminalSize[0],
		terminalSize[1], ".", "layerBackground")

	layers := ScreenLayers{layerBackground}

	for _, win := range windowsRO {
		fmt.Println("render: winId >", win.winId, "< xLeft:", win.xLeft, "yTop:", win.yTop, "width:", win.width, "height:", win.height)
		screenNow := ScreenLayerCreate(
			win.xLeft, win.yTop,
			win.width, win.height, win.backgroundDefault, win.winId)

		// structure the character input into one COLUMN, (a visible structure)
		textBlockVisible := []LineChars{LineChars{}}
		lineNumFirstVisible := len(win.lines) - win.height

		// this solution handles the too long lines, and flow the text into the next line
		// FILL LINES ////////////////////////////////////////////////////////
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
		} // FILL LINES //////////////////////////////////////////////////////

		for _, lineCharsVisible := range textBlockVisible {
			fmt.Println("DEBUG LINE:", lineCharsVisible)
		}
		// Load the theoretically visible text into the windows' column structure
		len_textBlockVisible := len(textBlockVisible)
		lineNumFirstVisibleInWindows := 0
		if len_textBlockVisible > win.height {

			// win.height +1 explanation:
			/* if we have 5 lines, and the window height is 4:
			1 invisible line
			2 VISIBLE \
			3 VISIBLE .\ these lines are visible in the window
			4 VISIBLE ./
			5 VISIBLE /
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
		screenNow = ScreenLayerBorderRender(screenNow, win)
		layers = append(layers, screenNow)
	}
	return layers
}

func ScreenLayerBorderRender(screen ScreenLayer_CharMatrix, winParent Window) ScreenLayer_CharMatrix {
	borderSizeLeft := len(winParent.borderLeft)
	borderSizeTop := len(winParent.borderTop)
	borderSizeRight := len(winParent.borderRight)
	borderSizeBottom := len(winParent.borderBottom)

	widthOrig := len(screen.matrix)
	heightOrig := len(screen.matrix[0])

	xLeft := screen.xLeft - borderSizeLeft
	yTop := screen.yTop - borderSizeTop
	width := borderSizeLeft + widthOrig + borderSizeRight
	height := borderSizeTop + heightOrig + borderSizeBottom

	screenBordered := ScreenLayerCreate(xLeft, yTop, width, height, winParent.backgroundDefault, winParent.winId)
	// copy the data from the source screen to the target
	for x := 0; x < widthOrig; x++ {
		for y := 0; y < heightOrig; y++ {
			charInserted := screen.matrix[x][y]
			screenBordered.matrix[x+borderSizeLeft][y+borderSizeTop] = charInserted
		}
	}

	defaultBorderRune := '*'
	// DRAW LEFT BORDER - the string can be longer than 1 char!

	/*  ABccccDE  so a wide string can be a border, too
	    AB....DE
	    AB....DE
	    ABCCCCDE
	*/

	for y := 0; y < height; y++ {
		for x := 0; x < borderSizeLeft; x++ {
			runeInserted := defaultBorderRune
			necessaryBorderStringLen := x + 1 // if the borderString is too short, don't read
			if necessaryBorderStringLen <= len(winParent.borderLeft) {
				runeInserted = rune(winParent.borderLeft[x])
			}
			charInserted := Char{runeVal: runeInserted}
			screenBordered.matrix[x][y] = charInserted
		}

		// the right border is reversed: read the string from right to left
		positionRightInLayer := width - 1                       // the last column
		positionInBorderRight := len(winParent.borderRight) - 1 // the last char in the border string

		for positionInBorderRight >= 0 {
			runeInserted := rune(winParent.borderRight[positionInBorderRight])
			charInserted := Char{runeVal: runeInserted}
			screenBordered.matrix[positionRightInLayer][y] = charInserted
			positionInBorderRight -= 1
			positionRightInLayer -= 1
		}

	}

	// example borderTop:    "-v", // - is outer, - is inner
	for y := 0; y < borderSizeTop; y++ {
		for x := 0; x < width; x++ {
			runeInserted := rune(winParent.borderTop[y])
			charInserted := Char{runeVal: runeInserted}
			screenBordered.matrix[x][y] = charInserted
		}
	}

	for y := 0; y < borderSizeBottom; y++ {
		for x := 0; x < width; x++ {
			runeInserted := rune(winParent.borderBottom[y])
			charInserted := Char{runeVal: runeInserted}
			screenBordered.matrix[x][y] = charInserted
		}
	}

	return screenBordered
}

func ScreenLayersDisplayAll(layers ScreenLayers, newlineSeparator string, loopCounter int) {
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

	screenMerged := ScreenLayerCreate(
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
