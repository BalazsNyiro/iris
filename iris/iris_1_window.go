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
	lines             []LineChars
	backgroundDefault string
	winId             string
}

func (w Window) print() {
	fmt.Println("winId:  ", w.winId)
	for _, line := range w.lines {
		fmt.Println("winLine:", line.LineToStr())
	}
}
