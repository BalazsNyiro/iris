/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package iris

import "strings"

var TextProcessingNewlineSeparator = "\n"

type Char struct {
	runeVal rune
}

func (c Char) display() string {
	return string(c.runeVal)
}

/////////////////////////////////////////////////////////////////

type LineChars struct {
	Chars   []Char
	LineNum int
}

func (line LineChars) LineToStr() string {
	out := []string{}
	for _, Char := range line.Chars {
		out = append(out, Char.display())
	}
	return strings.Join(out, "")
}
func LineCharsFromStr(text string, lineNum int) LineChars {
	line := LineChars{LineNum: lineNum}
	for _, runeVal := range text {
		line.Chars = append(line.Chars, Char{runeVal: runeVal})
	}
	return line
}

// maybe it has to be an object, with attributes
type TextBlock struct {
	Lines       []LineChars
	NextLineNum int
}

func TextBlockFromStr(txt string) TextBlock {
	text := TextBlock{}
	for _, line := range strings.Split(txt, TextProcessingNewlineSeparator) {
		text.Lines = append(text.Lines, LineCharsFromStr(line, text.NextLineNum))
		text.NextLineNum += 1
	}
	return text
}

func TextAppendIntoLastLine(txtNew string, textBlockPtr *TextBlock) {
	nextLineNum := len(textBlockPtr.Lines)

	firstNewline := true
	for _, line := range strings.Split(txtNew, TextProcessingNewlineSeparator) {
		lineCharsNew := LineCharsFromStr(line, nextLineNum)
		if firstNewline {
			lineNumLastId := len(textBlockPtr.Lines) - 1
			for _, char := range lineCharsNew.Chars {
				textBlockPtr.Lines[lineNumLastId].Chars =
					append(textBlockPtr.Lines[lineNumLastId].Chars, char)
			}
		} else {
			textBlockPtr.Lines = append(textBlockPtr.Lines, lineCharsNew)
			textBlockPtr.NextLineNum += 1
			firstNewline = false
		}
	}
}

func TextAppendIntoNewNextLine(txtNew string, textBlockPtr *TextBlock) {
	nextLineNum := len(textBlockPtr.Lines)

	for _, line := range strings.Split(txtNew, TextProcessingNewlineSeparator) {
		textBlockPtr.Lines = append(textBlockPtr.Lines, LineCharsFromStr(line, nextLineNum))
		textBlockPtr.NextLineNum += 1
	}
}
