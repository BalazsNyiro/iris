// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

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
    */
    Attr     map[string]string
    Children []DomObj
    Parent   *DomObj
}

func ObjNew(id, width, height string) DomObj {
    attr := map[string]string{"id": id, "width": width, "height": height}
    return DomObj{Attr: attr}
}

func DocumentCreate(id, width, height string) DomObj {
    root := ObjNew(id, width, height)
    return root
}
