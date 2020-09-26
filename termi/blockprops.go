package termi

import "github.com/gdamore/tcell"

type blockProps struct {
	inBlockPad int
	bg         tcell.Style
	fg         tcell.Style
}

var blkPropMap = map[int]blockProps{
	2: blockProps{
		inBlockPad: 6,
		bg:         tcell.StyleDefault.Background(colorAlmostWhite),
		fg:         tcell.StyleDefault.Background(colorDarkGray),
	},
	4: blockProps{
		inBlockPad: 6,
		bg:         tcell.StyleDefault.Background(tcell.Color229),
		fg:         tcell.StyleDefault.Background(colorDarkGray),
	},

	8: blockProps{
		inBlockPad: 6,
		bg:         tcell.StyleDefault.Background(tcell.Color222),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	16: blockProps{
		inBlockPad: 3,
		bg:         tcell.StyleDefault.Background(tcell.Color215),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	32: blockProps{
		inBlockPad: 3,
		bg:         tcell.StyleDefault.Background(tcell.Color209),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	64: blockProps{

		inBlockPad: 3,
		bg:         tcell.StyleDefault.Background(tcell.Color196),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	128: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault.Background(tcell.Color229),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	256: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault.Background(tcell.Color229),
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	512: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	1024: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	2048: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	4096: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},

	8192: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Background(tcell.ColorWhite),
	},
}
