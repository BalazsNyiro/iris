package iris

func NewLine() string { return "\n" }

type RenderedColumn []string
type RenderedScreen []RenderedColumn // there are vertical columns next to each other.

/*
I store everything in strings.
 , is a list separator so never use it as a key or a value
*/
type Window map[string]string
type Windows map[string]Window

func WindowsNew() Windows {
	Win := Windows{}
	return WinNew(Win, "root")
}

func WinNew(windows Windows, id string) Windows {
	return windows
}
