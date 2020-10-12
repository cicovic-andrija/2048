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

	termGame := &TermGame{
		game:   game,
		header: &header{width: board.width, style: whiteOnBlackDefault},
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
			"%s\nScore: %d | Undos %d",
			t.game.Player, t.game.Score(), t.game.UndosLeft(),
		)
	case core.GameOverWin:
		t.header.text = fmt.Sprintf("%s WINS!\nScore: %d", t.game.Player, t.game.Score())
		t.header.style = whiteOnGreen
	case core.GameOver:
		t.header.text = "GAME OVER!\nScore: 0"
		t.header.style = whiteOnRed
	}

	t.redrawHeader()
}

func (t *TermGame) redrawHeader() {
	for i, str := range strings.Split(t.header.text, "\n") {
		drawRect(t.header.width, 1, i, 0, t.screen, t.header.style)
		drawString(str, i, 0, t.screen, t.header.style)
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
	if t.game.Phase != core.NotStarted {
		return fmt.Errorf("terminal game has already started")
	}

	t.redrawComponents()

	// event loop
	outcome := core.Continue
	for outcome == core.Continue {
		t.updateHeader(outcome)
		t.screen.Show()

		switch ev := t.screen.PollEvent().(type) {

		case *tcell.EventResize:
			t.redrawComponents()
			continue

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				outcome = core.GameOver
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
		}
	}

	t.updateHeader(outcome)
	t.screen.Show()
	t.waitEsc()

	t.screen.Fini()
	return nil
}
