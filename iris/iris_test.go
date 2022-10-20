package iris

// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

import (
    "fmt"
    "testing"
)

func Test_terminal_detect(t *testing.T) {
    fmt.Println(terminal_detect())
    fmt.Println("iris response:", Page())
}
