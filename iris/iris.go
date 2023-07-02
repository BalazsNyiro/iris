// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	"os"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 10
var TimeIntervalTerminalSizeDetectMillisec = 100

func UserInterfaceStart(windows Windows, windowsChars WindowsChars) {

	///////////////////////////////////////////////////
	// keypress detection is based on this example:
	// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
	// thank you.
	ch_user_input := make(chan string)
	go func(ch chan string) {
		terminal_console_disable_input_buffering()
		terminal_console_character_hide()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch_user_input)
	///////////////////////////////////////////////////

	ch_terminal_size_change_detect := make(chan string)
	go func(ch chan string) {
		for {
			// TODO: detect terminal size change here
			TimeSleep(TimeIntervalTerminalSizeDetectMillisec)
		}
	}(ch_terminal_size_change_detect)

	terminal_console_clear()

	matrixCharsComposedStr_prev := ""

	for {

		// the windows content can be updated from an outsider source without direct user input
		windows = WinCoordsCalculateUpdate(windows)

		matrixCharsComposed := MatrixCharsCompose(windows, windowsChars, []string{"Terminal", "Child"}, " ")

		matrixCharsComposedStr := matrixCharsComposed.toString()
		if matrixCharsComposedStr != matrixCharsComposedStr_prev {
			fmt.Print(terminal_console_cursor_pos_home())
			fmt.Print(matrixCharsComposedStr)
		}

		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			// fmt.Println("Keys pressed:", stdin)
			if stdin == "q" {
				action = "quit"
			}
			// vim navigation keys
			if strings.Contains("lhjk", stdin) {
				winActiveId := windows["prgState"]["winActiveId"]
				DebugInfoSave(windows)
				if stdin == "l" {
					windows[winActiveId][KeyXshift] = StrMath(windows[winActiveId][KeyXshift], "+", "1")
				}
				if stdin == "h" {
					windows[winActiveId][KeyXshift] = StrMath(windows[winActiveId][KeyXshift], "-", "1")
				}
				if stdin == "j" {
					windows[winActiveId][KeyYshift] = StrMath(windows[winActiveId][KeyYshift], "+", "1")
				}
				if stdin == "k" {
					windows[winActiveId][KeyYshift] = StrMath(windows[winActiveId][KeyYshift], "-", "1")
				}
			}
		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			fmt.Println("terminal size change:", terminal_size_change)
		default: //               or not coming
			_ = ""
		}
		if action == "quit" {
			terminal_console_character_show()
			break
		}
		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)

		matrixCharsComposedStr_prev = matrixCharsComposedStr
	}
}
