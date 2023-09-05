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

	text := TextFromStr("a\nbc\ndef")
	compare_int_pair(funName, len(text.Lines), 3, t)
	compare_int_pair(funName, len(text.Lines[0].Chars), 1, t)
	compare_int_pair(funName, len(text.Lines[1].Chars), 2, t)
	compare_int_pair(funName, len(text.Lines[2].Chars), 3, t)
	compare_int_pair(funName, text.NextLineNum, 3, t)

	compare_rune_pair(funName, text.Lines[2].Chars[2].runeVal, 'f', t)

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
