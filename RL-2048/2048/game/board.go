package game

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"time"
)

var DebugLogLevel bool

const (
	_rows = 4
	_cols = 4

	// this is the sequence which is used to clear the screen :magic
	_clearScreenSequence = "\033[H\033[2J" // this works in mac. Might need other string for other OS

	probabilitySpace = 100
	probabilityOfTwo = 80 // probabilityOfTwo times 2 will come as new element out of  probabilitySpace1
)

type IBoard interface {
	Display()
	AddElement()
	TakeInput()
	IsOver() bool
	CountScore() (int, int)
	Move(dir Dir)
}

type SBoard struct {
	Matrix [][]int
	over   bool
	newRow int
	newCol int
}

func (b *SBoard) CountScore() (int, int) {
	total := 0
	maximum := 0
	matrix := b.Matrix
	for i := 0; i < _rows; i++ {
		for j := 0; j < _cols; j++ {
			total += matrix[i][j]
			maximum = max(maximum, matrix[i][j])
		}
	}
	return maximum, total
}

func max(one int, two int) int {
	if one > two {
		return one
	}
	return two
}

func (b *SBoard) IsOver() bool {
	empty := 0
	for i := 0; i < _rows; i++ {
		for j := 0; j < _cols; j++ {
			if b.Matrix[i][j] == 0 {
				empty++
			}
		}
	}
	return empty == 0 || b.over
}

func (b *SBoard) TakeInput() {
	var dir Dir
	dir, err := GetCharKeystroke()
	if err != nil {
		if errors.Is(err, errEndGame) {
			b.over = true
			return
		} else {
			log.Fatal("error while taking input for game: %v", err)
			return
		}
	}
	if DebugLogLevel {
		log.Printf("the dir is: %v \n", dir)
	}

	if dir == NO_DIR {
		// this makes pressing any keys other than Move-set doesn't make any change in the game
		b.TakeInput() // retry to get a valid direction
	}
	b.Move(dir)
}

// AddElement : it first finds the empty slots in the SBoard. They are the one with 0 value
// The it places a new cell randomly in one of those empty places
// The new value to put is also calculated randomly
func (b *SBoard) AddElement() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	val := r1.Int() % probabilitySpace
	if val <= probabilityOfTwo {
		val = 2
	} else {
		val = 4
	}

	empty := 0
	for i := 0; i < _rows; i++ {
		for j := 0; j < _cols; j++ {
			if b.Matrix[i][j] == 0 {
				empty++
			}
		}
	}
	elementCount := r1.Int()%empty + 1
	index := 0

	for i := 0; i < _rows; i++ {
		for j := 0; j < _cols; j++ {
			if b.Matrix[i][j] == 0 {
				index++
				if index == elementCount {
					b.newRow = i
					b.newCol = j
					b.Matrix[i][j] = val
					return
				}
			}
		}
	}
	return
}

// Display this is the method which draws the SBoard
// SBoard contains a Matrix which has cells. Each cell is a number.
// A Cell with 0 is considered empty
// to display number pretty, we make use of left and right padding
// Grid is formed using Ascii characters and some amount of test-&-see
func (b *SBoard) Display() {
	d := color.New(color.FgBlue, color.Bold)
	//b.Matrix = getRandom()
	fmt.Println(_clearScreenSequence)
	for i := 0; i < len(b.Matrix); i++ {
		printHorizontal()
		fmt.Printf("|")
		for j := 0; j < len(b.Matrix[0]); j++ {
			fmt.Printf("%3s", "")
			if b.Matrix[i][j] == 0 {
				fmt.Printf("%-6s|", "")
			} else {
				if i == b.newRow && j == b.newCol {
					d.Printf("%-6d|", b.Matrix[i][j])
				} else {
					fmt.Printf("%-6d|", b.Matrix[i][j])
				}
			}
		}
		fmt.Printf("%4s", "")
		fmt.Println()
	}
	printHorizontal()
}

// printHorizontal prints a grid row
func printHorizontal() {
	for i := 0; i < 40; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func New() IBoard {
	matrix := make([][]int, 0)
	for i := 0; i < _rows; i++ {
		matrix = append(matrix, make([]int, _cols))
	}
	return &SBoard{
		Matrix: matrix,
	}
}

var errEndGame = errors.New("GameOverError")
