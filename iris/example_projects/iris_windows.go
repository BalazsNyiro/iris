package main

// run: go run iris_windows.go

// run: go run jyp_example_usage.go

func main() {
	widthSys, heightSys := iris.TerminalDimensionsWithSyscall()
	fmt.Println("syscall:", widthSys, heightSys)

	windows := iris.WindowsNewState(widthSys, heightSys)
	windows = iris.WinNew(windows, "Child", "0", "0", "8", "4", "C")
	windows = iris.WinSourceLoad(windows, "Child", "simpleText", "apple orange banana")
	windows["prgState"]["winActiveId"] = "Child"

	iris.UserInterfaceStart(windows)
}

/*
So what is a GUI?

- sender: it displays something to the user
- receiver: the user can give back new info

  the user input can be:
    - SELECTED with a SELECTOR from a displayed collection (list, buttons, pie chart, horizontal/vertical menu system)
    - SELECTED with a SELECTOR from a hidden collection (shortcuts, key combinations that can start something in the bg program)
    - INSERTED into a text field
    -
*/
