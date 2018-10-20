package main

import (
	"fmt"
	"math/rand"
)

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

func Cells(width, height int) [][]Cell {
	cells := make([][]Cell, height)
	for i := 0; i < height; i++ {
		cells[i] = make([]Cell, width)
	}
	return cells
}

func (b *Board) Init() {
	b.grid = Cells(b.width, b.height)
	b.link()
}

func (b *Board) link() {
	around := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	// link
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			b.grid[i][j].UnLink()
			aroundCells := []*Cell{}
			for _, a := range around {
				y := i + a[0]
				x := j + a[1]

				if x < 0 || y < 0 || x >= b.width || y >= b.height {
					continue
				}

				aroundCells = append(aroundCells, &b.grid[y][x])
			}
			b.grid[i][j].Link(aroundCells)
		}
	}
}

func (b *Board) Resize(width, height int) {
	newGrid := Cells(width, height)
	var minHeight int
	if b.height > height {
		minHeight = height
	} else {
		minHeight = b.height
	}
	for i := 0; i < minHeight; i++ {
		copy(newGrid[i], b.grid[i])
	}
	b.grid = newGrid
	b.width = width
	b.height = height

	b.link()
}

func (b *Board) State() [][]Cell {
	return b.grid
}

func (b *Board) Random() {
	grid := make([][]bool, b.height)
	for i := 0; i < b.height; i++ {
		grid[i] = make([]bool, b.width)
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			grid[i][j] = rand.Int()%2 == 0
		}
	}
	b.Set(0, 0, grid)
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
