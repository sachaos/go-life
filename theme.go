package main

import (
	"github.com/gdamore/tcell"
)

type Theme struct {
	Name       string
	BackGround tcell.Color
	Colors     []tcell.Color
	Levels     []int
}

var ThemeOcean = Theme{
	Name:       "Ocean",
	BackGround: tcell.NewRGBColor(0, 0, 100),
	Colors: []tcell.Color{
		tcell.NewRGBColor(50, 50, 255),
		tcell.NewRGBColor(100, 100, 255),
		tcell.NewRGBColor(150, 150, 255),
		tcell.NewRGBColor(255, 255, 255),
	},
	Levels: []int{1, 1, 1, 1},
}

var ThemeFire = Theme{
	Name:       "Fire",
	BackGround: tcell.NewRGBColor(100, 0, 0),
	Colors: []tcell.Color{
		tcell.NewRGBColor(255, 50, 50),
		tcell.NewRGBColor(255, 100, 100),
		tcell.NewRGBColor(255, 150, 150),
		tcell.NewRGBColor(255, 255, 255),
	},
	Levels: []int{1, 1, 1, 1},
}

var ThemeMatrix = Theme{
	Name:       "Matrix",
	BackGround: tcell.NewRGBColor(0, 0, 0),
	Colors: []tcell.Color{
		tcell.NewRGBColor(50, 255, 50),
	},
	Levels: []int{1},
}

var ThemeBlackAndWhite = Theme{
	Name:       "Black and White",
	BackGround: tcell.NewRGBColor(0, 0, 0),
	Colors: []tcell.Color{
		tcell.NewRGBColor(255, 255, 255),
	},
	Levels: []int{1},
}

var ThemeWhiteAndBlack = Theme{
	Name:       "White and Black",
	BackGround: tcell.NewRGBColor(255, 255, 255),
	Colors: []tcell.Color{
		tcell.NewRGBColor(0, 0, 0),
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
