// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

/* docs:
https://en.wikipedia.org/wiki/ANSI_escape_code
https://stackoverflow.com/questions/37774983/clearing-the-screen-by-printing-a-character
*/

func screen_cursor_pos_home() string {
	return "\033[H"
}

func screen_clear() string {
	return "\033[2J"
}
