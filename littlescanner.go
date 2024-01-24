package main

import (
	"fmt"
	"os"
)

type TokenType int

const (
	// Token Types

	TT_ILLEGAL TokenType = iota
	TT_EOF
	TT_NEW_LINE
	TT_SPACE

	TT_IDENT  // main
	TT_INT    // 12345
	TT_STRING // "abc"

	TT_ADD    // +
	TT_SUB    // -
	TT_MUL    // *
	TT_QUO    // /
	TT_REM    // %
	TT_AND    // &
	TT_OR     // |
	TT_XOR    // ^
	TT_SHL    // <<
	TT_SHR    // >>
	TT_LAND   // &&
	TT_LOR    // ||
	TT_EQL    // ==
	TT_LSS    // <
	TT_GTR    // >
	TT_ASSIGN // =
	TT_NEQ    // !=
	TT_LEQ    // <=
	TT_GEQ    // >=

	TT_LPAREN // (
	TT_RPAREN // )
	TT_COMMA  // ,

	TT_WHILE
	TT_BREAK
	TT_CONTINUE
	TT_IF
	TT_ELSE
	TT_FUNC
	TT_RETURN
	TT_END
	TT_VAR
)

type TokenInfo struct {
	tokType TokenType
	tokStr  string
}

func GenerateToken(src []byte) (TokenInfo, int) {
	var bytesConsumed int = 0

	var currTok TokenInfo
	currTok.tokType = TT_ILLEGAL

	if len(src) == 0 {

		currTok.tokType = TT_EOF
		bytesConsumed = 0
		return currTok, bytesConsumed

	} else if src[0] == 0x20 {

		currTok.tokType = TT_SPACE
		bytesConsumed = 1
		return currTok, bytesConsumed

	} else if src[0] == 0x0a {

		currTok.tokType = TT_NEW_LINE
		bytesConsumed = 1
		return currTok, bytesConsumed

	} else if len(src) > 2 && src[0] == 0x2f && src[1] == 0x2f {

		currTok.tokType = TT_NEW_LINE
		bytesConsumed = 0
		for _, b := range src {
			bytesConsumed++
			if b == 0x0a {
				return currTok, bytesConsumed
			}
		}

	}

	var srcLine string

	for i, c := range src {
		if c == 0x0a {
			srcLine = string(src[:i])
			break
		}
	}

	TokensStrings := map[TokenType]string{
		TT_ADD: "+",
		TT_SUB: "-",
		TT_MUL: "*",
		TT_QUO: "/",
		TT_REM: "%",
		TT_AND: "&",
		TT_OR:  "|",
		TT_XOR: "^",
		TT_SHL: "<<",
		TT_SHR: ">>",
		TT_EQL: "==",
		TT_LSS: "<",
		TT_GTR: ">",
		TT_NEQ: "!=",
		TT_LEQ: "<=",
		TT_GEQ: ">=",

		TT_ASSIGN: "=",
		TT_LPAREN: "(",
		TT_RPAREN: ")",
		TT_COMMA:  ",",

		TT_WHILE:    "while",
		TT_BREAK:    "break",
		TT_CONTINUE: "continue",
		TT_IF:       "if",
		TT_ELSE:     "else",
		TT_FUNC:     "func",
		TT_RETURN:   "return",
		TT_END:      "end",
		TT_VAR:      "var",
	}

	for tokType, tokStr := range TokensStrings {
		if (len(srcLine) >= len(tokStr) && srcLine[:len(tokStr)] == tokStr) &&
			(currTok.tokType == TT_ILLEGAL || len(currTok.tokStr) < len(tokStr)) {

			currTok.tokType = tokType
			currTok.tokStr = tokStr
			bytesConsumed = len(currTok.tokStr)

		}
	}

	if currTok.tokType != TT_ILLEGAL {
		return currTok, bytesConsumed
	}

	isDigit := func(c byte) bool {
		return c >= 0x30 && c <= 0x39
	}

	isAplabet := func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) || (c == 0x5f)
	}

	var i int = 0

	if isAplabet(srcLine[i]) {

		currTok.tokType = TT_IDENT
		for (i < len(srcLine)) && (isAplabet(srcLine[i]) || isDigit(srcLine[i])) {
			i++
		}
		currTok.tokStr = srcLine[:i]
		bytesConsumed = len(currTok.tokStr)

	} else if isDigit(srcLine[i]) {

		currTok.tokType = TT_INT
		for (i < len(srcLine)) && (isAplabet(srcLine[i]) || isDigit(srcLine[i])) {
			i++
		}
		currTok.tokStr = srcLine[:i]
		bytesConsumed = len(currTok.tokStr)

	} else if srcLine[i] == 0x27 {

		i++
		for i < len(srcLine) {
			if srcLine[i] == 0x27 {
				currTok.tokType = TT_STRING
				currTok.tokStr = srcLine[1:i]
				bytesConsumed = len(currTok.tokStr) + 2
				break
			} else {
				i++
			}
		}

	}

	return currTok, bytesConsumed
}

func GenerateTokens(src []byte) []TokenInfo {
	var toks []TokenInfo

	var isDone bool = false

	for !isDone {
		currTok, bytesConsumed := GenerateToken(src)

		if currTok.tokType == TT_ILLEGAL {
			fmt.Println("Error while tokenizing")
			os.Exit(1)
		}

		if currTok.tokType != TT_SPACE {
			toks = append(toks, currTok)
		}

		if currTok.tokType == TT_EOF {
			isDone = true
		} else {
			src = src[bytesConsumed:]
		}
	}

	return toks
}
