package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/sachaos/go-life/format/rle"
	"github.com/sachaos/go-life/preset"
	"github.com/urfave/cli"
	"io"
	"log"
)

func initScreen() tcell.Screen {
	// init screen
	encoding.Register()
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	s.EnableMouse()

	return s
}

func startGame(themes []Theme, presets []preset.Preset, themeIndex int, defaultCells [][]bool) error {
	rand.Seed(time.Now().Unix())

	s := initScreen()
	defer s.Fini()

	// init board
	width, height := s.Size()
	b := NewBoard(height, width/2)

	if len(defaultCells) == 0 {
		b.Random()
	} else {
		pheight := len(defaultCells)
		pwidth := len(defaultCells[0])
		if pwidth > width || pheight > height {
			return fmt.Errorf("Specified pattern is too big\n")
		}

		b.Set((width/2-pwidth)/2, (height-pheight)/2, defaultCells)
	}

	// init ticker
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	event := make(chan Event)

	game := Game{
		screen:     s,
		board:      b,
		themes:     themes,
		presets:    presets,
		ticker:     ticker,
		event:      event,
		themeIndex: themeIndex,
	}

	go inputLoop(s, event)

	return game.Loop()
}

func listPresets(c *cli.Context, presets []preset.Preset) error {
	for _, preset := range presets {
		fmt.Println(preset.Name)
	}
	return nil
}

func listThemes(c *cli.Context, themes []Theme) error {
	for _, theme := range themes {
		fmt.Println(theme.Name)
	}
	return nil
}

func main() {
	presets, err := preset.LoadPresets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	themes := []Theme{
		ThemeBlackAndWhite,
		ThemeOcean,
		ThemeFire,
		ThemeMatrix,
		ThemeWhiteAndBlack,
	}

	app := cli.NewApp()
	app.Name = "go-life"
	app.Usage = "Conway's Game of Life"
	app.Version = "0.3.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
		cli.StringFlag{
			Name:  "theme",
			Value: "BlackAndWhite",
			Usage: "Theme name",
		},
		cli.StringFlag{
			Name:  "pattern",
			Usage: "Pattern name (e.g. glider, glider-gun)",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "Pattern file",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetOutput(os.Stderr)
		} else {
			file, err := os.Open(os.DevNull)
			if err != nil {
				return err
			}
			log.SetOutput(file)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "themes",
			Usage: "list themes",
			Action: func(c *cli.Context) error {
				return listThemes(c, themes)
			},
		},
		{
			Name:  "presets",
			Usage: "list presets",
			Action: func(c *cli.Context) error {
				return listPresets(c, presets)
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("pattern") != "" && c.String("file") != "" {
			return fmt.Errorf("Using pattern and file option is not permitted")
		}

		themeIndex := -1
		for i, theme := range themes {
			if theme.Name == c.String("theme") {
				themeIndex = i
				break
			}
		}
		if themeIndex == -1 {
			return fmt.Errorf("Invalid theme name: %s\n", c.String("theme"))
		}

		var defaultCells [][]bool
		specifiedPattern := c.String("pattern")
		if specifiedPattern != "" {
			for _, p := range presets {
				if p.Name == specifiedPattern {
					defaultCells = p.Cells
					break
				}
			}
			if len(defaultCells) == 0 {
				return fmt.Errorf("Invalid pattern name: %s\n", specifiedPattern)
			}
		}

		fileName := c.String("file")
		if fileName != "" {
			var file io.Reader
			if fileName == "-" {
				file = os.Stdin
			} else {
				file, err = os.Open(fileName)
				if err != nil {
					return err
				}
			}

			defaultCells = rle.Parse(file)
		}

		return startGame(themes, presets, themeIndex, defaultCells)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
