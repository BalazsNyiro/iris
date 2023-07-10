// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"fmt"
	iris "github.com/BalazsNyiro/iris/iris/TRASH_OLD_VERSION"
	"os"
	"strings"
)

var TimeIntervalUserInterfaceRefreshTimeMillisec = 10
var TimeIntervalTerminalSizeDetectMillisec = 100

type Windows map[string]Window
type Window struct {
	id string

	// top-left coord: 0, 0 in the root terminal
	top    int
	bottom int
	left   int
	right  int
	lines  []string
}

func UserInterfaceStart(ch_data_input chan string) {
	ui_init()
	ch_user_input := make(chan string)
	go channel_read_user_input(ch_user_input)

	ch_terminal_size_change_detect := make(chan [2]int)
	go channel_read_terminal_size_change_detect(ch_terminal_size_change_detect)

	// windows is a read-only variable everywhere,
	windows := Windows{} // modified/updated ONLY here:
	go data_input_interpret(ch_data_input, &windows)

	for {
		action := ""
		select { //                https://gobyexample.com/select
		case stdin, _ := <-ch_user_input: //  the message is coming...
			fmt.Println("Keys pressed:", stdin)
			action = action_of_user_input(stdin)

		case terminal_size_change, _ := <-ch_terminal_size_change_detect: //  the message is coming...
			fmt.Println("terminal size change:", terminal_size_change)
			// TODO: where do you store the new terminal size?

		default: //               or not coming
			_ = ""
		}

		if action == "quit" {
			UserInterfaceExit()
			break
		}

		TimeSleep(TimeIntervalUserInterfaceRefreshTimeMillisec)
	}
}

func data_input_interpret(ch_data_input chan string, windows *Windows) {

	for {
		select {
		case dataInput, _ := <-ch_data_input:
			// fmt.Println("data input:", dataInput)
			if strings.HasPrefix(dataInput, "select:win") {
				winId := select_win(dataInput, windows)
				if winId != "" {
					fmt.Println("after select:win, addSimpleText", (*windows)[winId].lines)

				}
			}
		default:
			_ = ""
		}
	}
}

/*
'add:simpleText:' is always the last added elem, everything after it is added automatically
into the lines
*/
func select_win(dataInput string, windows *Windows) string {
	winId := ""
	addSimpleTextDetectedLine := -1

	for lineNum, lineOrig := range strings.Split(dataInput, "\n") {
		line := strings.TrimSpace(lineOrig)
		fmt.Println("select_win, line:", line)
		elems := strings.Split(line, ":")

		// select:win:nameOfWin
		if elems[0] == "select" && elems[1] == "win" {
			if len(elems) == 3 {
				winId = strings.TrimSpace(elems[2])
				if _, exist := (*windows)[winId]; !exist {
					(*windows)[winId] = Window{}
				}
			}
		}

		if winId == "" {
			continue
		}
		if elems[0] == "add" && elems[1] == "simpleText" {
			addSimpleTextDetectedLine = lineNum
			win := (*windows)[winId]
			win.lines = append(win.lines, strings.SplitN(lineOrig, "add:simpleText:", 1)[1])
			(*windows)[winId] = win
			break
		}

	}
	if addSimpleTextDetectedLine > -1 {
		for lineNum, lineOrig := range strings.Split(dataInput, "\n") {
			if lineNum > addSimpleTextDetectedLine {
				win := (*windows)[winId]
				win.lines = append(win.lines, lineOrig)
				(*windows)[winId] = win
			}
		}
	}
	return winId
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
	terminal_console_clear()
	terminal_console_input_buffering_disable()
	terminal_console_character_hide()
}

func UserInterfaceExit() {
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

func channel_read_terminal_size_change_detect(ch chan [2]int) {
	widthSys, heightSys := 0, 0
	for {
		widthSysNow, heightSysNow := iris.TerminalDimensionsWithSyscall()
		if widthSysNow != widthSys || heightSysNow != heightSys {
			widthSys = widthSysNow
			heightSys = heightSysNow
			ch <- [2]int{widthSys, heightSys}
		}
		TimeSleep(TimeIntervalTerminalSizeDetectMillisec)
	}
}
