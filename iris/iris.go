// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 100
var TimeIntervalTerminalSizeDetectMillisec = 100

func UserInterfaceStart(windows Windows) {

	///////////////////////////////////////////////////
	// keypress detection is based on this example:
	// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
	// thank you.
	ch_user_input := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
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
			time.Sleep(time.Millisecond * time.Duration(TimeIntervalTerminalSizeDetectMillisec))
		}
	}(ch_terminal_size_change_detect)

	screen_clear()

	for {
		fmt.Print(screen_cursor_pos_home())
		screenComposed := ScreensComposeToScreen(windows, []string{"Terminal", "Child"})
		fmt.Print(screenComposed.toString())

		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			fmt.Println("Keys pressed:", stdin)
			if stdin == "q" {
				action = "quit"
			}
			// vim navigation keys
			if strings.Contains("lhjk", stdin) {
				winActiveId := windows["prgState"]["winActiveId"]
				fmt.Println("win active id", winActiveId)
				if stdin == "l" {
					windows[winActiveId][KeyXleftCalculated] = "5"
					windows[winActiveId][KeyXrightCalculated] = "5"
				}
			}
		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			fmt.Println("terminal size change:", terminal_size_change)
		default: //               or not coming
			fmt.Println("No user input..")
		}
		if action == "quit" {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(TimeIntervalUserInterfaceRefreshTimeMillisec))

	}
}
