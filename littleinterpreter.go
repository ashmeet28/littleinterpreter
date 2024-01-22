package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	data = append(data, 0x0a)

	fmt.Println(GenerateBytecode(GenerateTokens(data)))
}
