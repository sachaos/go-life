package main

import (
	"fmt"
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
