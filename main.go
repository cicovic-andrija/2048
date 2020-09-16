package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type D int

const (
	Left D = iota
	Right
	Up
	Down
)

type Outcome int

const (
	Continue Outcome = iota
	GameOver
	GameOverWin
)

const (
	MinSize   = 4
	MaxSize   = 5
	MinTarget = 8
	MaxTarget = 8192
)

var (
	rng *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// assumes p is [0.0, 1.0]
func ptrue(p float64) bool {
	if rng.Float64() < p {
		return true
	}
	return false
}

func randBlockval() int {
	if ptrue(0.25) {
		return 4
	}
	return 2
}

type Game struct {
	Player string
	Score  int
	Target int

	board [][]int
	size  int
}

func NewGame(player string, size int, target int) (*Game, error) {
	if player == "" {
		return nil, fmt.Errorf("player name cannot be nil")
	}

	if size < MinSize || size > MaxSize {
		return nil, fmt.Errorf("invalid size: %d, allowed range [%d, %d]",
			size, MinSize, MaxSize)
	}

	if target < MinTarget || target > MaxTarget || target&(target-1) != 0 {
		return nil,
			fmt.Errorf("invalid target: %d, allowed values: 8,16,...,8192",
				target)
	}

	board := make([][]int, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
	}

	game := &Game{
		Player: player,
		Score:  0,
		Target: target,
		board:  board,
		size:   size,
	}

	game.spawn(Left)

	return game, nil
}

func itoa(i int) string {
	if i == 0 {
		return " "
	}
	return strconv.Itoa(i)
}

func (g *Game) String() string {
	const (
		horizontalLine = "+------+------+------+------+------+"
	)
	var str strings.Builder
	str.WriteString(g.Player + ": " + strconv.Itoa(g.Score) + "\n")
	str.WriteString(horizontalLine + "\n")
	for _, row := range g.board {
		str.WriteString("| ")
		for _, blkval := range row {
			str.WriteString(fmt.Sprintf("%-4s | ", itoa(blkval)))
		}
		str.WriteString("\n" + horizontalLine + "\n")
	}
	return str.String()
}

func (g *Game) spawn(d D) {
	quota := func(free int) int {
		factor := 1
		if free > g.size/2 {
			factor = 2
		}
		return min(free, 1+min(free, g.size*factor/5))
	}

	vertfree := func(col int) int {
		c := 0
		for i := 0; i < g.size; i++ {
			if g.board[i][col] == 0 {
				c++
			}
		}
		return c
	}

	horizfree := func(row int) int {
		c := 0
		for _, blkval := range g.board[row] {
			if blkval == 0 {
				c++
			}
		}
		return c
	}

	switch d {
	case Left, Right:
		col := 0
		if d == Right {
			col = g.size - 1
		}
		for q := quota(vertfree(col)); q > 0; {
			for i := 0; i < g.size; i++ {
				if g.board[i][col] == 0 && ptrue(0.5) {
					g.board[i][col] = randBlockval()
					if q--; q == 0 {
						break
					}
				}
			}
		}
	case Up, Down:
		row := 0
		if d == Down {
			row = g.size - 1
		}
		for q := quota(horizfree(row)); q > 0; {
			for i := 0; i < g.size; i++ {
				if g.board[row][i] == 0 && ptrue(0.5) {
					g.board[row][i] = randBlockval()
					if q--; q == 0 {
						break
					}
				}
			}
		}
	}
}

func (g *Game) canMergeAnyNeighbor(i int, j int) bool {
	blkval := g.board[i][j]
	// left neighbor
	if j > 0 && g.board[i][j-1] == blkval {
		return true
	}
	// right neighbor
	if j < g.size-1 && g.board[i][j+1] == blkval {
		return true
	}
	// up neighbor
	if i > 0 && g.board[i-1][j] == blkval {
		return true
	}
	// down neighbor
	if i < g.size-1 && g.board[i+1][j] == blkval {
		return true
	}

	return false
}

func (g *Game) calcOutcome() Outcome {
	canContinue := false

	for i, row := range g.board {
		for j, blkval := range row {
			if blkval == g.Target {
				return GameOverWin
			}
			if !canContinue && (blkval == 0 || g.canMergeAnyNeighbor(i, j)) {
				canContinue = true
			}
		}
	}

	if canContinue {
		return Continue
	}
	return GameOver
}

func (g *Game) pushRight(lblock int, rblock int, fence int, row []int) {
	if rblock != fence {
		row[fence], row[rblock] = row[rblock], 0
	}
	if row[fence] == 0 || row[lblock] == 0 {
		row[fence] += row[lblock]
		row[lblock] = 0
		return
	}
	if row[fence] == row[lblock] {
		row[fence] <<= 1
		row[lblock] = 0
		g.Score += row[fence]
		return
	}
	if lblock != fence-1 {
		row[fence-1], row[lblock] = row[lblock], 0
	}
}

func (g *Game) PushRight() Outcome {
	for _, row := range g.board {
		fence := g.size - 1
		for fence > 0 {
			rblock := fence
			for rblock > 1 && row[rblock] == 0 {
				rblock--
			}
			lblock := rblock - 1
			for lblock > 0 && row[lblock] == 0 {
				lblock--
			}
			g.pushRight(lblock, rblock, fence, row)
			fence--
		}
	}

	g.spawn(Left)
	return g.calcOutcome()
}

func (g *Game) pushLeft(lblock int, rblock int, fence int, row []int) {
	if lblock != fence {
		row[fence], row[lblock] = row[lblock], 0
	}
	if row[fence] == 0 || row[rblock] == 0 {
		row[fence] += row[rblock]
		row[rblock] = 0
		return
	}
	if row[fence] == row[rblock] {
		row[fence] <<= 1
		row[rblock] = 0
		g.Score += row[fence]
		return
	}
	if rblock != fence+1 {
		row[fence+1], row[rblock] = row[rblock], 0
	}
}

func (g *Game) PushLeft() Outcome {
	for _, row := range g.board {
		fence := 0
		for fence < g.size-1 {
			lblock := fence
			for lblock < g.size-2 && row[lblock] == 0 {
				lblock++
			}
			rblock := lblock + 1
			for rblock < g.size-1 && row[rblock] == 0 {
				rblock++
			}
			g.pushLeft(lblock, rblock, fence, row)
			fence++
		}
	}

	g.spawn(Right)
	return g.calcOutcome()
}

func (g *Game) pushUp(ublock int, dblock int, fence int, col int) {
	if ublock != fence {
		g.board[fence][col], g.board[ublock][col] = g.board[ublock][col], 0
	}
	if g.board[fence][col] == 0 || g.board[dblock][col] == 0 {
		g.board[fence][col] += g.board[dblock][col]
		g.board[dblock][col] = 0
		return
	}
	if g.board[fence][col] == g.board[dblock][col] {
		g.board[fence][col] <<= 1
		g.board[dblock][col] = 0
		g.Score += g.board[fence][col]
		return
	}
	if dblock != fence+1 {
		g.board[fence+1][col], g.board[dblock][col] = g.board[dblock][col], 0
	}
}

func (g *Game) PushUp() Outcome {
	for col := 0; col < g.size; col++ {
		fence := 0
		for fence < g.size-1 {
			ublock := fence
			for ublock < g.size-2 && g.board[ublock][col] == 0 {
				ublock++
			}
			dblock := ublock + 1
			for dblock < g.size-1 && g.board[dblock][col] == 0 {
				dblock++
			}
			g.pushUp(ublock, dblock, fence, col)
			fence++
		}
	}

	g.spawn(Down)
	return g.calcOutcome()
}

func (g *Game) pushDown(ublock int, dblock int, fence int, col int) {
	if dblock != fence {
		g.board[fence][col], g.board[dblock][col] = g.board[dblock][col], 0
	}
	if g.board[fence][col] == 0 || g.board[ublock][col] == 0 {
		g.board[fence][col] += g.board[ublock][col]
		g.board[ublock][col] = 0
		return
	}
	if g.board[fence][col] == g.board[ublock][col] {
		g.board[fence][col] <<= 1
		g.board[ublock][col] = 0
		g.Score += g.board[fence][col]
		return
	}
	if ublock != fence-1 {
		g.board[fence-1][col], g.board[ublock][col] = g.board[ublock][col], 0
	}
}

func (g *Game) PushDown() Outcome {
	for col := 0; col < g.size; col++ {
		fence := g.size - 1
		for fence > 0 {
			dblock := fence
			for dblock > 1 && g.board[dblock][col] == 0 {
				dblock--
			}
			ublock := dblock - 1
			for ublock > 0 && g.board[ublock][col] == 0 {
				ublock--
			}
			g.pushDown(ublock, dblock, fence, col)
			fence--
		}
	}

	g.spawn(Up)
	return g.calcOutcome()
}

func main() {
	NewTextGame("Andrija", 5, 2048)
}

func NewTextGame(player string, size int, target int) {
	game, err := NewGame(player, size, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in game initialization: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", game.String())
	reader := bufio.NewReader(os.Stdin)
	outcome := Continue
	for outcome == Continue {
		// read a command
		fmt.Printf("Command: ")
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		// execute the command
		switch char {
		case 'd', 'D':
			outcome = game.PushRight()
			fmt.Printf("%s", game.String())
		case 'a', 'A':
			outcome = game.PushLeft()
			fmt.Printf("%s", game.String())
		case 'w', 'W':
			outcome = game.PushUp()
			fmt.Printf("%s", game.String())
		case 's', 'S':
			outcome = game.PushDown()
			fmt.Printf("%s", game.String())
		case 'e', 'E', 'q', 'Q':
			outcome = GameOver
		}
	}

	if outcome == GameOverWin {
		fmt.Printf("%s WINS! Score: %d\n", game.Player, game.Score)
	} else {
		fmt.Printf("GAME OVER!\n")
	}
}
