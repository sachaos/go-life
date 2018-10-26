package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/sachaos/go-life/preset"
	"log"
)

func main() {
	log.Print("start")
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

	rand.Seed(time.Now().Unix())

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
	defer s.Fini()

	s.EnableMouse()

	// init board
	width, height := s.Size()
	b := NewBoard(height, width/2)
	b.Init()
	b.Random()

	// init ticker
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	game := Game{
		screen:  s,
		board:   b,
		themes:  themes,
		presets: presets,
		ticker:  ticker,
		event:   make(chan Event),
	}

	go func() {
		for {
			log.Print("for loop")
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					x, y := ev.Position()
					game.event <- Event{Type: swtichState, X: x / 2, Y: y}
				case tcell.Button3:
					x, y := ev.Position()
					game.event <- Event{Type: putPreset, X: x / 2, Y: y}
				default:
					continue
				}
			case *tcell.EventResize:
				game.event <- Event{Type: resize}
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEnter {
					log.Print("event: enter")
					game.event <- Event{Type: step}
				} else if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					log.Print("event: exit")
					game.event <- Event{Type: done}
				} else if ev.Rune() == ' ' {
					game.event <- Event{Type: switchStop}
				} else if ev.Rune() == 'c' {
					game.event <- Event{Type: clear}
				} else if ev.Rune() == 'p' {
					game.event <- Event{Type: switchPreset}
				} else if ev.Rune() == 'r' {
					game.event <- Event{Type: reset}
				} else if ev.Rune() == 't' {
					game.event <- Event{Type: switchTheme}
				} else if ev.Rune() == 'h' {
					game.event <- Event{Type: switchHide}
				}
			default:
				continue
			}
		}
	}()

	if err := game.Loop(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
