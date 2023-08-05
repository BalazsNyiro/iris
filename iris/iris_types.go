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
	"errors"
	"fmt"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 1000 // 10 is the prod value
var TimeIntervalTerminalSizeDetectMillisec = 100        // 100 is the prod value

// Char is the smallest object.
// Future: a complex obj with foreground/bg colors, display attributes
type Char struct {
	runeVal rune
}

func (c Char) display() string {
	return string(c.runeVal)
}

/////////////////////////////////////////////////////////////////

type LineChars []Char

func (line LineChars) LineToStr() string {
	out := []string{}
	for _, Char := range line {
		out = append(out, Char.display())
	}
	return strings.Join(out, "")
}
func LineCharsFromStr(txt string) LineChars {
	Line := LineChars{}
	for _, runeVal := range txt {
		Line = append(Line, Char{runeVal: runeVal})
	}
	return Line
}

/////////////////////////////////////////////////////////////////

type MessageAndCharactersForWindowsUpdate struct {
	msg     string
	addLine LineChars
}

// ///////////////////////////////////////////////////////////////
type Windows map[string]Window

func (wins Windows) printAll() {
	for _, win := range wins {
		win.print()
	}
}

// A window is a logical unit.
// It has settings, and content, but doesn't now
// anything how it will be rendered
type Window struct {
	id string

	// top-left coord: 0, 0 in the root terminal
	yTop              int
	xLeft             int
	width             int
	height            int
	lines             []LineChars
	backgroundDefault string
	winId             string

	// every decorator receives the decorator-section, when decorators are split with the decorator separator,
	// \n is the decorator-sections separator,
	// \t is the decorator elem separator
	decorators string // decoratorFun:borderSimple\tborder_left:|\tborder_right:|\ndecoratorFun:Other/key:val/key2:val2
}

func (w Window) print() {
	fmt.Println("winId:  ", w.winId)
	for _, line := range w.lines {
		fmt.Println("winLine:", line.LineToStr())
	}
}

/////////////////////////////////////////////////////////////////

type ScreenLayers []ScreenLayer_CharMatrix

func (layers ScreenLayers) getLayer(layerIdWanted string) (ScreenLayer_CharMatrix, error) {
	for _, layer := range layers {
		if layerIdWanted == layer.layerId {
			return layer, nil
		}
	}
	return LayerEmtpyIfWeHaveErrors(), errors.New("unknown layerId")
}

type ScreenLayer_CharMatrix struct {
	xLeft   int
	yTop    int
	matrix  []ScreenColumn
	layerId string
}

func (sl ScreenLayer_CharMatrix) print(dataInputLineSeparator string) {
	fmt.Println("\nlayerId:", sl.layerId)
	matrixHeight := len(sl.matrix[0])
	for y := 0; y < matrixHeight; y++ {
		for _, column := range sl.matrix {
			fmt.Print(column[y].display())
		}
		fmt.Print(dataInputLineSeparator)
	}
}

func (layer ScreenLayer_CharMatrix) layerToTxt(lineSep string) string {
	yMax := len(layer.matrix[0])
	columns := layer.matrix

	output := []string{}
	for y := 0; y < yMax; y++ {
		for _, column := range columns {
			output = append(output, column[y].display())
		}
		if y < yMax-1 {
			output = append(output, lineSep)
		}
	}
	return strings.Join(output, "")
}

type ScreenColumn []Char

/////////////////////////////////////////////////////////////////
