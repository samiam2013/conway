package main

import (
	"fmt"
	"log/slog"

	"golang.org/x/term"
)

func main() {
	// get the width and height of the terminal
	width, height, err := term.GetSize(0)
	if err != nil {
		slog.Error("failed to get terminal size", "error", err)
		return
	}
	fmt.Printf("width: %d height: %d\n", width, height)
	// draw blocks every other column
	i := 0
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if (h+i)%2 == 0 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
			i++
		}
		fmt.Println()
	}
}
