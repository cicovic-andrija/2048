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

func NewTerminalGraphicsGame() {
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

	testDrawBox(s, 17, 7)

	// wait for esc
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			testDrawBox(s, 17, 7)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}
	}
}
