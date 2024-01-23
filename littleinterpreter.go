package main

import (
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	data = append(data, 0x0a)

	vm := VMCreate(GenerateBytecode(GenerateTokens(data)))
	VMIsDebuggerOn = true
	VMRun(vm)
}
