// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	"os"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 10
var TimeIntervalTerminalSizeDetectMillisec = 100

func UserInterfaceStart() {
	ui_init()
	ch_user_input := make(chan string)
	go channel_read_user_input(ch_user_input)

	ch_terminal_size_change_detect := make(chan string)
	go channel_read_terminal_size_change_detect(ch_terminal_size_change_detect)
	terminal_console_clear()

	for {
		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			fmt.Println("Keys pressed:", stdin)
			action = action_of_user_input(stdin)

		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			fmt.Println("terminal size change:", terminal_size_change)

		default: //               or not coming
			_ = ""
		}

		if action == "quit" {
			ui_exit()
			break
		}

		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)
	}
}

func action_of_user_input(stdin string) string {
	action := ""
	if stdin == "q" {
		action = "quit"
	}

	if stdin == "l" {
	}
	if stdin == "h" {
	}
	if stdin == "j" {
	}
	if stdin == "k" {
	}
	return action
}

func ui_init() {
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func ui_exit() {
	terminal_console_character_show()
	terminal_console_input_buffering_enable()
}

// /////////////////////////////////////////////////
// keypress detection is based on this example:
// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
// thank you.
func channel_read_user_input(ch chan string) {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		ch <- string(b)
	}
} ///////////////////////////////////////////////////

func channel_read_terminal_size_change_detect(ch chan string) {
	for {
		// TODO: detect terminal size change here
		TimeSleep(TimeIntervalTerminalSizeDetectMillisec)
	}
}
