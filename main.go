package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func putString(s tcell.Screen, x, y int, str string) {
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	for i, byte := range str {
		s.SetCell(x+i, y, st, byte)
	}
}

func main() {
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

	done := make(chan struct{})
	stopSwtich := make(chan struct{})
	reset := make(chan struct{})
	step := make(chan struct{})
	clear := make(chan struct{})
	resize := make(chan struct{})
	hideMessage := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					x, y := ev.Position()
					b.Get(x/2, y).Switch()
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
				} else if ev.Rune() == 'r' {
					reset <- struct{}{}
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
		for i, row := range b.State() {
			for j, cell := range row {
				st := tcell.StyleDefault.Background(ThemeBlackAndWhite.Color(cell.LiveTime()))
				if cell.State() {
					s.SetCell(j*2, i, st, ' ')
					s.SetCell(j*2+1, i, st, ' ')
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
		case <-hideMessage:
			hide = !hide
		case <-ticker.C:
			if !stop {
				b.Next()
			} else if hide == false {
				putString(s, 0, 0, "SPC: start, Enter: next, c: clear, r: random, h: hide this message")
				putString(s, 0, 1, "LeftClick: switch state, RightClick: select preset")
			}
			s.Show()
		}
	}
}
