// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
    "fmt"
)

func Page() string {
    fmt.Println("iris")
    TerminalDimensionsSttySize()
    return "page from iris"
}
