// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
// TESTED MANUALLY
func TerminalDimensionsWithSyscall() (int, int) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col), int(ws.Row)
}

// https://stackoverflow.com/questions/263890/how-do-i-find-the-width-height-of-a-terminal-window
// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
// TESTED MANUALLY, can't detect terminal size from 'go test'
func TerminalDimensionsSttySize() (int, int) {
	out, _ := shell("stty size")
	out = strings.TrimSpace(out)

	row := 0
	col := 0

	if strings.Contains(out, " ") {
		split := strings.Split(out, " ")
		rowDetected, err := strconv.Atoi(split[0])
		colDetected, err2 := strconv.Atoi(split[1])

		if err != nil || err2 != nil {
			// log.Fatal("stty size, not integer reply:", row, col)
			// no available terminal size
		} else {
			row = rowDetected
			col = colDetected
		}
	}

	return col, row
}

// https://zetcode.com/golang/exec-command/
// TESTED manually
func shell(commandAndParams string) (string, error) {
	return shellCore(commandAndParams, "")
}

// TESTED manually
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
