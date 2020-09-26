package termi

import (
	"github.com/cicovic-andrija/2048/core"
	"github.com/gdamore/tcell"
)

const (
	digitFieldWidth    = 5
	digitFieldHeight   = 7
	blockWidth         = core.MaxBlockDigits*(digitFieldWidth-1) + 1
	blockHeight        = digitFieldHeight
	horizontalBlockGap = 2
	verticalBlockGap   = 1
)

type board struct {
	game *core.Game

	boardWidth  int
	boardHeight int

	// absolute coordinates of some top-left reference point
	refx int
	refy int

	screen tcell.Screen
	bg     tcell.Style
}

func newBoard(game *core.Game, tlx int, tly int, screen tcell.Screen) *board {
	return &board{
		game:        game,
		boardWidth:  game.Size*(blockWidth+horizontalBlockGap) + horizontalBlockGap,
		boardHeight: game.Size*(blockHeight+verticalBlockGap) + verticalBlockGap,
		refx:        tlx,
		refy:        tly,
		screen:      screen,
		bg:          tcell.StyleDefault.Background(colorGray),
	}
}

func (b *board) redraw() {
	drawRect(b.boardWidth, b.boardHeight, b.refx, b.refy, b.screen, b.bg)
	for i := 0; i < b.game.Size; i++ {
		for j := 0; j < b.game.Size; j++ {
			x := b.refx + verticalBlockGap + i*(blockHeight+verticalBlockGap)
			y := b.refy + horizontalBlockGap + j*(blockWidth+horizontalBlockGap)
			if val := b.game.Block(i, j); val != 0 {
				drawRect(blockWidth, blockHeight, x, y, b.screen, blkPropMap[val].bg)
				drawNumber(val, x, y+blkPropMap[val].inBlockPad, b.screen, blkPropMap[val].fg)
			} else {
				drawRect(blockWidth, blockHeight, x, y, b.screen, emptyCellStyle)
			}
		}
	}
}

func (b *board) push(dir core.Direction) core.Outcome {
	outcome := b.game.Push(dir)
	b.redraw()
	return outcome
}

func (b *board) undo() {
	if b.game.Undo() {
		b.redraw()
	}
}
