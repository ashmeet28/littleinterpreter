package main

import (
	"fmt"
	"os"
	"strconv"
)

type SymType int

const (
	ST_ILLEGAL SymType = iota

	ST_FUNC
	ST_VAR
)

type SymInfo struct {
	symType     SymType
	symIdent    string
	symAddr     uint32
	symScope    int
	symArgCount int
}

func GenerateBytecode(toks []TokenInfo) []byte {
	var bytecode []byte

	var symTable []SymInfo

	var blankLiterals []int

	var GLOBAL_SCOPE int = 1
	var currScope int = GLOBAL_SCOPE

	var currTokIndex int

	emitInst := func(op byte) {
		bytecode = append(bytecode, op)
	}

	emitPushLitInst := func(v uint32) {
		bytecode = append(bytecode, OP_PUSH_LITERAL,
			uint8(v&0xff), uint8((v>>8)&0xff), uint8((v>>16)&0xff), uint8((v>>24)&0xff))
	}

	emitPushLitBlackInst := func() {
		bytecode = append(bytecode, OP_PUSH_LITERAL, 0, 0, 0, 0)
		blankLiterals = append(blankLiterals, len(bytecode)-4)
	}

	fillPushLitBlackInst := func(v uint32) {
		i := blankLiterals[len(blankLiterals)-1]
		bytecode[i] = uint8(v & 0xff)
		bytecode[i+1] = uint8((v >> 8) & 0xff)
		bytecode[i+2] = uint8((v >> 16) & 0xff)
		bytecode[i+3] = uint8((v >> 24) & 0xff)
	}

	peek := func() TokenInfo {
		return toks[currTokIndex]
	}

	advance := func() TokenInfo {
		currTokIndex++
		return toks[currTokIndex-1]
	}

	consume := func(tokType TokenType) TokenInfo {
		if peek().tokType != tokType {
			fmt.Println("Error while compiling")
			os.Exit(1)
		}
		return advance()
	}

	addVarToSymTable := func(varIdent string, varScope int) int {
		var varAddr uint32

		for _, s := range symTable {
			if s.symType == ST_VAR {
				if s.symScope == GLOBAL_SCOPE && varScope == GLOBAL_SCOPE {
					varAddr++
				} else if s.symScope > GLOBAL_SCOPE && varScope > GLOBAL_SCOPE {
					varAddr++
				}
			}
		}

		var s SymInfo
		s.symType = ST_VAR
		s.symIdent = varIdent
		s.symAddr = varAddr
		s.symScope = varScope
		s.symArgCount = 0

		var index int = len(symTable)
		symTable = append(symTable, s)
		return index
	}

	addFuncToSymTable := func(funcIdent string) int {
		var s SymInfo
		s.symType = ST_FUNC
		s.symIdent = funcIdent
		s.symAddr = uint32(len(bytecode))
		s.symScope = GLOBAL_SCOPE
		s.symArgCount = 0

		var index int = len(symTable)
		symTable = append(symTable, s)
		return index
	}

	removeSymFromSymTable := func() {
		var newSymTable []SymInfo
		for _, b := range symTable {
			if b.symScope <= currScope {
				newSymTable = append(newSymTable, b)
			}
		}
		symTable = newSymTable
	}

	findSym := func(symIdent string) SymInfo {
		var currSym SymInfo
		currSym.symType = ST_ILLEGAL
		currSym.symScope = GLOBAL_SCOPE
		for _, s := range symTable {
			if s.symIdent == symIdent &&
				(s.symScope > currSym.symScope || currSym.symType == ST_ILLEGAL) {
				currSym = s
			}
		}
		return currSym
	}

	isTokBinaryOp := func(tok TokenInfo) bool {
		tokTypes := []TokenType{
			TT_ADD, TT_SUB, TT_MUL, TT_QUO, TT_REM, TT_AND, TT_OR, TT_XOR, TT_SHL,
			TT_SHR, TT_LAND, TT_LOR, TT_EQL, TT_LSS, TT_GTR, TT_NEQ, TT_LEQ, TT_GEQ,
		}
		for _, t := range tokTypes {
			if t == tok.tokType {
				return true
			}
		}
		return false
	}

	BinaryTokOpcode := map[TokenType]byte{
		TT_ADD:  OP_ADD,
		TT_SUB:  OP_SUB,
		TT_MUL:  OP_MUL,
		TT_QUO:  OP_QUO,
		TT_REM:  OP_REM,
		TT_AND:  OP_AND,
		TT_OR:   OP_OR,
		TT_XOR:  OP_XOR,
		TT_SHL:  OP_SHL,
		TT_SHR:  OP_SHR,
		TT_LAND: OP_LAND,
		TT_LOR:  OP_LOR,
		TT_EQL:  OP_EQL,
		TT_LSS:  OP_LSS,
		TT_GTR:  OP_GTR,
		TT_NEQ:  OP_NEQ,
		TT_LEQ:  OP_LEQ,
		TT_GEQ:  OP_GEQ,
	}

	var compileExpr func()
	var compileGrouping func()
	var compileUnary func(bool)
	var compileBinary func()

	compileExpr = func() {
		compileUnary(false)
	}

	compileGrouping = func() {
		consume(TT_LPAREN)
		var c int = 0
		if peek().tokType == TT_MUL {
			for peek().tokType == TT_MUL {
				consume(TT_MUL)
				c++
			}
			compileGrouping()
			for c != 0 {
				emitInst(OP_LOAD_MEM)
				c--
			}
		} else {
			compileUnary(false)
		}
		consume(TT_RPAREN)
	}

	compileUnary = func(isRightOfBinary bool) {
		switch peek().tokType {
		case TT_IDENT:
			emitPushLitInst(findSym(consume(TT_IDENT).tokStr).symAddr)
			emitInst(OP_LOAD_LOCAL)
		case TT_INT:
			v, _ := strconv.ParseInt(consume(TT_INT).tokStr, 0, 64)
			emitPushLitInst(uint32(v))
		default:
			compileGrouping()
		}

		if (!isRightOfBinary) && isTokBinaryOp(peek()) {
			compileBinary()
		}
	}

	compileBinary = func() {
		opTok := advance()
		compileUnary(true)
		emitInst(BinaryTokOpcode[opTok.tokType])
		if isTokBinaryOp(peek()) {
			compileBinary()
		}
	}

	emitPushLitBlackInst()
	emitInst(OP_JUMP)

	for peek().tokType != TT_EOF {
		switch peek().tokType {
		case TT_VAR:
			consume(TT_VAR)
			addVarToSymTable(consume(TT_IDENT).tokStr, currScope)
			consume(TT_NEW_LINE)

		case TT_FUNC:
			consume(TT_FUNC)
			i := addFuncToSymTable(consume(TT_IDENT).tokStr)
			currScope++
			consume(TT_LPAREN)
			for peek().tokType != TT_NEW_LINE {
				if peek().tokType != TT_RPAREN {
					addVarToSymTable(consume(TT_IDENT).tokStr, currScope)
					symTable[i].symArgCount++
				}
				advance()
			}
			consume(TT_NEW_LINE)

		case TT_IDENT:
			sym := findSym(consume(TT_IDENT).tokStr)
			if sym.symType == ST_VAR {
				consume(TT_ASSIGN)
				compileExpr()
				emitPushLitInst(sym.symAddr)
				emitInst(OP_STORE_LOCAL)
				consume(TT_NEW_LINE)
			} else {
				compileExpr()
				emitInst(OP_POP_LITERAL)
			}
		case TT_END:
			consume(TT_END)
			currScope--
			removeSymFromSymTable()
		case TT_NEW_LINE:
			consume(TT_NEW_LINE)
		default:
			fmt.Println("Error while compiling")
			os.Exit(1)
		}
	}

	fillPushLitBlackInst(findSym("main").symAddr)

	return bytecode
}