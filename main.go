package main

import (
	"fmt"
	"initiative/internal/ui"
	"os"
)

func main() {
	p := ui.NewProgram()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
