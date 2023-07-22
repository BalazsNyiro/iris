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

	layer := LayerCreate(xLeft, yTop, width, height, txtlayerDefault)
	layerRendered := layer.layerToTxt("\n")
	fmt.Println(layerRendered)

	compare_str_pair("Test_layer_create 1", layerRendered, "LLLL\nLLLL", t)
}

func Test_layer_render(t *testing.T) {
	windows := Windows{} // modified/updated ONLY here:

	data_input := MessageAndCharacters{
		msg: `select:win:logs-left
						msgId:1
						set:top:3
						set:bottom:6
						set:left:2
						set:right:8`,
		addLine: LineFromStr("testRender1")}
	dataInputProcessLineByLine(data_input, &windows, "\n")

	data_input2 := MessageAndCharacters{
		msg: `select:win:second
						msgId:2
						set:top:2
						set:bottom:4
						set:left:1
						set:right:4`,
		addLine: LineFromStr("secondTxt")}
	dataInputProcessLineByLine(data_input2, &windows, "\n")

	windows.printAll()
}

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}
