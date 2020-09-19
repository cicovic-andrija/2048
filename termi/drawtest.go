package termi

import "github.com/gdamore/tcell"

func testDrawNumbers(s tcell.Screen) {
	s.Clear()
	sx, sy := 1, 1
	for n := 0; n < 10; n++ {
		for i := 0; i < 16; i++ {
			st := tcell.StyleDefault
			if bitmap[n]&(1<<i) != 0 {
				st = tcell.StyleDefault.Reverse(true)
			}
			s.SetContent(sy+i%3, sx+i/3, ' ', nil, st)
		}
		sy += 4
	}
	s.Show()
}

func testDrawBoxedNumber(s tcell.Screen) {
	fg := tcell.StyleDefault
	bg := tcell.StyleDefault.Reverse(true)
	s.Clear()
	for x := 0; x < blockHeight; x++ {
		for y := 0; y < blockWidth; y++ {
			s.SetContent(y, x, ' ', nil, bg)
		}
	}
	drawNumber(2048, 0, 0, s, fg)
	s.Show()
}
