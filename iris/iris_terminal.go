// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import "os/exec"

/* docs:
https://en.wikipedia.org/wiki/ANSI_escape_code
https://stackoverflow.com/questions/37774983/clearing-the-screen-by-printing-a-character
*/

func terminal_console_cursor_pos_home() string { // basic fun
	return "\033[H"
}

func terminal_console_clear() string { // basic fun
	return "\033[2J"
}

// unbuffered input manager in go:
// https://stackoverflow.com/questions/48831750/unbuffered-input-manager-in-go
func terminal_console_input_buffering_disable() { // basic fun
	exec.Command("stty", "--file=/dev/tty", "cbreak", "min", "1").Run()
}

func terminal_console_input_buffering_enable() { // basic fun
	exec.Command("stty", "--file=/dev/tty", "-cbreak", "min", "1").Run()
}

func terminal_console_character_show() { // basic fun
	exec.Command("stty", "--file=/dev/tty", "echo").Run()
}

// do not display entered characters in the console
func terminal_console_character_hide() { // basic fun
	exec.Command("stty", "--file=/dev/tty", "-echo").Run()
}
