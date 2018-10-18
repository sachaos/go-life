package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

type Cell struct {
	state     bool
	nextState bool

	aroundCells []*Cell
}

func (c *Cell) State() bool {
	return c.state
}

func (c *Cell) Switch() {
	c.state = !c.state
}

func (c *Cell) Set(state bool) {
	c.state = state
}

func (c *Cell) CalcNextState() {
	var nextState bool
	count := 0
	for _, ac := range c.aroundCells {
		if ac.state {
			count++
		}
	}

	if c.state {
		if count <= 1 || count >= 4 {
			nextState = false
		} else {
			nextState = true
		}
	} else {
		if count == 3 {
			nextState = true
		} else {
			nextState = false
		}
	}
	c.nextState = nextState
}

func (c *Cell) Flush() {
	c.state = c.nextState
}

type Board struct {
	height int
	width  int

	grid [][]Cell
	// cells []*Cell
}

func NewBoard(height, width int) *Board {
	board := Board{
		height: height,
		width:  width,
	}

	return &board
}

func (b *Board) Init() {
	b.grid = make([][]Cell, b.height)
	for i := 0; i < b.height; i++ {
		b.grid[i] = make([]Cell, b.width)
	}

	around := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	// link
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			for _, a := range around {
				y := i + a[0]
				x := j + a[1]

				if x < 0 || y < 0 || x >= b.width || y >= b.height {
					continue
				}

				b.grid[y][x].aroundCells = append(b.grid[y][x].aroundCells, &b.grid[i][j])
			}
		}
	}
}

func (b *Board) State() [][]bool {
	g := make([][]bool, b.height)
	for i := 0; i < b.height; i++ {
		g[i] = make([]bool, b.width)
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			g[i][j] = b.grid[i][j].State()
		}
	}

	return g
}

func (b *Board) Next() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			b.grid[i][j].CalcNextState()
		}
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			b.grid[i][j].Flush()
		}
	}
}

func (b *Board) Print() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.grid[i][j].state {
				fmt.Printf("O")
			} else {
				fmt.Printf("X")
			}
		}
		fmt.Printf("\n")
	}
}

func (b *Board) Set(x, y int, bgrid [][]bool) {
	for i, row := range bgrid {
		for j, state := range row {
			b.grid[y+i][x+j].Set(state)
		}
	}
}

func main() {
	// init screen
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()

	width, height := s.Size()

	// init board
	b := NewBoard(height, width/2)
	b.Init()
	b.Set(1, 1, [][]bool{
		{false, false, false, true, true, false, false, false},
		{false, false, true, false, false, true, false, false},
		{false, true, false, false, false, false, true, false},
		{true, false, false, false, false, false, false, true},
		{true, false, false, false, false, false, false, true},
		{false, true, false, false, false, false, true, false},
		{false, false, true, false, false, true, false, false},
		{false, false, false, true, true, false, false, false},
	})

	// b.Set(80, 40, [][]bool{
	// 	{false, true, false, false, false, false, false, false},
	// 	{false, false, false, true, false, false, false, false},
	// 	{true, true, false, false, true, true, true, false},
	// })

	for {
		s.Clear()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		for i, row := range b.State() {
			for j, state := range row {
				if state {
					s.SetCell(j*2, i, st, ' ')
					s.SetCell(j*2+1, i, st, ' ')
				}
			}
		}
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			continue
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				b.Next()
			} else {
				s.Fini()
				os.Exit(0)
			}
		default:
			continue
		}
	}
}
