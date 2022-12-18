// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package iris

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var Digits = "0123456789"

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

func OsDetect() string {
	os := runtime.GOOS
	if strings.Contains("windows|darwin|linux", os) {
		return os
	}
	return "linux" // if we have an exotic os, we will handle it as linux
	// return "unknown"
}

// TESTED - mainly used in test functions
func IsNumber(txt string) bool {
	plusMinusDetected := false
	normalCharDetected := false
	txt = strings.TrimSpace(txt)
	if len(txt) == 0 {
		return false // empty string is not a number
	}
	for id, rune := range txt {
		if id == 0 && (rune == '+' || rune == '-') {
			plusMinusDetected = true
			continue
		}
		if !strings.Contains(Digits, string(rune)) {
			return false
		} else {
			normalCharDetected = true
		}
	}
	// only plusMinus is detected
	if plusMinusDetected && !normalCharDetected {
		return false
	}
	return true
}

// wrapper, not tested
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// wrapper, not tested
func Atoi(txt string) int {
	num, error := strconv.Atoi(txt)
	if error == nil {
		return num
	}
	fmt.Println("Atoi error: ", error)
	return 0
}

// TESTED
func StrMath(a, operator, b string) string {
	a_int := Atoi(a)
	b_int := Atoi(b)
	if operator == "-" {
		return Itoa(a_int - b_int)
	}
	if operator == "+" {
		return Itoa(a_int + b_int)
	}
	if operator == "*" {
		return Itoa(a_int * b_int)
	}
	if operator == "/" {
		if b_int != 0 {
			return Itoa(a_int / b_int)
		} else {
			fmt.Println("zero division", a_int, operator, b_int)
		}
	}
	fmt.Println("Math Error: ", a_int, operator, b_int)
	return "0"
}

// TESTED
func StrDoubleSpacesRemove(txt string) string {
	for strings.Contains(txt, "  ") {
		txt = strings.Replace(txt, "  ", " ", -1)
	}
	return txt
}

// TESTED
func StrListRemoveEmptyElems(list []string, useTrim bool) []string {
	cleaned := []string{}
	for _, elem := range list {
		if useTrim {
			elem = strings.TrimSpace(elem)
		}
		if len(elem) > 0 {
			cleaned = append(cleaned, elem)
		}
	}
	return cleaned
}

// TESTED
func ExprOperatorIsValid(operatorChecked string) bool {
	for _, operatorKnown := range strings.Split("+,-,*,/", ",") {
		if operatorChecked == operatorKnown {
			return true
		}
	}
	return false
}

func debug_info_save(windows Windows) {

	f, _ := os.Create("debug_iris.txt")
	defer f.Close()

	f.Write([]byte("===============\n"))
	for key, val := range windows["prgState"] {
		message := fmt.Sprintf("%s: %s\n", key, val)
		data := []byte(message)
		f.Write(data)
	}
	for _, winName := range windows_get_win_names_publics_sorted(windows) {
		winInfo := windows[winName]
		data := []byte(fmt.Sprintf("win public: %s (%s, %s)\n",
			winName,
			winInfo[KeyXleftCalculated],
			winInfo[KeyYtopCalculated]))
		f.Write(data)
	}
}

// collect win names, I want to sort it!
func windows_get_win_names_publics_sorted(windows Windows) []string {
	winNames := []string{}
	for winName, _ := range windows {
		if win_name_is_publics(winName) {
			winNames = append(winNames, winName)
		}
	}
	sort.Strings(winNames)
	return winNames
}

// from windows -> windows public list
// map keys are always unsorted!
func windows_keep_publics(windows Windows) Windows {
	windows_publics := Windows{}
	for winName, value := range windows {
		if win_name_is_publics(winName) {
			windows_publics[winName] = value
		}
	}
	return windows_publics
}

// windows Names -> windows Names: remove internal/non-public window names
func win_names_keep_publics(winNames []string, sort_the_names bool) []string {
	publicNames := []string{}
	for _, name := range winNames {
		if win_name_is_publics(name) {
			publicNames = append(publicNames, name)
		}
	}
	if sort_the_names {
		sort.Strings(publicNames)
	}
	return publicNames
}

// window name is public?
func win_name_is_publics(winName string) bool {
	public := true
	if winName == "prgState" { // prgState is an internal key-value storage
		public = false
	}
	return public
}
