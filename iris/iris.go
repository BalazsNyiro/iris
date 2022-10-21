// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "syscall"
    "unsafe"
)

func Page() string {
    fmt.Println("iris")
    TerminalDetect()
    return "page from iris"
}

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func TerminalDimensionsWithSyscall() (uint, uint, uint, uint) {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return uint(ws.Col), uint(ws.Row), uint(ws.Xpixel), uint(ws.Ypixel)
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
