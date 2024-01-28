package main

import (
	"log"
	"os"
)

func VMHandleECALL(vm VMState) VMState {
	var callType uint32 = vm.mem[4]
	var addr uint32 = vm.mem[5]
	var bufSize uint32 = vm.mem[6]

	switch callType {
	case 8:
		data, err := os.ReadFile(os.Args[2])

		if err != nil {
			log.Fatal(err)
		}

		data = append(data, 0x0a, 0x00)
		for len(data) != 0 && addr < bufSize {
			vm.mem[addr] = uint32(data[0])
			data = data[1:]
			addr++
		}
	}

	vm.status = VM_STATUS_RUNNING
	return vm
}
