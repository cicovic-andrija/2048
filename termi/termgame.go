package termi

import (
	"fmt"
	"strings"

	"github.com/cicovic-andrija/2048/core"
	"github.com/gdamore/tcell"
)

type header struct {
	text  string
	width int
	style tcell.Style
}

type TermGame struct {
	game *core.Game

	header *header
	board  *board
	screen tcell.Screen

	// absolute coordinates of some top-left reference point
	refx int
	refy int
}

func NewTermGame(game *core.Game, tlx int, tly int) (*TermGame, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.HideCursor()
	screen.DisableMouse()
	screen.SetStyle(whiteOnBlackDefault)

	board := newBoard(game, tlx+2, tly, screen)

	header := &header{
		text:  "NEW GAME\nUse arrow keys to play / Ctrl+U to undo / Esc to quit",
		width: board.width,
		style: whiteOnGreen,
	}

	termGame := &TermGame{
		game:   game,
		header: header,
		board:  board,
		screen: screen,
		refx:   tlx,
		refy:   tly,
	}
	return termGame, nil
}

func (t *TermGame) updateHeader(outcome core.Outcome) {
	switch outcome {
	case core.Continue:
		t.header.text = fmt.Sprintf(
			"%s\nScore: %d / Undos %d",
			t.game.Player, t.game.Score(), t.game.UndosLeft(),
		)
		t.header.style = whiteOnBlue
	case core.GameOverWin:
		t.header.text = fmt.Sprintf("%s WINS! Score: %d\nPress Esc to exit", t.game.Player, t.game.Score())
		t.header.style = whiteOnGreen
	case core.GameOver:
		t.header.text = "GAME OVER! Score: 0\nPress Esc to exit"
		t.header.style = whiteOnRed
	}

	t.redrawHeader()
}

func (t *TermGame) promptConfirmQuit() {
	t.header.text = "Are you sure you want to quit?\nConfirm by pressing Esc or press any other key to continue"
	t.header.style = whiteOnRed
	t.redrawHeader()
}

func (t *TermGame) redrawHeader() {
	for i, str := range strings.Split(t.header.text, "\n") {
		drawRect(t.header.width, 1 /* height */, t.refx+i, t.refy, t.screen, t.header.style)
		drawString(str, t.refx+i, t.refy, t.screen, t.header.style)
	}
}

func (t *TermGame) waitEsc() {
	for {
		switch ev := t.screen.PollEvent().(type) {
		case *tcell.EventResize:
			t.redrawComponents()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return
			}
		}
	}
}

func (t *TermGame) redrawComponents() {
	t.redrawHeader()
	t.board.redraw()
	t.screen.Sync()
}

func (t *TermGame) Run() error {
	var (
		outcome       = core.Continue
		quitRequested = false
	)

	if t.game.Phase != core.NotStarted {
		return fmt.Errorf("terminal game has already started")
	}

	t.redrawComponents()

	// event loop
	for outcome == core.Continue {
		switch ev := t.screen.PollEvent().(type) {
		case *tcell.EventResize:
			t.redrawComponents()
			continue

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				if quitRequested {
					outcome = core.GameOver
				} else { // ask for second Esc to quit
					quitRequested = true
					t.promptConfirmQuit()
					t.screen.Show()
					continue
				}
			case tcell.KeyRight:
				outcome = t.board.push(core.Right)
			case tcell.KeyLeft:
				outcome = t.board.push(core.Left)
			case tcell.KeyUp:
				outcome = t.board.push(core.Up)
			case tcell.KeyDown:
				outcome = t.board.push(core.Down)
			case tcell.KeyCtrlU:
				t.board.undo()
			}

			t.updateHeader(outcome)
			t.screen.Show()
			quitRequested = false
		}
	}

	t.waitEsc()
	t.screen.Fini()
	return nil
}
