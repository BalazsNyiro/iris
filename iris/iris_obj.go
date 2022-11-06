package iris

import (
	"strconv"
	"strings"
)

var FromLeftToRight = "FromLeftToRight"
var FromTopToDown = "FromTopToDown"

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

// columns of strings!
type RenderedColumn []string
type RenderedScreen []RenderedColumn // there are vertical columns next to each other.

// wrapper - with default settings
func (obj DomObj) Render() RenderedScreen {
	return obj.RenderScreenMatrix(true)
}

// YOU CAN TEST THIS
// every object knows about his own rendered output
func (obj DomObj) RenderScreenMatrix(useBasicFiller bool) RenderedScreen {

	objWidth := obj.attribCalculated("width")
	objHeight := obj.attribCalculated("height")

	if objWidth < 1 || objHeight < 1 {
		screenDisplayed := RenderedScreen{}
		return screenDisplayed
	}

	screenDisplayed := render_screen_columns_create(objWidth)

	for x := 0; x < objWidth; x++ {
		for y := 0; y < objHeight; y++ {
			if useBasicFiller {
				screenDisplayed[x] = append(screenDisplayed[x], obj.Attr["basicBackgroundFiller"])
			}
		}
	}
	return screenDisplayed
}

func render_screen_columns_create(columnNum int) RenderedScreen {
	screenDisplayed := RenderedScreen{}
	for x := 0; x < columnNum; x++ {
		screenDisplayed = append(screenDisplayed, RenderedColumn{})
	}
	return screenDisplayed
}

func ScreenToTxt(screen RenderedScreen) string {
	columnsNum := len(screen)
	if columnsNum < 1 {
		return ""
	}
	rowsNum := len(screen[0])
	if rowsNum < 1 {
		return ""
	}

	out := []string{}
	for y := 0; y < rowsNum; y++ {
		for x := 0; x < columnsNum; x++ {
			out = append(out, screen[x][y])
		}
		if y < rowsNum-1 {
			out = append(out, "\n")
		}
	}
	return strings.Join(out[:], "")
}

func ObjNew(id, width, height, bornX, bornY, basicBgFiller string, parentPointer *DomObj) DomObj {
	attr := map[string]string{
		"id": id, // all id is unique. Equivalent with dom Id

		"width":             width,           // value examples: 20, 20%,
		"height":            height,          // or it can be empty: ''
		"bornX":             bornX,           // default unit: character.
		"bornY":             bornY,           // 20 means: 20 character
		"bornHorizontalWay": FromLeftToRight, // or FromRightToLeft
		"bornVerticalWay":   FromTopToDown,   // or FromDownToTop

		"basicBackgroundFiller": basicBgFiller, // this string represents the object in debug/dev mode
		// the displayed filler is one char wide. But: it can be longer than once char, with color codes for example

	}
	return DomObj{Attr: attr, Parent: parentPointer}
}

// Tested
func (obj DomObj) attribCalculated(key string) int {
	return valueStringToNumber(obj.Attr[key], obj.Parent, key)
}

// Tested
func valueStringToNumber(valStr string, baseObjectPointer *DomObj, baseObjAttrName string) int {
	valStr = strings.TrimSpace(valStr)
	valCalculated := 0
	if strings.Contains(valStr, "%") {
		valStr = strings.Replace(valStr, "%", "", -1)
		val, err := strconv.Atoi(valStr)
		if err == nil && baseObjectPointer != nil {
			valCalculated = valueStringToNumber(
				baseObjectPointer.Attr[baseObjAttrName],
				baseObjectPointer.Parent,
				baseObjAttrName) * val / 100
		}
	} else { // no % in the value
		val, err := strconv.Atoi(valStr)
		if err == nil {
			valCalculated = val
		}
	}
	return valCalculated
}
