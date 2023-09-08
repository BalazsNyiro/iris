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
	"testing"
)

var newline = "\n"

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

func Test_windows_render(t *testing.T) {
	//funName := " Test_windows_display"
	textBlock := TextBlockFromStr("1a\n2bc\n3def\n4ghijklmno")
	winId := 0

	win := Window{yTop: 0, xLeft: 0, width: 4, height: 3, textBlockPtr: &textBlock, backgroundDefault: "b", winIdNum: winId}
	win.print()

	win.render()

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
