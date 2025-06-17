package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"golang.org/x/term"
)

// [height][width]
type world struct {
	Places [][]bool
}

func (w *world) Width() int {
	if len(w.Places) == 0 {
		return 0
	}
	return len(w.Places[0])
}

func (w *world) Height() int {
	return len(w.Places)
}

func NewWorld(height, width int) world {
	new := make([][]bool, height)
	for h := range height {
		new[h] = make([]bool, width)
		for w := range width {
			new[h][w] = false
		}
	}
	return world{Places: new}
}

func (w *world) SetCoord(x, y int, val bool) error {
	// convert the cartesian offset to the indexes of the slice
	// set that value
	if w == nil || w.Height() <= 0 || w.Width() <= 0 {
		return errors.New("w must be non-nil height and width must be > 0")
	}
	// TODO can this be accomplished without all the dumb casting?
	halfW := int(math.Floor(float64(w.Width()) / 2.0))
	halfH := int(math.Floor(float64(w.Height()) / 2.0))

	placeW := halfW + x
	placeH := halfH + y

	if placeW > w.Width() {
		return errors.New("x coord must be between +- 1/2 * width")
	}
	if placeH > w.Height() {
		return errors.New("y coord must be between +- 1/2 * height")
	}

	w.Places[placeH][placeW] = true
	return nil
}

func (w *world) String() string {
	sb := strings.Builder{}
	for _, row := range w.Places {
		for _, val := range row {
			c := " "
			if val {
				c = "█"
			}
			if _, err := sb.WriteString(c); err != nil {
				return "error writing value in world string building"
			}
		}
		if _, err := sb.WriteString("\n"); err != nil {
			return "error writing newline in world string building"
		}
	}
	return sb.String()
}

func (w *world) Evolve() error {
	// for each space inside the map, check the value of all the spaces around
	// < 2 neighbors dies
	// 2 - 3 neighbors survives
	// > 3 neighbors dies
	newWorld := NewWorld(w.Height(), w.Width())
	for hRow := range len(w.Places) {
		for wCol := range len(w.Places[0]) {
			count, err := w.CountNeighbors(hRow, wCol)
			if err != nil {
				return errors.Join(errors.New("failed to count neighbors for evolution"), err)
			}
			val := w.Places[hRow][wCol]
			if val {
				switch {
				case count < 2:
					newWorld.Places[hRow][wCol] = false
				case count == 2 || count == 3:
					newWorld.Places[hRow][wCol] = true
				default:
					newWorld.Places[hRow][wCol] = false
				}
			} else if !val && count == 3 {
				newWorld.Places[hRow][wCol] = true
			}
		}
	}
	*w = newWorld
	return nil
}

func (w *world) CountNeighbors(height, width int) (int, error) {
	if height < 0 || height > len(w.Places) {
		return 0, errors.New("height for count out of bounds")
	}
	if width < 0 || (len(w.Places) > 0 && width > len(w.Places[0])) {
		return 0, errors.New("width out of bounds for count")
	}

	count := 0
	// if there's a row above
	if height > 0 {
		// check the three spaces
		if width > 0 {
			// there's one to the top left
			if w.Places[height-1][width-1] {
				count++
			}
		}
		if w.Places[height-1][width] {
			// there's one in the top middle
			count++
		}
		if width < (w.Width() - 1) {
			// there's one in the top right
			if w.Places[height-1][width+1] {
				count++
			}
		}
	}
	// if there's one to the left
	if width > 0 {
		if w.Places[height][width-1] {
			count++
		}
	}
	// TODO: do I count (shouldn't)
	// if w.Places[height][width] {
	// 	count++
	// }
	// if there's one to the right
	if width < (w.Width() - 1) {
		if w.Places[height][width+1] {
			count++
		}
	}
	// if there's a row below
	if height < (w.Height() - 1) {
		// check the three spaces
		if width > 0 {
			// there's one to the bottom left
			if w.Places[height+1][width-1] {
				count++
			}
		}
		if w.Places[height+1][width] {
			// there's one in the bottom middle
			count++
		}
		if width < (w.Width() - 1) {
			// there's one in the bottom right
			if w.Places[height+1][width+1] {
				count++
			}
		}
	}
	return count, nil
}

func main() {
	// fmt.Println("Press enter at the end to exit the program")
	// defer func() {
	// 	// hold the terminal open waiting for input
	// 	_, _ = fmt.Scanf("%s")
	// }()
	// time.Sleep(time.Second * 1)
	// get the width and height of the terminal
	width, height, err := term.GetSize(0)
	if err != nil {
		slog.Error("failed to get terminal size", "error", err)
		return
	}
	slog.Info("dimensions", "width", width, "height", height)
	// draw blocks every other column
	// lines := drawCheckerBoard(width, height)
	// fmt.Print(lines)

	slog.Info("building new world")
	newWorld := NewWorld(height, width)
	type inputRow struct {
		x   int
		y   int
		val bool
	}
	inputs := []inputRow{
		{0, 0, true},
		{0, -1, true},
		{1, -2, true},
		{1, 0, true},
		{2, 0, true}}
	slog.Info("starting loop")
	for _, input := range inputs {
		if err := newWorld.SetCoord(input.x, input.y, input.val); err != nil {
			slog.Error("failed to set coordinate", "error", err)
			return
		}
	}

	err = nil
	for err == nil {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(newWorld.String())
		err = newWorld.Evolve()
	}
	slog.Error("failed to evolve the world", "error", err)
}

// this is just a helper function to get started
// func drawCheckerBoard(width, height int) string {
// 	lines := strings.Builder{}
// 	for h := range height {
// 		char := " "
// 		if h%2 == 0 {
// 			char = "█"
// 		}
// 		if _, err := lines.Write([]byte(char)); err != nil {
// 			panic(fmt.Sprintf("failed to draw checkerboard leading char: %s", err))
// 		}
// 		// width - 1 because making the first line alternate takes one away
// 		for range width - 1 {
// 			if char == " " {
// 				char = "█"
// 			} else {
// 				char = " "
// 			}
// 			if _, err := lines.Write([]byte(char)); err != nil {
// 				panic(fmt.Sprintf("failed to write checkerboard trailing char: %s", err))
// 			}
// 		}
// 		if _, err := lines.Write([]byte("\n")); err != nil {
// 			panic(fmt.Sprintf("failed to write checkerboard newline char: %s", err))
// 		}
// 	}
// 	return strings.TrimRight(lines.String(), "\n")
// }
