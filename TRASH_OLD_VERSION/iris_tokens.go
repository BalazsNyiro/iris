package iris

import "strings"

// if the next operator is "": there is no more operator
// TESTED
func TokenOperatorNext(tokens []string) (int, string) { // DIA
	operatorNext := "unknown"
	tokens = StrListRemoveEmptyElems(tokens, true)
	// math operator precedence: * / are the first
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("*,/", token) {
			operatorNext = token
			return id, operatorNext
		}
	}
	for id, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) == 1 && strings.Contains("+,-", token) {
			operatorNext = token
			return id, operatorNext
		}
	}
	return -1, operatorNext
}

func TokenReplaceWinPlaceholders(windows Windows, tokens []string) []string { // DIA
	tokens = StrListRemoveEmptyElems(tokens, true)
	for id, token := range tokens {
		if len(token) > 4 && token[0:4] == "win:" { // win:Terminal:xRightCalculated
			splitted := strings.Split(token, ":")
			winName := splitted[1]
			attrib := splitted[2]

			tokens[id] = "0" // set the normal value if key/attrib exists:
			if winObj, keyInMap := windows[winName]; keyInMap {
				if valueAttrib, attribInMap := winObj[attrib]; attribInMap {
					tokens[id] = valueAttrib
				}
			}
		}
	}
	return tokens
}

// TESTED
func TokenParametersCollect(tokens []string, tokenId int) (string, string, int, int, string) { // DIA
	errMsg := ""
	valueLeft := ""
	valueRight := ""

	idValueLeft, idValueRight := tokenId-1, tokenId+1
	if idValueLeft < 0 {
		errMsg = errMsg + "express param left id < 0:" + Int2Str(idValueLeft) + ";"
	}
	if idValueRight < 0 {
		errMsg = errMsg + "express param right id < 0:" + Int2Str(idValueRight) + ";"
	}
	idMax := len(tokens) - 1
	idMaxStr := Int2Str(idMax)
	if idValueLeft > idMax {
		errMsg = errMsg + "express param left id > len(tokens)-1:" + Int2Str(idValueLeft) + " len tokens: " + idMaxStr + ";"
	}
	if idValueRight > idMax {
		errMsg = errMsg + "express param right id > len(tokens)-1:" + Int2Str(idValueRight) + " len tokens: " + idMaxStr + ";"
	}
	if errMsg == "" {
		valueLeft = tokens[idValueLeft]
		valueRight = tokens[idValueRight]
	}

	return valueLeft, valueRight, idValueLeft, idValueRight, errMsg
}
