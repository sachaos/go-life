package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/sachaos/go-life/preset"
)

type Game struct {
	screen      tcell.Screen
	board       *Board
	themes      []Theme
	themeIndex  int
	presets     []preset.Preset
	presetIndex int
	event       chan Event
	hide        bool
	stop        bool
	ticker      *time.Ticker
}

func (g *Game) display() {
	g.screen.Clear()
	bst := tcell.StyleDefault.Background(g.themes[g.themeIndex].BackGround)
	for i, row := range g.board.State() {
		for j, cell := range row {
			st := tcell.StyleDefault.Background(g.themes[g.themeIndex].Color(cell.LiveTime()))
			if cell.State() {
				g.screen.SetCell(j*2, i, st, ' ')
				g.screen.SetCell(j*2+1, i, st, ' ')
			} else {
				g.screen.SetCell(j*2, i, bst, ' ')
				g.screen.SetCell(j*2+1, i, bst, ' ')
			}
		}
	}
}

func (g *Game) Loop() error {
	for {
		g.display()

		select {
		case ev := <-g.event:
			switch ev.Type {
			case switchState:
				if (ev.X < g.board.width) && (ev.Y < g.board.height) && (ev.X >= 0) && (ev.Y >= 0) {
					{
						g.board.Get(ev.X, ev.Y).Switch()
					}
				}
			case putPreset:
				g.board.Set(ev.X, ev.Y, g.presets[g.presetIndex].Cells)
			case resize:
				stopState := g.stop
				g.stop = true
				width, height := g.screen.Size()
				g.board.Resize(width/2, height)
				g.stop = stopState
			case step:
				g.board.Next()
				g.screen.Show()
			case done:
				return nil
			case reset:
				stopState := g.stop
				g.stop = true
				g.board.Random()
				g.stop = stopState
			case switchStop:
				g.switchStop()
			case clear:
				g.board.Init()
				g.screen.Show()
			case switchPreset:
				g.switchPreset()
			case switchHide:
				g.switchHide()
			case switchTheme:
				g.switchTheme()
			default:
				return fmt.Errorf(ev.Type)
			}
		case <-g.ticker.C:
			if !g.stop {
				g.board.Next()
			} else if !g.hide {
				g.displayMessage()
			}
			g.screen.Show()
		}
	}
}

func putString(s tcell.Screen, x, y int, str string) {
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	for i, byte := range str {
		s.SetCell(x+i, y, st, byte)
	}
}

func (g *Game) switchHide() {
	g.hide = !g.hide
}

func (g *Game) switchStop() {
	g.stop = !g.stop
}

func (g *Game) displayMessage() {
	_, height := g.screen.Size()
	putString(g.screen, 0, 0, "SPC: start, Enter: next, c: clear, r: random, h: hide this message & status")
	putString(g.screen, 0, 1, "LeftClick: switch state, RightClick: insert preset")
	putString(g.screen, 0, 2, fmt.Sprintf("p: switch preset, Current: \"%s\"", g.presets[g.presetIndex].Name))
	putString(g.screen, 0, 3, fmt.Sprintf("t: switch theme, Current: \"%s\"", g.themes[g.themeIndex].Name))
	putString(g.screen, 0, height-1, fmt.Sprintf("Time: %d", g.board.Time()))
}

func (g *Game) switchPreset() {
	if g.presetIndex < len(g.presets)-1 {
		g.presetIndex++
	} else {
		g.presetIndex = 0
	}
}

func (g *Game) switchTheme() {
	if g.themeIndex < len(g.themes)-1 {
		g.themeIndex++
	} else {
		g.themeIndex = 0
	}
}
