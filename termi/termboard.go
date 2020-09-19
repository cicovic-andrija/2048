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

type TermBoard struct {
	nblocks int // number of blocks in one dimension

	boardWidth  int
	boardHeight int

	// absolute coordinates of some top-left reference point
	refx int
	refy int

	screen tcell.Screen
	bg     tcell.Style
}

// FIXME: assume size is in limits for now
func NewTermBoard(nblocks int, tlx int, tly int, screen tcell.Screen) *TermBoard {
	return &TermBoard{
		nblocks:     nblocks,
		boardWidth:  nblocks*(blockWidth+horizontalBlockGap) + horizontalBlockGap,
		boardHeight: nblocks*(blockHeight+verticalBlockGap) + verticalBlockGap,
		refx:        tlx,
		refy:        tly,
		screen:      screen,
		bg:          tcell.StyleDefault.Background(tcell.ColorGray),
	}
}

func (b *TermBoard) draw(game *core.Game) {
	b.screen.Clear()
	drawRect(b.boardWidth, b.boardHeight, b.refx, b.refy, b.screen, b.bg)
	for i := 0; i < b.nblocks; i++ {
		for j := 0; j < b.nblocks; j++ {
			x := b.refx + verticalBlockGap + i*(blockHeight+verticalBlockGap)
			y := b.refy + horizontalBlockGap + j*(blockWidth+horizontalBlockGap)
			drawRect(blockWidth, blockHeight, x, y, b.screen, tcell.StyleDefault)
			if val := game.Block(i, j); val != 0 {
				drawNumber(val, x, y+blkPropMap[val].inBlockPad, b.screen, b.bg)
			} else {
				drawRect(blockWidth, blockHeight, x, y, b.screen, tcell.StyleDefault)
			}
		}
	}
	b.screen.Show()
}
