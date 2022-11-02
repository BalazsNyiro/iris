// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"strconv"
	"strings"
)

/*
	all DOM objects are containers
    An object can have childrens (other objects)
    An object can be rendered on the screen.
	Document: the Root container elem.

	I follow a Document Object model to represent a CLI gui
	I will use HTML to describe the structure

*/
type DomObj struct {
	/*
		in a gui we can use 'px'. here we use: char as base width unit

		Area: A DomObj has an area in the screen.
		=========================================
		All children is created from one coordinate point:
		        the object tries to form the wanted fix sizes (fix width or height)
		        or a relative size: (height_available-1) for example

	*/
	Attr     map[string]string
	Children []DomObj
	Parent   *DomObj
}

func ObjNew(id, width, height, bornX, bornY string, parentPointer *DomObj) DomObj {
	attr := map[string]string{
		"id": id, // all id is unique. Equivalent with dom Id

		"width":             width,             // value examples: 20, 20%,
		"height":            height,            // or it can be empty: ''
		"bornX":             bornX,             // default unit: character.
		"bornY":             bornY,             // 20 means: 20 character
		"bornHorizontalWay": "FromLeftToRight", // or FromRightToLeft
		"bornVerticalWay":   "FromTopToDown",   // or FromDownToTop
	}
	return DomObj{Attr: attr, Parent: parentPointer}
}
func (obj DomObj) WidthCalculated() int {
	return valueStringToNumber(obj.Attr["width"], obj.Parent, "width")
}

func valueStringToNumber(valStr string, parentPtr *DomObj, parentAttrKey string) int {
	valStr = strings.TrimSpace(valStr)
	valCalculated := 0
	if strings.Contains(valStr, "%") {
		valStr = strings.Replace(valStr, "%", "", -1)
		val, err := strconv.Atoi(valStr)
		if err == nil && parentPtr != nil {
			valCalculated = valueStringToNumber(
				parentPtr.Attr[parentAttrKey],
				parentPtr.Parent,
				parentAttrKey) * val / 100
		}
	} else {
		val, err := strconv.Atoi(valStr)
		if err == nil {
			valCalculated = val
		}
	}
	return valCalculated
}

func DocumentCreate(id, width, height, terminalWidth, terminalHeight string) DomObj {
	parentTerminal := ObjNew("terminal", terminalWidth, terminalHeight, "", "", nil)
	root := ObjNew("documentRoot", width, height, "0", "0", &parentTerminal)
	return root
}
