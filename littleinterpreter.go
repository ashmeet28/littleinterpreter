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
	VMIsDebuggerOn = true
	VMRun(VMCreate(GenerateBytecode(GenerateTokens(data))))
}
