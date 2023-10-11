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

// Char: basic Character object
type Char struct {
	runeVal rune
}

func (c Char) display() string {
	return string(c.runeVal)
}

/////////////////////////////////////////////////////////////////

// LineChars: Horizontal character representation.
// LineNum can be important, because to save memory, early lines can be removed
// from the TextBlock - and in this situation you cannot know the original line number of the inserted lines.
// the first line's LineNum is 0.
// if LineNum == -1, it means that the real line number is unknown at the moment of LineNum creation
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
	// if lineNum == -1 it means the lineNum is unknown
	line := LineChars{LineNum: lineNum}
	for _, runeVal := range text {
		line.Chars = append(line.Chars, Char{runeVal: runeVal})
	}
	return line
}

/////////////////////////////////////////////////////////////////

// MessageAndCharactersForTextBlockUpdate is a channel transported message for the TextBlocks.
// it can contain statements/commands to manage TextBlocks,
// and it has a character container 'addLine' if you want to pass characters to the TextBlock
type MessageAndCharactersForTextBlockUpdate struct {
	msg     string
	addLine LineChars
}

// ///////////////////////////////////////////////////////////////

// TextBlock: the most important data structure, a logical unit of lines.
// one TextBlock can be displayed in more Windows, and in these Windows you can see
// different parts of the same TextBlock
type TextBlock struct {
	Lines       []LineChars
	NextLineNum int
}

// TextBlocks: key based TextBlock storage
type TextBlocks map[string]TextBlock

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
