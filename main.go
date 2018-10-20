package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

type Theme struct {
	Colors []tcell.Color
	Levels []int
}

var theme = Theme{
	Colors: []tcell.Color{
		tcell.NewRGBColor(50, 50, 255),
		tcell.NewRGBColor(100, 100, 255),
		tcell.NewRGBColor(125, 125, 255),
		tcell.NewRGBColor(150, 150, 255),
		tcell.NewRGBColor(175, 175, 255),
		tcell.NewRGBColor(205, 205, 255),
		tcell.NewRGBColor(230, 230, 255),
		tcell.NewRGBColor(255, 255, 255),
	},
	Levels: []int{1, 1, 1, 1, 1, 1, 1},
}

func (t *Theme) Color(time int) tcell.Color {
	total := 0
	for i, l := range t.Levels {
		total += l
		if time < total {
			return t.Colors[i]
		}
	}
	return t.Colors[len(t.Colors)-1]
}

type Cell struct {
	state     bool
	nextState bool
	liveTime  int

	aroundCells []*Cell
}

func (c *Cell) State() bool {
	return c.state
}

func (c *Cell) LiveTime() int {
	return c.liveTime
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
	if c.state && c.nextState {
		c.liveTime++
	} else {
		c.liveTime = 0
	}
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

func (b *Board) State() [][]Cell {
	return b.grid
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
			if y+i < b.height && x+j < b.width {
				b.grid[y+i][x+j].Set(state)
			}
		}
	}
}

func (b *Board) Get(x, y int) *Cell {
	return &b.grid[y][x]
}

func main() {
	rand.Seed(time.Now().Unix())
	stop := false
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
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()

	width, height := s.Size()

	// init board
	b := NewBoard(height, width/2)
	b.Init()
	// b.Set(1, 1, [][]bool{
	// 	{false, false, false, true, true, false, false, false},
	// 	{false, false, true, false, false, true, false, false},
	// 	{false, true, false, false, false, false, true, false},
	// 	{true, false, false, false, false, false, false, true},
	// 	{true, false, false, false, false, false, false, true},
	// 	{false, true, false, false, false, false, true, false},
	// 	{false, false, true, false, false, true, false, false},
	// 	{false, false, false, true, true, false, false, false},
	// })

	b.Set(80, 40, [][]bool{
		{false, true, false, false, false, false, false, false},
		{false, false, false, true, false, false, false, false},
		{true, true, false, false, true, true, true, false},
	})

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan struct{}, 0)
	stopSwtich := make(chan struct{}, 0)
	reset := make(chan struct{}, 0)
	step := make(chan struct{}, 0)
	clear := make(chan struct{}, 0)
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					x, y := ev.Position()
					b.Get(x/2, y).Switch()
				default:
					continue
				}
			case *tcell.EventResize:
				continue
			case *tcell.EventKey:
				if ev.Rune() == 'r' {
					reset <- struct{}{}
				} else if ev.Key() == tcell.KeyEnter {
					step <- struct{}{}
				} else if ev.Rune() == ' ' {
					stopSwtich <- struct{}{}
				} else if ev.Key() == tcell.KeyEsc {
					done <- struct{}{}
				} else if ev.Rune() == 'c' {
					clear <- struct{}{}
				}
			default:
				continue
			}
		}
	}()

	for {
		s.Clear()
		for i, row := range b.State() {
			for j, cell := range row {
				st := tcell.StyleDefault.Background(theme.Color(cell.LiveTime()))
				if cell.State() {
					s.SetCell(j*2, i, st, ' ')
					s.SetCell(j*2+1, i, st, ' ')
				}
			}
		}

		select {
		case <-reset:
			stop = !stop
			b.Init()
			grid := make([][]bool, height)
			for i := 0; i < height; i++ {
				grid[i] = make([]bool, width)
			}

			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {
					grid[i][j] = rand.Int()%2 == 0
				}
			}
			b.Set(0, 0, grid)
			stop = !stop
		case <-stopSwtich:
			stop = !stop
		case <-step:
			b.Next()
			s.Show()
		case <-done:
			s.Fini()
			os.Exit(0)
		case <-clear:
			b.Init()
			s.Show()
		case <-ticker.C:
			if !stop {
				b.Next()
			}
			s.Show()
		}
	}
}
