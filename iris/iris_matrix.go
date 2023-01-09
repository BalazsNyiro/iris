package iris

import "strings"

type MatrixCoord [2]int

type MatrixCharsRenderedWithFgBgSettings map[MatrixCoord]CharObj

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
			coordinate := MatrixCoord{x, y}
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
			coordinate := MatrixCoord{x, y}
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
			coordinate := MatrixCoord{x, y}
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
	for _, winName := range WinNamesSortByAttribute(windows, winNamesToComposite, KeyLayerNum, "number") {
		matrixActualWin := windows[winName].RenderToMatrixCharsOfWin(windowsChars, matrixFiller)

		winLocalXLeftCalculated := Str2Int(windows[winName][KeyXleftCalculated]) // read the values only once,
		winLocalYTopCalculated := Str2Int(windows[winName][KeyYtopCalculated])   // avoid to be changed in the for loop

		for yInWin := 0; yInWin < matrixActualWin.height; yInWin++ {
			for xInWin := 0; xInWin < matrixActualWin.width; xInWin++ {
				coordInWinLocal := MatrixCoord{xInWin, yInWin}
				coordInRootTerminal := MatrixCoord{winLocalXLeftCalculated + xInWin, winLocalYTopCalculated + yInWin}
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
