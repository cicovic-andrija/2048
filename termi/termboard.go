package termi

import (
	"github.com/cicovic-andrija/2048/core"
	"github.com/gdamore/tcell"
)

const (
	digitFieldWidth  = 5
	digitFieldHeight = 7
	blockWidth       = core.MaxBlockDigits*(digitFieldWidth-1) + 1
	blockHeight      = digitFieldHeight
)

type TermBoard struct {
	size int // number of blocks in one dimension

	boardWidth  int
	boardHeight int

	absTopLeftx int
	absTopLefty int

	screen tcell.Screen
	bg     tcell.Style
}

func NewTermBoard(size int, screen tcell.Screen) *TermBoard {
	return nil
}
