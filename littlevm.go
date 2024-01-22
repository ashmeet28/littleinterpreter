package main

// import "fmt"

// const (
// 	VM_STATUS_UNKNOWN int = iota
// 	VM_STATUS_READY
// 	VM_STATUS_RUNNING
// 	VM_STATUS_HALT
// 	VM_STATUS_ERROR
// )

var (
	OP_NOP   byte = 1
	OP_ECALL byte = 2

	OP_ADD  byte = 4  // +
	OP_SUB  byte = 5  // -
	OP_MUL  byte = 6  // *
	OP_QUO  byte = 7  // /
	OP_REM  byte = 8  // %
	OP_AND  byte = 9  // &
	OP_OR   byte = 10 // |
	OP_XOR  byte = 11 // ^
	OP_SHL  byte = 12 // <<
	OP_SHR  byte = 13 // >>
	OP_LAND byte = 14 // &&
	OP_LOR  byte = 15 // ||
	OP_EQL  byte = 16 // ==
	OP_LSS  byte = 17 // <
	OP_GTR  byte = 18 // >
	OP_NEQ  byte = 19 // !=
	OP_LEQ  byte = 20 // <=
	OP_GEQ  byte = 21 // >=

	OP_PUSH_LITERAL byte = 32
	OP_POP_LITERAL  byte = 33

	OP_LOAD_GLOBAL  byte = 34
	OP_STORE_GLOBAL byte = 35

	OP_LOAD_LOCAL  byte = 36
	OP_STORE_LOCAL byte = 37

	OP_STORE_MEM byte = 38
	OP_LOAD_MEM  byte = 39

	OP_JUMP   byte = 40
	OP_BRANCH byte = 41

	OP_CALL   byte = 42
	OP_RETURN byte = 43
)

// type VMState struct {
// 	bytecode []byte

// 	pc  uint32   // Program counter
// 	g   []uint32 // Global variables
// 	s   []uint32 // Local variables and literals stack
// 	sp  uint32   // Stack pointer
// 	sfp uint32   // Stack frame pointer
// 	rs  []uint32 // Return stack
// 	rsp uint32   // Return stack pointer
// 	rv  uint32   // Return value
// 	mem []uint32 // Virtual Memory

// 	status int
// }

// var VMIsDebuggerOn bool = false

// func VMExecInst(vm VMState) VMState {
// 		var op byte = vm.bytecode[vm.pc]

// 		switch op {
// 		case OP_NOP:
// 			vm.pc++

// 		case OP_ECALL:
// 			vm.pc++
// 			vm.status = VM_STATUS_HALT

// 		case OP_ADD:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] + vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_SUB:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] - vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_AND:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] & vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_OR:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] | vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_XOR:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] ^ vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_SHL:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] << vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_SHR:
// 			vm.s[vm.sp-2] = vm.s[vm.sp-2] >> vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_PUSH_LITERAL:
// 			vm.s[vm.sp] = uint32(vm.byteCode[vm.pc+1]) | (uint32(vm.byteCode[vm.pc+2]) << 8) | (uint32(vm.byteCode[vm.pc+3]) << 16) | (uint32(vm.byteCode[vm.pc+4]) << 24)
// 			vm.sp++
// 			vm.pc += 5

// 		case OP_PUSH_LOCAL:
// 			vm.s[vm.sp-1] = vm.s[vm.s[vm.sp-1]+vm.sfp]
// 			vm.pc++

// 		case OP_PUSH_GLOBAL:
// 			vm.s[vm.sp-1] = vm.g[vm.s[vm.sp-1]]
// 			vm.pc++

// 		case OP_PUSH_FUNC_ARG:
// 			vm.s[vm.sp] = vm.fas[vm.fasp-1]
// 			vm.sp++
// 			vm.fasp--
// 			vm.pc++

// 		case OP_PUSH_FUNC_RET_VAL:
// 			vm.s[vm.sp] = vm.frv
// 			vm.sp++
// 			vm.pc++

// 		case OP_POP_LITERAL:
// 			vm.sp--
// 			vm.pc++

// 		case OP_POP_LOCAL:
// 			vm.s[vm.s[vm.sp-2]+vm.sfp] = vm.s[vm.sp-1]
// 			vm.sp -= 2
// 			vm.pc++

