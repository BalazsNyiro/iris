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

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}
