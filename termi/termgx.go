package termi

import (
	"github.com/gdamore/tcell"
)

const (
	colorAlmostWhite = tcell.Color254
	colorLightGray   = tcell.Color250
	colorGray        = tcell.Color246
	colorDarkGray    = tcell.Color242
)

var (
	bitmap = map[int]uint16{
		0: 0x7b6f,
		1: 0x4926,
		2: 0x73e7,
		3: 0x79e7,
		4: 0x49ed,
		5: 0x79cf,
		6: 0x7bcf,
		7: 0x4927,
		8: 0x7bef,
		9: 0x79ef,
	}

	emptyCellStyle = tcell.StyleDefault.
			Background(colorLightGray)
	whiteOnBlackDefault = tcell.StyleDefault.
				Background(tcell.ColorBlack).
				Foreground(tcell.ColorWhite)
)

// assumes n > 0
// ignore return value, it is used internally
func drawNumber(n int, tlx int, tly int, s tcell.Screen, st tcell.Style) int {
	if n/10 > 0 {
		tly = drawNumber(n/10, tlx, tly, s, st) + digitFieldWidth - 1
	}
	digit := n % 10
	for i := 0; i < 16; i++ {
		if bitmap[digit]&(1<<i) != 0 {
			s.SetContent(tly+i%3+1, tlx+i/3+1, ' ', nil, st)
		}
	}
	return tly
}

func drawRect(w int, h int, tlx int, tly int, s tcell.Screen, st tcell.Style) {
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			s.SetContent(tly+y, tlx+x, ' ', nil, st)
		}
	}
}

func drawString(str string, tlx int, tly int, s tcell.Screen, st tcell.Style) {
	for _, c := range str {
		s.SetContent(tly, tlx, c, nil, st)
		tly++
	}
}
