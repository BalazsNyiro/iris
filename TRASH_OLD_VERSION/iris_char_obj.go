// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

type CharObj struct {
	ColorFg Color
	ColorBg Color
	CharVal rune
}

func (c CharObj) render() string {
	// FIXME: COLOR prefixes
	// FIXME: COLOR postfixes, cleanings?
	return string(c.CharVal)
}

func (c CharObj) charVal() string {
	return string(c.CharVal)
}

func CharObjNew(oneRune rune) CharObj {
	return CharObjNewTotal(oneRune, "black", "white")
}

func CharObjNewTotal(oneRune rune, bgColorName, fgColorName string) CharObj {
	return CharObj{CharVal: oneRune, ColorBg: Color{Name: fgColorName}, ColorFg: Color{Name: bgColorName}}
}
