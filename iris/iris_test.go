package iris

import (
	"fmt"
	"testing"
)

// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

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

func Test_layer_create(t *testing.T) {
	xLeft := 2
	yTop := 3
	width := 4
	height := 2
	txtlayerDefault := "L"

	layer := LayerCreate(xLeft, yTop, width, height, txtlayerDefault, "testLayerCreate")
	layerRendered := layer.layerToTxt("\n")
	fmt.Println(layerRendered)

	compare_str_pair("Test_layer_create 1", layerRendered, "LLLL\nLLLL", t)
}

func Test_layer_render(t *testing.T) {
	windows := Windows{} // modified/updated ONLY here:
	dataInputLineSeparator := "\n"
	data_input := MessageAndCharactersForWindowsUpdateForWindowsUpdateForWindowsUpdate{
		msg: `select:win:logs-left
						msgId:1
						set:xLeft:2
						set:yTop:3
						set:width:8
						set:height:6 `,
		addLine: LineCharsFromStr("testRender1")}
	dataInputProcessLineByLine(data_input, &windows, dataInputLineSeparator)

	data_input2 := MessageAndCharactersForWindowsUpdateForWindowsUpdate{
		msg: `select:win:second
						msgId:2
						set:xLeft:1
						set:yTop:2
						set:width:4
						set:height:4 `,
		addLine: LineCharsFromStr("secondText")}
	dataInputProcessLineByLine(data_input2, &windows, dataInputLineSeparator)

	windows.printAll()
	terminalSizeActual := [2]int{9, 6}
	layers := LayersRenderFromWindows(windows, terminalSizeActual)
	fmt.Println("layers in test:", layers)

	for _, layer := range layers {
		layer.print(dataInputLineSeparator)
	}
}

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}
