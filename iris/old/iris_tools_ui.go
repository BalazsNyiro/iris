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
	"fmt"
	"os"
	"strings"
)

// /////////////////////////////////////////////////
// keypress detection is based on this example:
// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
// thank you.
func channelRead_userInput(ch chan string) {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		ch <- string(b)
	}
} ///////////////////////////////////////////////////

func channelRead_terminalSizeChangeDetect(ch chan [2]int) {
	widthSys, heightSys := 0, 0
	for {
		widthSysNow, heightSysNow := TerminalDimensionsWithSyscall()
		if widthSysNow != widthSys || heightSysNow != heightSys {
			widthSys = widthSysNow
			heightSys = heightSysNow
			ch <- [2]int{widthSys, heightSys}
		}
		TimeSleep(TimeIntervalTerminalSizeDetectMillisec)
	}
}

func channelRead_dataInputInterpret(ch_data_input chan MessageAndCharactersForWindowsUpdate, windows *Windows, dataInputLineSeparator string) {
	for {
		select {
		case dataInput, _ := <-ch_data_input:
			fmt.Println("\ndata input:", dataInput)
			if strings.HasPrefix(strings.TrimSpace(dataInput.msg), "select:win") {
				dataInputProcessLineByLine(dataInput, windows, dataInputLineSeparator)
			}
		default:
			_ = ""
		}
	}
}

func dataInputProcessLineByLine(dataInput MessageAndCharactersForWindowsUpdate, windows *Windows, dataInputLineSeparator string) string {
	winId := ""
	fieldSeparator := ":"

	/*
				when an attribute is set, the line is split by separators,
			    to get the path what we need to do.
				For example:
				set:borderBottom:=

				it means: set borderBottom =
				so set the borderBottom value to =

				but what if we want to use the possible separator char, too?
				it can be escaped, for example, but it is ugly.

				So after a 'set:borderBottom:' the program knows that every char is accepted,
		        There is NO MORE SEPARATOR CHAR splitting.

	*/
	everythingIsValueFromHere := func(elems []string, useEverythingfromHere int) string {
		return strings.Join(elems[useEverythingfromHere:], fieldSeparator)
	}

	for _, lineOrig := range strings.Split(dataInput.msg, dataInputLineSeparator) {
		line := strings.TrimSpace(lineOrig)
		elems := strings.Split(line, fieldSeparator)

		// select:win:nameOfWin
		if elems[0] == "select" && elems[1] == "win" {
			winId = strings.TrimSpace(elems[2])
			if _, exist := (*windows)[winId]; !exist {
				(*windows)[winId] = Window{winId: winId}
			}

			// process only the first line here, then later add all other lines, too
			if len(dataInput.addLine) > 0 {
				win := (*windows)[winId]
				win.lines = append(win.lines, dataInput.addLine)
				(*windows)[winId] = win
			}

		}

		win := (*windows)[winId]

		if elems[0] == "set" {
			if elems[1] == "backgroundDefault" {
				win.backgroundDefault = elems[2]
			}
			if elems[1] == "xLeft" {
				win.xLeft = Str2Int(elems[2])
			}
			if elems[1] == "width" {
				win.width = Str2Int(elems[2])
			}
			if elems[1] == "yTop" {
				win.yTop = Str2Int(elems[2])
			}
			if elems[1] == "height" {
				win.height = Str2Int(elems[2])
			}

			// ######## borders ###############
			if elems[1] == "borderLeft" {
				// "[=>", [ is outer, > is inner
				win.borderLeft = elems[2]
			}

			// in right directions, the inner is in left, outer is in right
			if elems[1] == "borderRight" { // "<=]" < is inner, ] is outer
				win.borderRight = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderTop" { // "Oo." O is outer, . is inner
				win.borderTop = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderBottom" { // .:| . is inner, :middle, | is outer
				win.borderBottom = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderLeftTop" { // abc a is outer, c is inner
				win.borderLeftTop = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderRightTop" { // def  f is inner, d is outer
				win.borderRightTop = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderLeftBottom" { // ghi g is outer, i is inner
				win.borderLeftBottom = everythingIsValueFromHere(elems, 2)
			}
			if elems[1] == "borderRightBottom" { // lkj  l is inner, j is outer
				win.borderRightBottom = everythingIsValueFromHere(elems, 2)
			}
		}

		(*windows)[winId] = win

		if winId == "" {
			continue
		}

	}

	return winId
}

func action_of_user_input(stdin string) string {
	action := ""
	if stdin == "q" {
		action = "quit"
	}

	if stdin == "l" {
	}
	if stdin == "h" {
	}
	if stdin == "j" {
	}
	if stdin == "k" {
	}
	return action
}
