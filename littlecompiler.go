package main

import (
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

	var blankLitAddr []int

	var loopAddr []uint32

	var curBlocks []TokenType

	var GLOBAL_SCOPE int = 0
	var curScope int = GLOBAL_SCOPE

	var curTokIndex int

	emitInst := func(op byte) {
		bytecode = append(bytecode, op)
	}

	emitPushLitInst := func(v uint32) {
		bytecode = append(bytecode, OP_PUSH_LITERAL,
			uint8(v&0xff), uint8((v>>8)&0xff), uint8((v>>16)&0xff), uint8((v>>24)&0xff))
	}

	emitPushLitBlackInst := func() {
		bytecode = append(bytecode, OP_PUSH_LITERAL, 0, 0, 0, 0)
		blankLitAddr = append(blankLitAddr, len(bytecode)-4)
	}

	fillPushLitBlackInst := func(v uint32) {
		i := blankLitAddr[len(blankLitAddr)-1]
		bytecode[i] = uint8(v & 0xff)
		bytecode[i+1] = uint8((v >> 8) & 0xff)
		bytecode[i+2] = uint8((v >> 16) & 0xff)
		bytecode[i+3] = uint8((v >> 24) & 0xff)
		blankLitAddr = blankLitAddr[:len(blankLitAddr)-1]
	}

	peek := func() TokenInfo {
		return toks[curTokIndex]
	}

	advance := func() TokenInfo {
		curTokIndex++
		return toks[curTokIndex-1]
	}

	consume := func(tokType TokenType) TokenInfo {
		if peek().tokType != tokType {
			panic("Error while compiling")
		}
		return advance()
	}

	addVarToSymTable := func(varIdent string) int {
		var varAddr uint32

		for _, s := range symTable {
			if s.symType == ST_VAR {
				if s.symScope == GLOBAL_SCOPE && curScope == GLOBAL_SCOPE {
					varAddr++
				} else if s.symScope > GLOBAL_SCOPE && curScope > GLOBAL_SCOPE {
					varAddr++
				}
			}
		}

		var s SymInfo
		s.symType = ST_VAR
		s.symIdent = varIdent
		s.symAddr = varAddr
		s.symScope = curScope
		s.symArgCount = 0

		symTable = append(symTable, s)
		return len(symTable) - 1
	}

	addFuncToSymTable := func(funcIdent string) int {
		var s SymInfo
		s.symType = ST_FUNC
		s.symIdent = funcIdent
		s.symAddr = uint32(len(bytecode))
		s.symScope = GLOBAL_SCOPE
		s.symArgCount = 0

		symTable = append(symTable, s)
		return len(symTable) - 1
	}

	delSymFromSymTable := func() int {
		var symRemovedCount int = 0
		var newSymTable []SymInfo
		for _, b := range symTable {
			if b.symScope <= curScope {
				newSymTable = append(newSymTable, b)
			} else {
				symRemovedCount++
			}
		}
		symTable = newSymTable
		return symRemovedCount
	}

	findSym := func(symIdent string) SymInfo {
		var curSym SymInfo
		curSym.symType = ST_ILLEGAL
		curSym.symScope = GLOBAL_SCOPE
		for _, s := range symTable {
			if s.symIdent == symIdent &&
				(curSym.symScope < s.symScope || curSym.symType == ST_ILLEGAL) {
				curSym = s
			}
		}
		return curSym
	}

	isTokBinaryOp := func(tok TokenInfo) bool {
		tokTypes := []TokenType{
			TT_ADD, TT_SUB, TT_MUL, TT_QUO, TT_REM, TT_AND, TT_OR, TT_XOR,
			TT_SHL, TT_SHR, TT_EQL, TT_LSS, TT_GTR, TT_NEQ, TT_LEQ, TT_GEQ,
		}
		for _, t := range tokTypes {
			if t == tok.tokType {
				return true
			}
		}
		return false
	}

	BinaryTokOpcode := map[TokenType]byte{
		TT_ADD: OP_ADD,
		TT_SUB: OP_SUB,
		TT_MUL: OP_MUL,
		TT_QUO: OP_QUO,
		TT_REM: OP_REM,
		TT_AND: OP_AND,
		TT_OR:  OP_OR,
		TT_XOR: OP_XOR,
		TT_SHL: OP_SHL,
		TT_SHR: OP_SHR,
		TT_EQL: OP_EQL,
		TT_LSS: OP_LSS,
		TT_GTR: OP_GTR,
		TT_NEQ: OP_NEQ,
		TT_LEQ: OP_LEQ,
		TT_GEQ: OP_GEQ,
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
		if peek().tokType == TT_MUL {
			var c int = 0
			for peek().tokType == TT_MUL {
				consume(TT_MUL)
				c++
			}
			compileGrouping()
			for c != 0 {
				emitInst(OP_GET_MEM)
				c--
			}
		} else if peek().tokType == TT_ADD {
			consume(TT_ADD)
			v, _ := strconv.ParseInt("+"+consume(TT_INT).tokStr, 0, 64)
			emitPushLitInst(uint32(v))
		} else if peek().tokType == TT_SUB {
			consume(TT_SUB)
			v, _ := strconv.ParseInt("-"+consume(TT_INT).tokStr, 0, 64)
			emitPushLitInst(uint32(v))
		} else {
			compileUnary(false)
		}
		consume(TT_RPAREN)
	}

	compileUnary = func(isRightOfBinary bool) {
		switch peek().tokType {
		case TT_IDENT:
			s := findSym(consume(TT_IDENT).tokStr)
			if s.symType == ST_FUNC {
				consume(TT_LPAREN)
				for peek().tokType != TT_RPAREN {
					compileExpr()
					if peek().tokType != TT_RPAREN {
						consume(TT_COMMA)
					}
				}
				consume(TT_RPAREN)
				emitPushLitInst(uint32(s.symArgCount))
				emitPushLitInst(s.symAddr)
				emitInst(OP_CALL)
			} else {
				emitPushLitInst(s.symAddr)
				if s.symScope == GLOBAL_SCOPE {
					emitInst(OP_GET_GLOBAL)
				} else {
					emitInst(OP_GET_LOCAL)
				}
			}

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

	emitPushLitInst(0)
	emitPushLitBlackInst()
	emitInst(OP_CALL)
	emitInst(OP_ECALL)

	addFuncToSymTable("ecall")
	emitInst(OP_ECALL)
	emitPushLitInst(0)
	emitInst(OP_RETURN)

	for peek().tokType != TT_EOF {
		switch peek().tokType {
		case TT_VAR:
			consume(TT_VAR)
			varIdent := consume(TT_IDENT).tokStr

			if curScope == GLOBAL_SCOPE {
				addVarToSymTable(varIdent)
			} else {
				emitPushLitInst(0)

				if peek().tokType == TT_ASSIGN {
					consume(TT_ASSIGN)
					compileExpr()
					emitPushLitInst(symTable[addVarToSymTable(varIdent)].symAddr)
					emitInst(OP_SET_LOCAL)
				} else {
					addVarToSymTable(varIdent)
				}
			}
			consume(TT_NEW_LINE)

		case TT_FUNC:
			consume(TT_FUNC)
			curScope++
			curBlocks = append(curBlocks, TT_FUNC)
			i := addFuncToSymTable(consume(TT_IDENT).tokStr)
			consume(TT_LPAREN)
			for peek().tokType != TT_RPAREN {
				addVarToSymTable(consume(TT_IDENT).tokStr)
				symTable[i].symArgCount++
				if peek().tokType != TT_RPAREN {
					consume(TT_COMMA)
				}
			}
			consume(TT_RPAREN)
			consume(TT_NEW_LINE)

		case TT_IDENT:
			s := findSym(peek().tokStr)
			if s.symType == ST_VAR {
				consume(TT_IDENT)
				consume(TT_ASSIGN)
				compileExpr()
				emitPushLitInst(s.symAddr)
				if s.symScope == GLOBAL_SCOPE {
					emitInst(OP_SET_GLOBAL)
				} else {
					emitInst(OP_SET_LOCAL)
				}
				consume(TT_NEW_LINE)
			} else {
				compileExpr()
				emitInst(OP_POP_LITERAL)
				consume(TT_NEW_LINE)
			}

		case TT_MUL:
			var c int = 0
			for peek().tokType == TT_MUL {
				consume(TT_MUL)
				c++
			}
			consume(TT_LPAREN)
			s := findSym(consume(TT_IDENT).tokStr)
			consume(TT_RPAREN)
			consume(TT_ASSIGN)
			if peek().tokType == TT_STR {
				emitPushLitInst(0)
				b := []byte(consume(TT_STR).tokStr)
				for len(b) > 0 {
					emitPushLitInst(uint32(b[len(b)-1]))
					b = b[:len(b)-1]
				}
				emitPushLitInst(s.symAddr)
				if s.symScope == GLOBAL_SCOPE {
					emitInst(OP_GET_GLOBAL)
				} else {
					emitInst(OP_GET_LOCAL)
				}
				c--
				for c != 0 {
					emitInst(OP_GET_MEM)
					c--
				}
				emitInst(OP_SET_MEM_STR)
			} else {
				compileExpr()
				emitPushLitInst(s.symAddr)
				if s.symScope == GLOBAL_SCOPE {
					emitInst(OP_GET_GLOBAL)
				} else {
					emitInst(OP_GET_LOCAL)
				}
				c--
				for c != 0 {
					emitInst(OP_GET_MEM)
					c--
				}
				emitInst(OP_SET_MEM)
			}

		case TT_IF:
			consume(TT_IF)
			curScope++
			curBlocks = append(curBlocks, TT_IF)
			emitPushLitBlackInst()
			compileExpr()
			emitInst(OP_BRANCH)
			consume(TT_NEW_LINE)

		case TT_WHILE:
			consume(TT_WHILE)
			curScope++
			curBlocks = append(curBlocks, TT_WHILE)
			loopAddr = append(loopAddr, uint32(len(bytecode)))
			emitPushLitBlackInst()
			compileExpr()
			emitInst(OP_BRANCH)
			consume(TT_NEW_LINE)

		case TT_RETURN:
			consume(TT_RETURN)
			compileExpr()
			emitInst(OP_RETURN)
			consume(TT_NEW_LINE)

		case TT_END:
			consume(TT_END)
			curScope--
			c := delSymFromSymTable()
			for c != 0 {
				emitInst(OP_POP_LITERAL)
				c--
			}

			switch curBlocks[len(curBlocks)-1] {
			case TT_FUNC:
				if curScope != GLOBAL_SCOPE {
					panic("Error while compiling")
				}
				emitPushLitInst(0)
				emitInst(OP_RETURN)

			case TT_IF:
				fillPushLitBlackInst(uint32(len(bytecode)))

			case TT_WHILE:
				emitPushLitInst(loopAddr[len(loopAddr)-1])
				loopAddr = loopAddr[:len(loopAddr)-1]
				emitInst(OP_JUMP)
				fillPushLitBlackInst(uint32(len(bytecode)))

			default:
				panic("Error while compiling")
			}

			curBlocks = curBlocks[:len(curBlocks)-1]
			consume(TT_NEW_LINE)

		case TT_NEW_LINE:
			consume(TT_NEW_LINE)
		default:
			panic("Error while compiling")
		}
	}

	fillPushLitBlackInst(findSym("main").symAddr)

	return bytecode
}
