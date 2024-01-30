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

	return newTok, bytesConsumed
}
