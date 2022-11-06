package iris

import (
	"fmt"
	"strconv"
	"strings"
)

// long, pre-defined values:
var ValFromLeftToRight = "FromLeftToRight"
var ValFromTopToDown = "FromTopToDown"

// pre-defined keys:
var KeyBornXinParent = "bornXinParent"
var KeyBornYinParent = "bornYinParent"
var KeyBornHorizontalWay = "bornHorizontalWay"
var KeyBornVerticalWay = "bornVerticalWay "

// when you print the output in terminal you have to use
func NewLineInTerminalPrinting() string {
	return "\n"
}

func NewlineSeparatorInText() string {
	return "<br />"
}

func IsHtmlLineSeparatorInText(txt string) bool {
	txt = strings.Replace(txt, " ", "", -1)
	if txt == "<br>" || txt == "<br/>" {
		return true
	}
	return false
}

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
// these are (x,y) paired point coords
func (obj DomObj) RenderScreenMatrix(useBasicFiller bool) RenderedScreen {

	objWidth := obj.attribCalculated("width")
	objHeight := obj.attribCalculated("height")
	fmt.Println("TEST RenderScreenMatrix width, height:", objWidth, objHeight)

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

// the screen matrix can be addressed by x,y coordinates.
// this function is the raw text output
func (obj DomObj) RenderScreenMatrixToTxt(useBasicFiller bool) string {
	screenReceived := obj.RenderScreenMatrix(useBasicFiller)
	fmt.Println("TEST screen received:", screenReceived)
	return ScreenToTxt(screenReceived)
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
			out = append(out, NewLineInTerminalPrinting())
		}
	}
	return strings.Join(out[:], "")
}

func ObjNew(id, width, height, bornXinParentVal, bornYinParentVal, basicBgFiller string, parentPointer *DomObj) DomObj {
	attr := map[string]string{
		"id": id, // all id is unique. Equivalent with dom Id

		"width":  width,  // value examples: 20, 20%,
		"height": height, // or it can be empty: ''

		KeyBornXinParent:     bornXinParentVal,   // default unit: character.
		KeyBornYinParent:     bornYinParentVal,   // 20 means: 20 character
		KeyBornHorizontalWay: ValFromLeftToRight, // or FromRightToLeft
		KeyBornVerticalWay:   ValFromTopToDown,   // or FromDownToTop

		"basicBackgroundFiller": basicBgFiller, // this string represents the object in debug/dev mode
		// the displayed filler is one char wide. But: it can be longer than once char, with color codes for example

		"text": "", // if something has text, this is a LEAF elem,
		// it can't have children!
	}
	return DomObj{Attr: attr, Parent: parentPointer}
}

func (obj DomObj) WidthFixOrTextBased() int {
	attribConverted, objWidth := valueStringToNumber(obj.Attr["width"], obj.Parent, "width")
	if !attribConverted { // if not pre-defined width:
		objTxt := obj.Attr["text"]
		objWidth = 0
		if len(objTxt) > 0 { // find the widest line
			for _, row := range strings.Split(objTxt, NewlineSeparatorInText()) {
				if len(row) > objWidth {
					objWidth = len(row)
				}
			}
		}
	}
	return objWidth
}

func (obj DomObj) HeightFixOrTextBased() int {
	attribConverted, objHeight := valueStringToNumber(obj.Attr["height"], obj.Parent, "height")
	if !attribConverted { // if not pre-defined width:
		objTxt := obj.Attr["text"]
		objHeight = 0
		if len(objTxt) > 0 { // find the num of newlines in the text
			objHeight = len(strings.Split(objTxt, NewlineSeparatorInText()))
		}
	}
	return objHeight
}

//  obj relative positions in the parent:
func (obj DomObj) PositionsInParent() (int, int, int, int) {
	left, right, top, bottom := 0, 0, 0, 0
	toInt := strconv.Atoi

	objWidth := obj.WidthFixOrTextBased()
	objHeight := obj.HeightFixOrTextBased()

	if obj.Attr[KeyBornHorizontalWay] == ValFromLeftToRight {
		left, _ = toInt(obj.Attr[KeyBornXinParent]) // if left == 2, width = 3,
		right = left + objWidth - 1                 // then horizontal used elems: 2,3,4 -> rigth = 2 + width - 1
	} else {
		right, _ = toInt(obj.Attr[KeyBornXinParent])
		left = right - objWidth + 1
	}
	if obj.Attr[KeyBornVerticalWay] == ValFromTopToDown {
		top, _ = toInt(obj.Attr[KeyBornYinParent])
		bottom = top + objHeight - 1
	} else {
		bottom, _ = toInt(obj.Attr[KeyBornYinParent])
		top = bottom - objHeight + 1
	}
	return left, right, top, bottom
}

// Tested
// return with a numeric attribute value
func (obj DomObj) attribCalculated(key string) int {
	_, keyInAttr := obj.Attr[key]
	if !keyInAttr {
		return 0
	}
	if len(obj.Attr[key]) == 0 { // if the Attribute value is empty
		return 0
	}

	attribVal := 0

	// the attrib is not set, so it depends on the children!
	// or on the text width
	if key == "width" || key == "height" {

		if key == "width" {
			attribVal = obj.WidthFixOrTextBased()
		}
		if key == "height" {
			attribVal = obj.HeightFixOrTextBased()
		}
		// no txtLen, and no pre-sed width.
		/*
			x_min, x_max, y_min, y_max := 0, 0, 0, 0
			for child := range obj.Children {
			}
		*/
	} else { // something else than width/height attrib reading
		strAttribConverted, valueConverted := valueStringToNumber(obj.Attr[key], obj.Parent, key)
		if strAttribConverted {
			return valueConverted
		}
		return 0 // if we have conversion error, return with 0
	}

	return attribVal // if the attrib is
}

// Tested
func valueStringToNumber(valStr string, baseObjectPointer *DomObj, baseObjAttrName string) (bool, int) {
	valStr = strings.TrimSpace(valStr)
	valCalculated := 0
	stringConverted := false
	if strings.Contains(valStr, "%") {
		valStr = strings.Replace(valStr, "%", "", -1)
		percentVal, err := strconv.Atoi(valStr)
		if err == nil && baseObjectPointer != nil {
			strConverted, valCalculated2 := valueStringToNumber(
				baseObjectPointer.Attr[baseObjAttrName],
				baseObjectPointer.Parent,
				baseObjAttrName)
			valCalculated2 = valCalculated2 * percentVal / 100
			if strConverted {
				valCalculated = valCalculated2
				stringConverted = true
			}
		}
	} else { // no % sign
		if len(valStr) > 0 { // if valStr is defined, not-empty
			val, err := strconv.Atoi(valStr)
			if err == nil {
				valCalculated = val
				stringConverted = true
			}
		}
	}
	return stringConverted, valCalculated
}
