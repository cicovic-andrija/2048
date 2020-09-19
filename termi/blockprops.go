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
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	4: blockProps{
		inBlockPad: 6,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	8: blockProps{
		inBlockPad: 6,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	16: blockProps{
		inBlockPad: 3,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	32: blockProps{
		inBlockPad: 3,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	64: blockProps{
		inBlockPad: 3,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	128: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	256: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	512: blockProps{
		inBlockPad: 1,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	1024: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	2048: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	4096: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},

	8192: blockProps{
		inBlockPad: 0,
		bg:         tcell.StyleDefault,
		fg:         tcell.StyleDefault.Reverse(true),
	},
}
