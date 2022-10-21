// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
)

func Page() string {
    fmt.Println("iris")
    TerminalDetect()
    return "page from iris"
}

// https://stackoverflow.com/questions/263890/how-do-i-find-the-width-height-of-a-terminal-window
// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
func TerminalDetect() (int, int) {
    out, _ := shell("stty size")
    fmt.Println("terminal detect >>>", out)
    return 1, 2
}

// https://zetcode.com/golang/exec-command/
func shell(commandAndParams string) (string, error) {
    return shellCore(commandAndParams, "")
}

func shellCore(commandAndParams, input string) (string, error) {
    args := strings.Fields(commandAndParams)
    cmd := exec.Command(args[0], args[1:]...)
    if len(input) > 0 {
        cmd.Stdin = strings.NewReader(input)
    } else {
        cmd.Stdin = os.Stdin
    }
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Println(err)
    }
    return out.String(), err
}
