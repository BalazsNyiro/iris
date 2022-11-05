// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

func DocumentCreate(id, width, height, terminalWidth, terminalHeight string) DomObj {
	parentTerminal := ObjNew("terminal", terminalWidth, terminalHeight, "", "", "T", nil)
	root := ObjNew("documentRoot", width, height, "0", "0", "R", &parentTerminal)
	return root
}
