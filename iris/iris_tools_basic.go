/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package iris

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
// TESTED MANUALLY
func TerminalDimensionsWithSyscall() (int, int) { // basic fun
	return 20, 10
	/*
		ws := &winsize{}
		retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(syscall.Stdin),
			uintptr(syscall.TIOCGWINSZ),
			uintptr(unsafe.Pointer(ws)))

		if int(retCode) == -1 {
			panic(errno)
		}
		return int(ws.Col), int(ws.Row)

	*/
}

// https://stackoverflow.com/questions/263890/how-do-i-find-the-width-height-of-a-terminal-window
// https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
// TESTED MANUALLY, can't detect terminal size from 'go test'
func TerminalDimensionsSttySize() (int, int) { // basic fun
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
func shell(commandAndParams string) (string, error) { // basic fun
	return shellCore(commandAndParams, "")
}

// TESTED manually
func shellCore(commandAndParams, stdInput string) (string, error) { // basic fun
	args := strings.Fields(commandAndParams)
	cmd := exec.Command(args[0], args[1:]...)
	if len(stdInput) > 0 {
		cmd.Stdin = strings.NewReader(stdInput)
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

// Int2Str - convert an int to a string (wrapper)
func Int2Str(i int) string {
	return strconv.Itoa(i)
}

// IntMax - return with the maximum of the values
func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Str2Int - convert a string to an int
func Str2Int(txt string) int {
	num, error := strconv.Atoi(txt)
	if error == nil {
		return num
	}
	fmt.Println("Str2Int error: ", error)
	return 0
}

// StrDoubleSpacesRemove - remove double spaces
func StrDoubleSpacesRemove(txt string) string {
	for strings.Contains(txt, "  ") {
		txt = strings.Replace(txt, "  ", " ", -1)
	}
	return txt
}

// StrListRemoveEmptyElemsWithoutRealValue - remove the elems where there is no real value, only whitespaces
func StrListRemoveEmptyElemsWithoutRealValue(list []string, useTrim bool) []string {
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

// TimeSleep - wait/sleep for a given time
func TimeSleep(interval_millisec int) { // basic fun
	time.Sleep(time.Millisecond * time.Duration(interval_millisec))
}
