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

	line := LineCharsFromStr("abcd", 0)

	// 4 Chars are inserted into the line
	compare_int_pair("test_lines_create", len(line.Chars), 4, t)
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
