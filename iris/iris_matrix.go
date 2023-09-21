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
	"strings"
)

type ColumnChars []Char
type MatrixChars struct {
	matrix      []ColumnChars
	defaultRune rune
	width       int
	height      int
}

func (mx MatrixChars) LineFill(characters []Char, lineNumZeroBased int) {
	matrixWidth := len(mx.matrix)
	for id, lineChar := range characters {
		if id >= matrixWidth {
			break
		}
		mx.matrix[id][lineNumZeroBased] = lineChar
	}
}

func (mx MatrixChars) DisplayInConsoleToDebugOrAnalyse() {
	fmt.Println("matrix width:", mx.width, "matrix height:", mx.height)
	for y := 0; y < mx.height; y++ {
		for x := 0; x < mx.width; x++ {
			fmt.Print(mx.matrix[x][y].display())
		}
		fmt.Println()
	}
}
func (mx MatrixChars) Lines_representation_to_test_comparison() []string {
	// list of lines
	output := []string{}
	for y := 0; y < mx.height; y++ {
		line := []string{}
		for x := 0; x < mx.width; x++ {
			line = append(line, mx.matrix[x][y].display())
		}
		output = append(output, strings.Join(line[:], ""))
	}
	return output
}

func MatrixNew(width, height int, backgroundDefault string) MatrixChars {
	defaultRune := 'r'
	if len(backgroundDefault) > 0 {
		defaultRune = rune(backgroundDefault[0])
	}
	matrixNew := MatrixChars{defaultRune: defaultRune, width: width, height: height}
	for x := 0; x < width; x++ {
		column := ColumnChars{}
		for y := 0; y < height; y++ {
			column = append(column, Char{runeVal: defaultRune})
		}
		matrixNew.matrix = append(matrixNew.matrix, column)
	}
	// fmt.Println("new layer:", matrixNew)
	return matrixNew
}

func MatrixAdd(matrixBase, matrixAdded MatrixChars, insertX, insertY int) MatrixChars {
	// loop over Added matrix, and insert the points into Base
	for addedY := 0; addedY < matrixAdded.height; addedY++ {
		for addedX := 0; addedX < matrixAdded.width; addedX++ {
			calculatedXInBase := insertX + addedX
			calculatedYInBase := insertY + addedY
			if calculatedXInBase >= 0 && calculatedXInBase < matrixBase.width {
				if calculatedYInBase >= 0 && calculatedYInBase < matrixBase.height {
					matrixBase.matrix[calculatedXInBase][calculatedYInBase] = matrixAdded.matrix[addedX][addedY]
				}
			}
		}
	}
	return matrixBase
}
