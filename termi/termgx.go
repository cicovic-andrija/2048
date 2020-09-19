package termi

import (
	"github.com/gdamore/tcell"
)

var bitmap = map[int]uint16{
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

// assumes n > 0
// ignore return value, it is used internally
func drawNumber(n int, tlx int, tly int, s tcell.Screen, fg tcell.Style, bg tcell.Style) int {
	if n/10 > 0 {
		tly = drawNumber(n/10, tlx, tly, s, fg, bg) + digitFieldWidth - 1
	}
	digit := n % 10
	for i := 0; i < 16; i++ {
		if bitmap[digit]&(1<<i) != 0 {
			s.SetContent(tly+i%3+1, tlx+i/3+1, ' ', nil, fg)
		} else {
			s.SetContent(tly+i%3+1, tlx+i/3+1, ' ', nil, bg)
		}
	}
	return tly
}