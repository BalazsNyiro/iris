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
	"fmt"
	"strings"
	"testing"
)

func Test_lines_create(t *testing.T) {
	funName := "test_lines_create"

	line := LineCharsFromStr("abcd", 0)
	// 4 Chars are inserted into the line
	compare_int_pair(funName, len(line.Chars), 4, t)
	compare_rune_pair(funName, line.Chars[3].runeVal, 'd', t)
}

func Test_text_create(t *testing.T) {
	funName := "Test_text_create"

	text := TextBlockFromStr("a\nbc\ndef")
	compare_int_pair(funName, len(text.Lines), 3, t)
	compare_int_pair(funName, len(text.Lines[0].Chars), 1, t)
	compare_int_pair(funName, len(text.Lines[1].Chars), 2, t)
	compare_int_pair(funName, len(text.Lines[2].Chars), 3, t)
	compare_int_pair(funName, text.NextLineNum, 3, t)

	compare_rune_pair(funName, text.Lines[2].Chars[2].runeVal, 'f', t)

}

func Test_text_append(t *testing.T) {
	funName := "Test_text_append"
	textBlock := TextBlockFromStr("a\nbc\ndef")
	TextAppendIntoLastLine("ghi", &textBlock)

	compare_int_pair(funName, len(textBlock.Lines), 3, t)
	compare_int_pair(funName, len(textBlock.Lines[2].Chars), 6, t)

	TextAppendIntoNewNextLine("jklmn", &textBlock)
	compare_int_pair(funName, len(textBlock.Lines), 4, t)
	compare_int_pair(funName, len(textBlock.Lines[3].Chars), 5, t)

}

func Test_text_line_split_into_window_width_get_indexes(t *testing.T) {
	/* split a long text into a smaller window - find the last index */
	indexStart, indexEnd := text_line_last_displayed_segment__startIncluded_endIncluded(0, 4)
	compare_int_tuples("segment 1", []int{indexStart, indexEnd}, []int{-1, -1}, t)

	indexStart, indexEnd = text_line_last_displayed_segment__startIncluded_endIncluded(4, 4)
	compare_int_tuples("segment 2", []int{indexStart, indexEnd}, []int{0, 3}, t)

	indexStart, indexEnd = text_line_last_displayed_segment__startIncluded_endIncluded(8, 4)
	compare_int_tuples("segment 2", []int{indexStart, indexEnd}, []int{4, 7}, t)

	indexStart, indexEnd = text_line_last_displayed_segment__startIncluded_endIncluded(14, 4)
	compare_int_tuples("segment 2", []int{indexStart, indexEnd}, []int{12, 13}, t)
}

func Test_windows_render(t *testing.T) {
	//funName := " Test_windows_display"
	textBlock := TextBlockFromStr("1a\n2bc\n3def\n4ghijklmno")

	winId := 0
	win := Window{yTop: 0, xLeft: 0, width: 4, height: 5,
		textBlockPtr: &textBlock, backgroundDefault: ".", winIdNum: winId}
	win.print()

	matrixOfWin := win.matrixRender()

	fmt.Println("MATRIX DISPLAY AFTER WIN RENDER")
	matrixOfWin.DisplayInConsoleToDebugOrAnalyse()

	wantedRenderedOutput := `2bc.
	                         3def
	                         4ghi
	                         jklm
	                         no.. `
	compare_str_block("windowsRender1",
		matrixOfWin.Lines_representation_to_test_comparison(), wantedRenderedOutput, t)
}

func Test_root_matrix_and_windows_merged(t *testing.T) {

	// create a root matrix
	// create 2 windows.
	// matrixRender 2 matrixes from the window
	// merge everything
	// display it

	matrixRoot := MatrixNew(10, 8, "#")

	textBlock1 := TextBlockFromStr("1")
	textBlock2 := TextBlockFromStr("2")

	winId := 0
	win0 := Window{xLeft: 1, yTop: 1, width: 5, height: 4, textBlockPtr: &textBlock1, backgroundDefault: "?", winIdNum: winId}

	winId = 1
	win1 := Window{xLeft: 4, yTop: 3, width: 4, height: 5, textBlockPtr: &textBlock2, backgroundDefault: "-", winIdNum: winId}

	matrixWin0 := win0.matrixRender()
	matrixWin1 := win1.matrixRender()

	matrixMerged0 := MatrixAdd(matrixRoot, matrixWin0, win0.xLeft, win0.yTop)
	matrixMerged1 := MatrixAdd(matrixMerged0, matrixWin1, win1.xLeft, win1.yTop)

	matrixMerged1.DisplayInConsoleToDebugOrAnalyse()

	wantedRenderedOutput := `##########
						   	 #?????####
						  	 #?????####
						  	 #???----##
						  	 #1??----##
						 	 ####----##
						 	 ####----##
						 	 ####2---## `

	compare_str_block("matrixRenderTest1",
		matrixMerged1.Lines_representation_to_test_comparison(), wantedRenderedOutput, t)
}

func compare_str_block(callerInfo string, received []string, wanted string, t *testing.T) {
	// in the tests,  I use \n as a line separator, in wanted
	// detect lines in wanted, then compare lines
	for lineNum, lineWanted := range strings.Split(wanted, "\n") {
		lineWanted = strings.TrimSpace(lineWanted)
		lineReceived := received[lineNum]
		compare_str_pair(callerInfo, lineReceived, lineWanted, t)
	}
}

func Test_window_row_start_end_positions(t *testing.T) {
	windowWidth := 4
	windowRowPositions := window_row_start_end_positions(14, windowWidth)
	/*  good answers:
	position: [12 13]
	position: [8 11]
	position: [4 7]
	position: [0 3]
	*/

	compare_int_pair("rowStartEndPos1", windowRowPositions[0][0], 12, t)
	compare_int_pair("rowStartEndPos2", windowRowPositions[0][1], 13, t)
	compare_int_pair("rowStartEndPos3", windowRowPositions[1][0], 8, t)
	compare_int_pair("rowStartEndPos4", windowRowPositions[1][1], 11, t)
	compare_int_pair("rowStartEndPos5", windowRowPositions[2][0], 4, t)
	compare_int_pair("rowStartEndPos6", windowRowPositions[2][1], 7, t)
	compare_int_pair("rowStartEndPos7", windowRowPositions[3][0], 0, t)
	compare_int_pair("rowStartEndPos8", windowRowPositions[3][1], 3, t)
}
func compare_int_tuples(callerInfo string, received, wanted []int, t *testing.T) {
	if len(received) != len(wanted) {
		t.Fatalf("\nErr: %s received int tuple: %v\n  wanted int tuple: %v, error, different length, not comparable", callerInfo, received, wanted)
	}

	for i := 0; i < len(received); i++ {
		if received[i] != wanted[i] {
			t.Fatalf("\nErr: %s (id %v) received int in tuple: %v\n  wanted int in tuple: %v, error", callerInfo, i, received, wanted)
		}

	}
}

func compare_int_pair(callerInfo string, received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received int: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

func compare_rune_pair(callerInfo string, received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received rune = %v, wanted %v, error", callerInfo, received, wanted)
	}
}

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}
