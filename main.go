package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/sachaos/go-life/preset"
)

func putString(s tcell.Screen, x, y int, str string) {
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	for i, byte := range str {
		s.SetCell(x+i, y, st, byte)
	}
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
		ThemeWhiteAndBlack,
	}

	rand.Seed(time.Now().Unix())
	stop := false
	hide := false
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
		Background(tcell.ColorDarkCyan).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()

	width, height := s.Size()

	// init board
	b := NewBoard(height, width/2)
	b.Init()
	b.Random()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	presetIndex := 0
	themeIndex := 0

	done := make(chan struct{})
	stopSwtich := make(chan struct{})
	reset := make(chan struct{})
	step := make(chan struct{})
	clear := make(chan struct{})
	resize := make(chan struct{})
	hideMessage := make(chan struct{})
	switchPreset := make(chan struct{})
	switchTheme := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					x, y := ev.Position()
					b.Get(x/2, y).Switch()
				case tcell.Button3:
					x, y := ev.Position()
					b.Set(x/2, y, presets[presetIndex].Cells)
				default:
					continue
				}
			case *tcell.EventResize:
				resize <- struct{}{}
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEnter {
					step <- struct{}{}
				} else if ev.Key() == tcell.KeyEsc || ev.Rune() == 'q' {
					done <- struct{}{}
				} else if ev.Rune() == ' ' {
					stopSwtich <- struct{}{}
				} else if ev.Rune() == 'c' {
					clear <- struct{}{}
				} else if ev.Rune() == 'p' {
					switchPreset <- struct{}{}
				} else if ev.Rune() == 'r' {
					reset <- struct{}{}
				} else if ev.Rune() == 't' {
					switchTheme <- struct{}{}
				} else if ev.Rune() == 'h' {
					hideMessage <- struct{}{}
				}
			default:
				continue
			}
		}
	}()

	for {
		s.Clear()
		bst := tcell.StyleDefault.Background(themes[themeIndex].BackGround)
		for i, row := range b.State() {
			for j, cell := range row {
				st := tcell.StyleDefault.Background(themes[themeIndex].Color(cell.LiveTime()))
				if cell.State() {
					s.SetCell(j*2, i, st, ' ')
					s.SetCell(j*2+1, i, st, ' ')
				} else {
					s.SetCell(j*2, i, bst, ' ')
					s.SetCell(j*2+1, i, bst, ' ')
				}
			}
		}

		select {
		case <-reset:
			stopState := stop
			stop = true
			b.Random()
			stop = stopState
		case <-stopSwtich:
			stop = !stop
		case <-resize:
			stopState := stop
			stop = true
			width, height := s.Size()
			b.Resize(width/2, height)
			stop = stopState
		case <-step:
			b.Next()
			s.Show()
		case <-done:
			s.Fini()
			os.Exit(0)
		case <-clear:
			b.Init()
			s.Show()
		case <-switchPreset:
			if presetIndex < len(presets)-1 {
				presetIndex++
			} else {
				presetIndex = 0
			}
		case <-hideMessage:
			hide = !hide
		case <-switchTheme:
			if themeIndex < len(themes)-1 {
				themeIndex++
			} else {
				themeIndex = 0
			}
		case <-ticker.C:
			if !stop {
				b.Next()
			} else if hide == false {
				_, height = s.Size()
				putString(s, 0, 0, "SPC: start, Enter: next, c: clear, r: random, h: hide this message & status")
				putString(s, 0, 1, "LeftClick: switch state, RightClick: insert preset")
				putString(s, 0, 2, fmt.Sprintf("p: switch preset, Current: \"%s\"", presets[presetIndex].Name))
				putString(s, 0, 3, fmt.Sprintf("t: switch theme, Current: \"%s\"", themes[themeIndex].Name))
				putString(s, 0, height-1, fmt.Sprintf("Time: %d", b.Time()))
			}
			s.Show()
		}
	}
}