// 		case OP_POP_GLOBAL:
// 			vm.g[vm.s[vm.sp-2]] = vm.s[vm.sp-1]
// 			vm.sp -= 2
// 			vm.pc++

// 		case OP_POP_FUNC_ARG:
// 			vm.fas[vm.fasp] = vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.fasp++
// 			vm.pc++

// 		case OP_POP_FUNC_RET_VAL:
// 			vm.frv = vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.pc++

// 		case OP_EQ:
// 			if vm.s[vm.sp-2] == vm.s[vm.sp-1] {
// 				vm.s[vm.sp-2] = 1
// 			} else {
// 				vm.s[vm.sp-2] = 0
// 			}

// 			vm.sp--
// 			vm.pc++

// 		case OP_NE:
// 			if vm.s[vm.sp-2] != vm.s[vm.sp-1] {
// 				vm.s[vm.sp-2] = 1
// 			} else {
// 				vm.s[vm.sp-2] = 0
// 			}

// 			vm.sp--
// 			vm.pc++

// 		case OP_LT:
// 			if vm.s[vm.sp-2] < vm.s[vm.sp-1] {
// 				vm.s[vm.sp-2] = 1
// 			} else {
// 				vm.s[vm.sp-2] = 0
// 			}

// 			vm.sp--
// 			vm.pc++

// 		case OP_GE:
// 			if vm.s[vm.sp-2] >= vm.s[vm.sp-1] {
// 				vm.s[vm.sp-2] = 1
// 			} else {
// 				vm.s[vm.sp-2] = 0
// 			}

// 			vm.sp--
// 			vm.pc++

// 		case OP_JUMP:
// 			if vm.s[vm.sp-1] == 0 {
// 				vm.pc = vm.s[vm.sp-2]
// 			} else {
// 				vm.pc++
// 			}
// 			vm.sp -= 2

// 		case OP_CALL:
// 			vm.rs[vm.rsp] = vm.pc + 1
// 			vm.rs[vm.rsp+1] = vm.sfp

// 			vm.rsp += 2

// 			vm.pc = vm.s[vm.sp-1]
// 			vm.sp--
// 			vm.sfp = vm.sp

// 		case OP_RETURN:
// 			vm.sp = vm.sfp

// 			vm.sfp = vm.rs[vm.rsp-1]
// 			vm.pc = vm.rs[vm.rsp-2]
// 			vm.rsp -= 2

// 		default:
// 			vm.status = VM_STATUS_ERROR
// 		}

// 		if VMIsDebuggerOn {
// 			fmt.Println(op)

// 			fmt.Println("pc", vm.pc)

// 			fmt.Println("g", vm.g[:32])

// 			fmt.Println("s", vm.s[:32])
// 			fmt.Println("sp", vm.sp)
// 			fmt.Println("fp", vm.sfp)

// 			fmt.Println("rs", vm.rs[:32])
// 			fmt.Println("rsp", vm.rsp)
// 			fmt.Println("rv", vm.rv)

// 			fmt.Println("-")
// 			fmt.Println("-")
// 			fmt.Println("-")
// 			fmt.Println("-")
// 		}

// 	return vm
// }

// func VMRun(vm VMState) {
// 	if vm.status == VM_STATUS_READY {
// 		vm.status = VM_STATUS_RUNNING
// 	}

// 	for vm.status == VM_STATUS_RUNNING {
// 		vm = VMExecInst(vm)

// 		if vm.status == VM_STATUS_ERROR {
// 			fmt.Println("VM STATUS: ERROR")
// 		}
// 	}

// 	fmt.Println(vm.rv)
// }

// func VMCreate(bytecode []byte) VMState {
// 	var vm VMState

// 	for _, b := range bytecode {
// 		vm.bytecode = append(vm.bytecode, b)
// 	}

// 	vm.pc = 0

// 	vm.g = make([]uint32, 16777216)

// 	vm.s = make([]uint32, 16777216)
// 	vm.sp = 0
// 	vm.sfp = 0

// 	vm.rsp = 0
// 	vm.rs = make([]uint32, 16777216)
// 	vm.rv = 0

// 	vm.status = VM_STATUS_READY

// 	return vm
// }
