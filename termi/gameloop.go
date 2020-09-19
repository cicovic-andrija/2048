package termi

import (
	"fmt"
	"os"

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
	board.draw()

	// wait for esc
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			board.draw()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}
	}
}
