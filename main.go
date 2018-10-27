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

func startGame(themes []Theme, presets []preset.Preset) error {
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
		screen:  s,
		board:   b,
		themes:  themes,
		presets: presets,
		ticker:  ticker,
		event:   event,
	}

	go inputLoop(s, event)

	return game.Loop()
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

	app.Action = func(c *cli.Context) error {
		return startGame(themes, presets)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
