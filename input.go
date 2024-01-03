package main

import (
	"github.com/gdamore/tcell"
	"log"
)

func inputLoop(s tcell.Screen, event chan<- Event) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1:
				x, y := ev.Position()
				event <- Event{Type: switchState, X: x / 2, Y: y}
			case tcell.Button3:
				x, y := ev.Position()
				event <- Event{Type: putPreset, X: x / 2, Y: y}
			default:
				continue
			}
		case *tcell.EventResize:
			event <- Event{Type: resize}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				log.Print("event: enter")
				event <- Event{Type: step}
			} else if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				log.Print("event: exit")
				event <- Event{Type: done}
			} else if ev.Rune() == ' ' {
				event <- Event{Type: switchStop}
			} else if ev.Rune() == 'c' {
				event <- Event{Type: clear}
			} else if ev.Rune() == 'p' {
				event <- Event{Type: switchPreset}
			} else if ev.Rune() == 'r' {
				event <- Event{Type: reset}
			} else if ev.Rune() == 't' {
				event <- Event{Type: switchTheme}
			} else if ev.Rune() == 'h' {
				event <- Event{Type: switchHide}
			}
		default:
			continue
		}
	}
}
