package main

import (
	"fmt"
	"os"

	"SmokeLab/packages/engine"
)

func main() {
	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	fmt.Println(engine.NewGreetingService().Greet(name))
}
