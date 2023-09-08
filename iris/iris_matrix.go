/*
author: Balazs Nyiro, balazs.nyiro.ca@gmail.com

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package iris

type ColumnChars []Char
type MatrixChars struct {
	matrix      []ColumnChars
	defaultRune rune
}

func MatrixNew(width, height int, backgroundDefault string) MatrixChars {
	defaultRune := 'r'
	if len(backgroundDefault) > 0 {
		defaultRune = rune(backgroundDefault[0])
	}
	matrixNew := MatrixChars{defaultRune: defaultRune}
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
