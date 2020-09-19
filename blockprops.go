package main

import "github.com/gdamore/tcell"

type blockProps struct {
	ndigits int
	bg      tcell.Style
	fg      tcell.Style
}

var blkPropMap = [int]blockProps{
	2: blockProps{
		ndigits: 1,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	4: blockProps{
		ndigits: 1,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	8: blockProps{
		ndigits: 1,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	16: blockProps{
		ndigits: 2,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	32: blockProps{
		ndigits: 2,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	64: blockProps{
		ndigits: 2,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	128: blockProps{
		ndigits: 3,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	256: blockProps{
		ndigits: 3,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	512: blockProps{
		ndigits: 3,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	1024: blockProps{
		ndigits: 4,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	2048: blockProps{
		ndigits: 4,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	4096: blockProps{
		ndigits: 4,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},

	8192: blockProps{
		ndigits: 4,
		bg:      tcell.StyleDefault,
		fg:      tcell.StyleDefault.Reverse(true),
	},
}
