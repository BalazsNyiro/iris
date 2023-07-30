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
	"testing"
)

/*
	func Test_terminal_detect(t *testing.T) {
		widthStty, heightStty := TerminalDimensionsSttySize()
		fmt.Println("test stty:", widthStty, heightStty)
		widthSys, heightSys := TerminalDimensionsWithSyscall()
		fmt.Println("test syscall:", widthSys, heightSys)
	}

	func Test_document_create(t *testing.T) {
		rootObj := DocumentCreate("0", "50%", "50%", "40", "20")

}
*/

var newline = "\n"

func Test_layer_create(t *testing.T) {
	xLeft := 2
	yTop := 3
	width := 4
	height := 2
	txtlayerDefault := "L"

	layer := LayerCreate(xLeft, yTop, width, height, txtlayerDefault, "testLayerCreate")
	layerRendered := layer.layerToTxt(newline)
	fmt.Println(layerRendered)

	compare_str_pair("Test_layer_create 1", layerRendered, "LLLL"+newline+"LLLL", t)
}

func Test_layer_render(t *testing.T) {
	windows := Windows{} // modified/updated ONLY here:
	dataInputLineSeparator := "\n"
	data_input := MessageAndCharactersForWindowsUpdate{
		msg: `select:win:logs-left
						msgId:1
						set:xLeft:2
						set:yTop:3
						set:width:8
						set:height:6 `,
		addLine: LineCharsFromStr("testRender1")}
	dataInputProcessLineByLine(data_input, &windows, dataInputLineSeparator)

	data_input2 := MessageAndCharactersForWindowsUpdate{
		msg: `select:win:second
						msgId:2
						set:xLeft:1
						set:yTop:2
						set:width:4
						set:height:4 `,
		addLine: LineCharsFromStr("secondTextWithMoreLetters")}
	dataInputProcessLineByLine(data_input2, &windows, dataInputLineSeparator)

	// windows.printAll()
	terminalSizeActual := [2]int{9, 6}
	layers := LayersRenderFromWindows(windows, terminalSizeActual)
	fmt.Println("layers in test:", layers)

	for _, layer := range layers {
		layer.print(dataInputLineSeparator)
	}

	layerSecond, err := layers.getLayer("second")
	if err == nil {
		wanted := "thMo" + newline + "reLe" + newline + "tter" + newline + "srrr"
		compare_str_pair("Test layer SECOND, rendered txt", layerSecond.layerToTxt(newline), wanted, t)
	} else {
		fmt.Println("ERROR:", err)
	}

}

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}
