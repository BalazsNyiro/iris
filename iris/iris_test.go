package iris

// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

import (
	"fmt"
	"testing"
)

func Test_terminal_detect(t *testing.T) {
	widthStty, heightStty := TerminalDimensionsSttySize()
	fmt.Println("test stty:", widthStty, heightStty)
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("test syscall:", widthSys, heightSys)
}
