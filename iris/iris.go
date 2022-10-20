// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
    "strings"
)

func Page() string {
    fmt.Println("iris")
    terminal_detect()
    return "page from iris"
}

// https://stackoverflow.com/questions/263890/how-do-i-find-the-width-height-of-a-terminal-window
func terminal_detect() (int, int) {
    out, _ := shell("ls -la")
    out, _ = shell("stty size")
    fmt.Println(">>>", out)
    return 1, 2
}

// https://zetcode.com/golang/exec-command/
func shell(commandAndParams string) (string, error) {
    args := strings.Fields(commandAndParams)
    fmt.Println("args:", args)
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stdin = strings.NewReader("some input")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Println(err)
    }
    return out.String(), err
}
