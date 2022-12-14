// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import "os/exec"

/* docs:
https://en.wikipedia.org/wiki/ANSI_escape_code
https://stackoverflow.com/questions/37774983/clearing-the-screen-by-printing-a-character
*/

func terminal_console_cursor_pos_home() string {
	return "\033[H"
}

func terminal_console_clear() string {
	return "\033[2J"
}

func terminal_console_disable_input_buffering() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}

// do not display entered characters in the console
func terminal_console_character_hide() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}
