// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

type Color struct {
    // https://en.wikipedia.org/wiki/ANSI_escape_code
    Mode string // RGB6: calculate the result from rgb values
    // GRAY24: use a gray color code
    // NAMED16: first 16 color code, named values

    Red   int    // 0-5 step
    Green int    // 0-5 step
    Blue  int    // 0-5 step
    Gray  int    // 0-23 steps (color code range: 232 - 255)
    Name  string // named color value?
}

type CharObj struct {
    ColorFg Color
    ColorBg Color
    CharVal rune
}
