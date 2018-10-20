package main

import (
	"github.com/gdamore/tcell"
)

type Theme struct {
	Colors []tcell.Color
	Levels []int
}

var ThemeBlue = Theme{
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

var ThemeBlackAndWhite = Theme{
	Colors: []tcell.Color{
		tcell.NewRGBColor(255, 255, 255),
	},
	Levels: []int{1},
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
