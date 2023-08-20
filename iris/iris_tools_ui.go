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

	for _, lineOrig := range strings.Split(dataInput.msg, dataInputLineSeparator) {
		line := strings.TrimSpace(lineOrig)
		elems := strings.Split(line, ":")

		if len(elems) == 3 {

			// select:win:nameOfWin
			if elems[0] == "select" && elems[1] == "win" {
				winId = strings.TrimSpace(elems[2])
				if _, exist := (*windows)[winId]; !exist {
					(*windows)[winId] = Window{winId: winId,
						borderLeft:  "[>", // direction: from left to right
						borderRight: "<]", // < is inner, ] is outer

						// A is outer, C is inner
						borderTop:    "ABC", //"-v", // - is outer, v is inner
						borderBottom: ".:|", // . is inner, :middle, | is outer

						// from left-Top to right-down:  L
						borderLeftTop: "Lt", //           t

						// from right-top to left-down  R
						borderRightTop: "Rt", //       t

						//                            b
						borderLeftBottom: "Lb", // L

						//                            b
						borderRightBottom: "bR", //   R
					}
				}

				// process only the first line here, then later add all other lines, too
				if len(dataInput.addLine) > 0 {
					win := (*windows)[winId]
					win.lines = append(win.lines, dataInput.addLine)
					(*windows)[winId] = win
				}

			}

			win := (*windows)[winId]
			if elems[0] == "set" && elems[1] == "backgroundDefault" {
				win.backgroundDefault = elems[2]
			}
			if elems[0] == "set" && elems[1] == "xLeft" {
				win.xLeft = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "width" {
				win.width = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "yTop" {
				win.yTop = Str2Int(elems[2])
			}
			if elems[0] == "set" && elems[1] == "height" {
				win.height = Str2Int(elems[2])
			}
			(*windows)[winId] = win

		}

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
