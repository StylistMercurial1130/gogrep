package main

import (
	"fmt"
	"gogrep/statemachine"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("invalid usage of command")
		return
	}

	expr := args[0]
	word := args[1]

	status := statemachine.Evaluate(expr, word)

	if status {
		fmt.Printf("is match !")
	} else {
		fmt.Printf("is not match !")
	}
}
