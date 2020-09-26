package termi

import (
	"fmt"

	"github.com/cicovic-andrija/2048/core"
	"github.com/gdamore/tcell"
)

type State int

type TermGame struct {
	game *core.Game

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

	termGame := &TermGame{
		game:   game,
		board:  newBoard(game, tlx+2, tly, screen),
		screen: screen,
		refx:   tlx,
		refy:   tly,
	}
	return termGame, nil
}

func (t *TermGame) redrawHeader() {
	drawString(
		fmt.Sprintf("%s", t.game.Player),
		0, 0, t.screen, whiteOnBlackDefault,
	)
	drawString(
		fmt.Sprintf("Score: %d | Undos: %d", t.game.Score(), t.game.UndosLeft),
		1, 0, t.screen, whiteOnBlackDefault,
	)
}

func (t *TermGame) redrawComponents() {
	t.redrawHeader()
	t.board.redraw()
	t.screen.Sync()
}

func (t *TermGame) Run() {
	if t.game.Phase != core.NotStarted {
		return
	}

	t.redrawComponents()

	// event loop
	for outcome := core.Continue; outcome == core.Continue; {
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

		switch outcome {
		case core.Continue:
			t.redrawHeader()
			t.screen.Show()
		case core.GameOver:
			t.screen.Fini()
		case core.GameOverWin:
			t.screen.Fini()
		}
	}
}
