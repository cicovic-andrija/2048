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

func testDraw2048(s tcell.Screen) {
	bg := tcell.StyleDefault
	fg := tcell.StyleDefault.Reverse(true)
	s.Clear()
	drawNumber(2048, 0, 0, s, fg, bg)
	s.Show()
}

func testDrawBox(s tcell.Screen, w int, h int) {
	bg := tcell.StyleDefault
	fg := tcell.StyleDefault.Reverse(true)
	s.Clear()
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			s.SetContent(y, x, ' ', nil, fg)
		}
	}
	drawNumber(2048, 0, 0, s, bg, fg)
	s.Show()
}
