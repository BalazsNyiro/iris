/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

// A window is a logical unit.
// It has settings, and content, but doesn't now
// anything how it will be rendered
package iris

import "fmt"

type ColumnChars []Char
type MatrixChars []ColumnChars // 2 dimensional Char structure

type Window struct {
	id string

	// top-left coord: 0, 0 in the root terminal
	yTop              int
	xLeft             int
	width             int
	height            int
	textBlockPtr      *TextBlock
	backgroundDefault string
	winIdNum          int
	winHumanName      string
}

func (w Window) print() {
	fmt.Println("winId:  ", w.winIdNum)

	fmt.Println("textBlock lineNum", len(w.textBlockPtr.Lines))
	for _, line := range w.textBlockPtr.Lines {
		fmt.Println("winLine:", line.LineToStr())
	}
}
func (w Window) render() {
	// TODO: matrix := Matrix_empty()

	/* How can we render? :-) This is a game!

	an example textBlock.Lines (2 lines:)
	abcdefghijklm
	ABCDEFGHIJKLMN

	How can we render this? the window width = 4, height = 5

	A window map:
	  1234   the last line needs to be broken to win.width sections:
	1 ....
	2 ....      ABCD  <- totally filled FIRST window row
	3 ....	    EFGH  <- totally filled       window row
	4 ....	    IJKL  <- totally filled       window row
	5 ....	    MN    <- partial              window row

	one text line can fille more window row.

	so the textBlock.Lines needs to be splitted up smaller
	sections, and in reverse orders we need to upload the lines.
	*/

	fmt.Println("RENDER")
	fmt.Println("textBlock lineNum", len(w.textBlockPtr.Lines))
	for i := len(w.textBlockPtr.Lines) - 1; i >= 0; i-- {
		line := w.textBlockPtr.Lines[i]
		fmt.Println("winLine:", i, line.LineToStr())
	}
}

func window_row_start_end_positions(lengtOfTextBlockLine int, windowWidth int) [][2]int {

	positionValues := [][2]int{}
	lengthActual := lengtOfTextBlockLine

	for true {
		positionEnd := lengthActual - 1
		positionStart := positionEnd // a default value only, overwritten in if{}
		partialNotFullWindowRow_text_length := lengthActual % windowWidth

		if partialNotFullWindowRow_text_length > 0 {
			positionStart = positionEnd - partialNotFullWindowRow_text_length + 1
		} else {
			positionStart = positionEnd - windowWidth + 1
		}
		/*Concrete example: position 13 and position 14 are the partials.
		so we need to select (13, 14) first. but the length will be
		decreased with 2:  14-13+1 = 2
		*/
		lengthActual = lengthActual - (positionEnd - positionStart + 1)

		positionValues = append(positionValues, [2]int{positionStart, positionEnd})
		if positionStart == 0 {
			break
		}
	}

	return positionValues
}
