package rle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

func newCells(width, height int) [][]bool {
	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}
	return grid
}

func ParseFile(path string) ([][]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return [][]bool{}, err
	}
	return Parse(f), nil
}

type RunLength struct {
	RunCount  int
	Tag       rune
	IsNewLine bool
}

func Parse(r io.Reader) [][]bool {
	br := bufio.NewReader(r)
	var x, y int
	var header string

	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if line[0] == '#' {
			continue
		}

		if line[0] == 'x' {
			header = string(line)
			break
		}
	}

	_, err := fmt.Sscanf(header, "x = %d, y = %d", &x, &y)
	if err != nil {
		panic(err)
	}

	cells := newCells(x, y)

	var runLengthes []RunLength

	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if line[0] == '#' {
			continue
		}

		var parsingDigit string
		for _, r := range string(line) {
			if r == '!' {
				break
			}

			if r == '$' {
				runLengthes = append(runLengthes, RunLength{IsNewLine: true})
				continue
			}

			if unicode.IsDigit(r) {
				parsingDigit = parsingDigit + string(r)
				continue
			}

			parsedInt := 1
			if len(parsingDigit) > 0 {
				parsedInt, _ = strconv.Atoi(parsingDigit)
			}
			parsingDigit = ""

			runLengthes = append(runLengthes, RunLength{RunCount: parsedInt, Tag: r})
		}
	}

	var i, j int
	for _, runLength := range runLengthes {
		if runLength.IsNewLine {
			j++
			i = 0
			continue
		}

		if runLength.Tag == 'o' {
			for p := 0; p < runLength.RunCount; p++ {
				cells[j][i+p] = true
			}
		}
		i = i + runLength.RunCount
	}

	return cells
}
