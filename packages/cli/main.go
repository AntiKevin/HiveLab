package main

import (
	"fmt"
	"os"

	"SmokeLab/packages/core"
)

func main() {
	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	fmt.Println(core.NewGreetingService().Greet(name))
}
