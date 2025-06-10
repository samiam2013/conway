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
	defer func() {
		// hold the terminal open waiting for input
		_, _ = fmt.Scanf("%s")
	}()
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

}

// this is just a helper function to get started
func drawCheckerBoard(width, height int) string {
	lines := strings.Builder{}
	for h := 0; h < height; h++ {
		char := " "
		if h%2 == 0 {
			char = "█"
		}
		if _, err := lines.Write([]byte(char)); err != nil {
			panic(fmt.Sprintf("failed to draw checkerboard leading char: %s", err))
		}
		// width - 1 because making the first line alternate takes one away
		for w := 0; w < (width - 1); w++ {
			if char == " " {
				char = "█"
			} else {
				char = " "
			}
			if _, err := lines.Write([]byte(char)); err != nil {
				panic(fmt.Sprintf("failed to write checkerboard trailing char: %s", err))
			}
		}
		if _, err := lines.Write([]byte("\n")); err != nil {
			panic(fmt.Sprintf("failed to write checkerboard newline char: %s", err))
		}
	}
	return strings.TrimRight(lines.String(), "\n")
}
