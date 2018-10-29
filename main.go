package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/sachaos/go-life/preset"
	"github.com/urfave/cli"
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

func startGame(themes []Theme, presets []preset.Preset, themeIndex int) error {
	rand.Seed(time.Now().Unix())

	s := initScreen()
	defer s.Fini()

	// init board
	width, height := s.Size()
	b := NewBoard(height, width/2)

	b.Random()

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
	app.Version = "0.2.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
		cli.StringFlag{
			Name:  "theme",
			Value: "BlackAndWhite",
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
		themeIndex := -1
		for i, theme := range themes {
			if theme.Name == c.String("theme") {
				themeIndex = i
				break
			}
		}
		if themeIndex == -1 {
			return fmt.Errorf("Invalid theme name: %s", c.String("theme"))
		}
		return startGame(themes, presets, themeIndex)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
