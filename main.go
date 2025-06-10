package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/term"
)

// [height][width]
type world [][]place

type place bool

func NewWorld(height, width int) world {
	new := [][]place{}
	for h := range height {
		for w := range width {
			new[h][w] = false
		}
	}
	return new
}

func (w *world) SetCoord(x, y int, val bool) {
	// convert the cartesian offset to the indexes of the slice
	// set that value
}

func main() {
	fmt.Println("Press enter at the end to exit the program")
	time.Sleep(time.Second * 1)
	// get the width and height of the terminal
	width, height, err := term.GetSize(0)
	if err != nil {
		slog.Error("failed to get terminal size", "error", err)
		return
	}
	fmt.Printf("width: %d height: %d\n", width, height)
	// draw blocks every other column
	lines := drawCheckerBoard(width, height)

	fmt.Print(lines)

	// hold the terminal open waiting for input
	_, _ = fmt.Scanf("%s")
}

// this is just a helper function to get started
func drawCheckerBoard(width, height int) string {
	lines := ""
	for h := 0; h < height; h++ {
		char := " "
		if h%2 == 0 {
			char = "█"
		}
		lines += char
		// width - 1 because making the first line alternate takes one away
		for w := 0; w < (width - 1); w++ {
			if char == " " {
				char = "█"
			} else {
				char = " "
			}
			lines += char
		}
		lines += "\n"
	}
	lines = strings.TrimRight(lines, "\n")
	return lines
}
