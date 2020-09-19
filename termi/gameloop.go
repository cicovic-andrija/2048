package termi

import (
	"fmt"
	"os"

	"github.com/cicovic-andrija/2048/core"
	"github.com/gdamore/tcell"
)

func errorout(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func NewTerminalGraphicsGame(size int) {
	s, e := tcell.NewScreen()
	if e != nil {
		errorout(e)
	}
	if e = s.Init(); e != nil {
		errorout(e)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	s.SetStyle(defStyle)
	s.HideCursor()
	s.DisableMouse()

	board := NewTermBoard(size, 0, 0, s)
	game, _ := core.NewGame("Andrija", 4, 2048, 3)

	board.draw(game)
	// wait for esc
	for outcome := core.Continue; outcome == core.Continue; {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			board.draw(game)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			case tcell.KeyRight:
				outcome = game.Push(core.Right)
				board.draw(game)
			case tcell.KeyLeft:
				outcome = game.Push(core.Left)
				board.draw(game)
			case tcell.KeyUp:
				outcome = game.Push(core.Up)
				board.draw(game)
			case tcell.KeyDown:
				outcome = game.Push(core.Down)
				board.draw(game)
			}
		}
	}
}
