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
)

func channelRead_dataInputInterpret(ch_data_input chan MessageAndCharactersForTextBlockUpdate, textBlocks *TextBlocks, dataInputLineSeparator string) {
	for {
		select {
		case dataInput, _ := <-ch_data_input:
			fmt.Println("\ndata input:", dataInput)
			if strings.HasPrefix(strings.TrimSpace(dataInput.msg), "select:win") {
				dataInputProcessLineByLine(dataInput, textBlocks, dataInputLineSeparator)
			}
		default:
			_ = ""
		}
	}
}

func dataInputProcessLineByLine(dataInput MessageAndCharactersForTextBlockUpdate, textBlocks *TextBlocks, dataInputLineSeparator string) string {
	winId := ""
	fieldSeparator := ":"

	/*
				when an attribute is set, the line is split by separators,
			    to get the path what we need to do.
				For example:
				set:borderBottom:=

				it means: set borderBottom =
				so set the borderBottom value to =

				but what if we want to use the possible separator char, too?
				it can be escaped, for example, but it is ugly.

				So after a 'set:borderBottom:' the program knows that every char is accepted,
		        There is NO MORE SEPARATOR CHAR splitting.

	*/

	for _, lineOrig := range strings.Split(dataInput.msg, dataInputLineSeparator) {
		line := strings.TrimSpace(lineOrig)
		elems := strings.Split(line, fieldSeparator)

		// select:win:nameOfWin
		if elems[0] == "select" && elems[1] == "textBlock" {
			textBlockNameSelected := strings.TrimSpace(elems[2])
			if _, exist := (*textBlocks)[textBlockNameSelected]; !exist {
				(*textBlocks)[textBlockNameSelected] = TextBlock{}
			}
		}

	}

	return winId
}
