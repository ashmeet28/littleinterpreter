package main

type TokenType int

const (
	// Token Types

	TT_ILLEGAL TokenType = iota
	TT_EOF
	TT_NEW_LINE
	TT_SPACE

	TT_IDENT // main
	TT_INT   // 12345
	TT_STR   // "abc"

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

	var newTok TokenInfo

	if len(src) == 0 {

		newTok.tokType = TT_EOF
		bytesConsumed = 0
		return newTok, bytesConsumed

	} else if src[0] == 0x20 {

		newTok.tokType = TT_SPACE
		bytesConsumed = 1
		return newTok, bytesConsumed

	} else if src[0] == 0x0a {

		newTok.tokType = TT_NEW_LINE
		bytesConsumed = 1
		return newTok, bytesConsumed

	} else if len(src) > 2 && src[0] == 0x2f && src[1] == 0x2f {

		newTok.tokType = TT_NEW_LINE
		bytesConsumed = 0
		for _, b := range src {
			bytesConsumed++
			if b == 0x0a {
				return newTok, bytesConsumed
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
			(newTok.tokType == TT_ILLEGAL || len(newTok.tokStr) < len(tokStr)) {

			newTok.tokType = tokType
			newTok.tokStr = tokStr
			bytesConsumed = len(newTok.tokStr)

		}
	}

	if newTok.tokType != TT_ILLEGAL {
		return newTok, bytesConsumed
	}

	isDigit := func(c byte) bool {
		return c >= 0x30 && c <= 0x39
	}

	isAplabet := func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) || (c == 0x5f)
	}

	var i int = 0

	if isAplabet(srcLine[i]) {

		newTok.tokType = TT_IDENT
		for (i < len(srcLine)) && (isAplabet(srcLine[i]) || isDigit(srcLine[i])) {
			i++
		}
		newTok.tokStr = srcLine[:i]
		bytesConsumed = len(newTok.tokStr)

	} else if isDigit(srcLine[i]) {

		newTok.tokType = TT_INT
		for (i < len(srcLine)) && (isAplabet(srcLine[i]) || isDigit(srcLine[i])) {
			i++
		}
		newTok.tokStr = srcLine[:i]
		bytesConsumed = len(newTok.tokStr)

	} else if srcLine[i] == 0x22 {

		i++
		for i < len(srcLine) {
			if srcLine[i] == 0x22 {
				newTok.tokType = TT_STR
				newTok.tokStr = srcLine[1:i]
				bytesConsumed = len(newTok.tokStr) + 2
				break
			} else {
				i++
			}
		}

	}

	return newTok, bytesConsumed
}

func GenerateTokens(src []byte) []TokenInfo {
	var toks []TokenInfo

	var isDone bool = false

	for !isDone {
		newTok, bytesConsumed := GenerateToken(src)

		if newTok.tokType == TT_ILLEGAL {
			panic("Error while tokenizing")
		}

		if newTok.tokType != TT_SPACE {
			toks = append(toks, newTok)
		}

		if newTok.tokType == TT_EOF {
			isDone = true
		} else {
			src = src[bytesConsumed:]
		}
	}

	return toks
}
