package iris

import (
	"testing"
)

// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

/*
func Test_terminal_detect(t *testing.T) {
	widthStty, heightStty := TerminalDimensionsSttySize()
	fmt.Println("test stty:", widthStty, heightStty)
	widthSys, heightSys := TerminalDimensionsWithSyscall()
	fmt.Println("test syscall:", widthSys, heightSys)
}
func Test_document_create(t *testing.T) {
	rootObj := DocumentCreate("0", "50%", "50%", "40", "20")

}

*/
func Test_render(t *testing.T) {

	// calculated: width=4, height=2
	rootObj1 := DocumentCreate("0", "50%", "50%", "8", "4")
	rendered := rootObj1.RenderScreenMatrixToTxt(true)
	wanted := "RRRR" + NewLineInTerminalPrinting() + "RRRR"
	compare_str_pair(rendered, wanted, t)
}

func Test_value_string_to_number(t *testing.T) {

	valueFromEmpty := valueStringToNumber("", nil, "")
	compare_int_pair(valueFromEmpty, 0, t)

	valueNoPercentNoParent := valueStringToNumber("20", nil, "")
	compare_int_pair(valueNoPercentNoParent, 20, t)

	// % values need to know the parent measures to calculate relations
	rootObj := DocumentCreate("0", "60", "20", "", "0")
	valuePercentOfParentAttribute := valueStringToNumber("20%", &rootObj, "width")
	compare_int_pair(valuePercentOfParentAttribute, 12, t)

	rootObjHalfTerminal := DocumentCreate("0", "50%", "50%", "50", "20")
	valuePercentInHalfRoot := valueStringToNumber("20%", &rootObjHalfTerminal, "width")
	compare_int_pair(valuePercentInHalfRoot, 5, t)

	rootObj3 := DocumentCreate("0", "30%", "40%", "160", "180")
	child1 := ObjNew("child1", "25%", "50%", "0", "0", "1", &rootObj3)
	child2 := ObjNew("child2", "50%", "50%", "0", "0", "2", &child1)

	// 160*0.3*0.25*0.5
	compare_int_pair(child2.attribCalculated("width"), 6, t)

	// 180*0.4*0.5*0.5 = 18.0
	compare_int_pair(child2.attribCalculated("height"), 18, t)

	// rounded values /////////////////////////////////////////////////
	rootObjFloat := DocumentCreate("0", "53%", "45%", "14", "50")
	// 14*0.53 = 7.5 is the exact value - this is rounded:
	compare_int_pair(rootObjFloat.attribCalculated("width"), 7, t)
	// 22.5 is the exact value - this is rounded:
	compare_int_pair(rootObjFloat.attribCalculated("height"), 22, t)

	// width: 25% of 7 -> 1 because it's less than 2.
	childFloat1 := ObjNew("childFloat1", "25%", "20%", "0", "0", "F", &rootObjFloat)
	compare_int_pair(childFloat1.attribCalculated("width"), 1, t)
	// width: 20% of 22 -> 4 (exact: 4.4)
	compare_int_pair(childFloat1.attribCalculated("height"), 4, t)
}

// TEST FUNCTIONS ////////////////////////////////////////////////////////////
func compare_int_pair(received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR INT - received: %v\n  wanted: %v", received, wanted)
	}
}

func compare_str_pair(received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nERROR STR - received: %v\n  wanted: %v", received, wanted)
	}
}
