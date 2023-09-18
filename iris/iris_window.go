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

func text_line_last_displayed_segment(lenText int, winWidth int) (int, int) {
	// this fun assumes that the text length is minimum 1 char.
	// in that case, the start=0, end=0 is the minimum possible range, that could be given back
	// if the text length is less, so there is no text, return with -1, -1 by definition

	if lenText <= 0 {
		return -1, -1
	}
	if lenText == 1 {
		return 0, 0
	}
	if lenText > 1 && lenText <= winWidth {
		return 0, lenText - 1
	}

	if lenText > winWidth {
		numOfCharsInLastNotFullLine := lenText % winWidth

		// if winWidth == 4, lenText == 4, 8, 12, 16 -
		// so there is no empty line, everything is filled
		if numOfCharsInLastNotFullLine == 0 {
			indexStart := (lenText - 1) - (winWidth - 1)
			indexEnd := lenText - 1
			return indexStart, indexEnd
		} else {
			/*
					2 ....      ABCD  <- totally filled FIRST window row
					3 ....	    EFGH  <- totally filled       window row
					4 ....	    IJKL  <- totally filled       window row
					5 ....	    MN    <- partial              window row

					so in this case, the last line has partial chars, the window line is not filled.
				    numOfCharsInLastNotFullLine == 2

					in this case, the end is the last char, the first is the beginning
					of the reminder/numOfCharsInLastNotFullLine
			*/
			indexStart := (lenText - 1) - (numOfCharsInLastNotFullLine - 1)
			indexEnd := lenText - 1
			return indexStart, indexEnd
		}

	}

	// in worst case, if nothing is matching, return with this :-)
	return -2, -2
}

func (w Window) render() {
	//matrix := MatrixNew(w.width, w.height)

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

	one text line can fill more window row.

	so the textBlock.Lines needs to be splitted up smaller
	sections, and in reverse orders we need to upload the lines.

	1: Loop over the lines in reversed order, from last to the first ones
	2:
	*/

	fmt.Println("RENDER")
	fmt.Println("win render, textBlock lineNum", len(w.textBlockPtr.Lines))
	for i := len(w.textBlockPtr.Lines) - 1; i >= 0; i-- {
		line := w.textBlockPtr.Lines[i]
		lineStr := line.LineToStr()
		fmt.Println("winLine:", i, lineStr)

		widthSmallerAndSmaller := len(lineStr)
		for true {
			lastSegmentIndexStart, lastSegmentIndexEnd := text_line_last_displayed_segment(widthSmallerAndSmaller, w.width)
			if lastSegmentIndexStart < 0 {
				break
			}
			fmt.Println("lastSegment start:", lastSegmentIndexStart, " end:", lastSegmentIndexEnd)
			widthSmallerAndSmaller = lastSegmentIndexStart
			// the lastSegmentIndexStart is the last segment's start point.
			// because the indexes are zero based, the new width (which is 1 based)
			// is equal width the last 0 based start
		}
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
