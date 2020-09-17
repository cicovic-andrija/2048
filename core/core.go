package core

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	MinSize   = 4
	MaxSize   = 6
	MinTarget = 2048
	MaxTarget = 8192

	blockFourProbability float64 = 0.15
)

type Outcome int

const (
	Continue Outcome = iota
	GameOver
	GameOverWin
)

type Game struct {
	Player string // player's name
	Score  int    // player's score
	Target int    // end-game block
	Size   int    // size of the board

	board         [][]int    // matrix of Size x Size cells
	blockCnt      int        // number of free cells
	anyBlockMoved bool       // any block moved in the current turn?
	rng           *rand.Rand // random number generator
}

func NewGame(player string, size int, target int) (*Game, error) {
	// param validation

	if player == "" {
		return nil, errors.New("player name cannot be empty")
	}

	if size < MinSize || size > MaxSize {
		errmsg := fmt.Sprintf(
			"invalid size: %d, allowed range [%d, %d]",
			size,
			MinSize,
			MaxSize,
		)
		return nil, errors.New(errmsg)
	}

	if target < MinTarget || target > MaxTarget || target&(target-1) != 0 {
		errmsg := fmt.Sprintf(
			"invalid target: %d, allowed values: 2048, 4096, 8192",
			target,
		)
		return nil, errors.New(errmsg)
	}

	// create a new game

	board := make([][]int, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
	}

	game := &Game{
		Player:        player,
		Score:         0,
		Target:        target,
		Size:          size,
		board:         board,
		blockCnt:      0,
		anyBlockMoved: false,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// spawn two blocks at the start of the game
	game.spawn()
	game.spawn()

	return game, nil
}

func (g *Game) randBlock() int {
	if ptrue(blockFourProbability, g.rng) {
		return 4
	}
	return 2
}

func (g *Game) spawn() {
	if g.blockCnt == g.Size*g.Size {
		return
	}

	for {
		n := g.rng.Intn(g.Size * g.Size)
		i, j := n/g.Size, n%g.Size
		if g.board[i][j] == 0 {
			g.board[i][j] = g.randBlock()
			break
		}
	}

	g.blockCnt++
}

func (g *Game) canMergeAnyNeighbor(i int, j int, blkval int) bool {
	// left neighbor
	if j > 0 && g.board[i][j-1] == blkval {
		return true
	}
	// right neighbor
	if j < g.Size-1 && g.board[i][j+1] == blkval {
		return true
	}
	// up neighbor
	if i > 0 && g.board[i-1][j] == blkval {
		return true
	}
	// down neighbor
	if i < g.Size-1 && g.board[i+1][j] == blkval {
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
			if !canContinue &&
				(blkval == 0 || g.canMergeAnyNeighbor(i, j, blkval)) {

				canContinue = true
			}
		}
	}

	if canContinue {
		return Continue
	}
	return GameOver
}

func (g *Game) pushRight(lcell int, rcell int, fence int, row []int) {
	if rcell != fence && row[rcell] != 0 {
		row[fence], row[rcell] = row[rcell], 0
		g.anyBlockMoved = true
	}

	if row[fence] == 0 || row[lcell] == 0 {
		if row[lcell] != 0 {
			g.anyBlockMoved = true
		}
		row[fence] += row[lcell]
		row[lcell] = 0
		return
	}

	if row[fence] == row[lcell] {
		row[fence] <<= 1
		row[lcell] = 0
		g.Score += row[fence]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if lcell != fence-1 {
		row[fence-1], row[lcell] = row[lcell], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) pushLeft(lcell int, rcell int, fence int, row []int) {
	if lcell != fence && row[lcell] != 0 {
		row[fence], row[lcell] = row[lcell], 0
		g.anyBlockMoved = true
	}

	if row[fence] == 0 || row[rcell] == 0 {
		if row[rcell] != 0 {
			g.anyBlockMoved = true
		}
		row[fence] += row[rcell]
		row[rcell] = 0
		return
	}

	if row[fence] == row[rcell] {
		row[fence] <<= 1
		row[rcell] = 0
		g.Score += row[fence]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if rcell != fence+1 {
		row[fence+1], row[rcell] = row[rcell], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) pushUp(ucell int, dcell int, fence int, col int) {
	if ucell != fence && g.board[ucell][col] != 0 {
		g.board[fence][col], g.board[ucell][col] = g.board[ucell][col], 0
		g.anyBlockMoved = true
	}

	if g.board[fence][col] == 0 || g.board[dcell][col] == 0 {
		if g.board[dcell][col] != 0 {
			g.anyBlockMoved = true
		}
		g.board[fence][col] += g.board[dcell][col]
		g.board[dcell][col] = 0
		return
	}

	if g.board[fence][col] == g.board[dcell][col] {
		g.board[fence][col] <<= 1
		g.board[dcell][col] = 0
		g.Score += g.board[fence][col]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if dcell != fence+1 {
		g.board[fence+1][col], g.board[dcell][col] = g.board[dcell][col], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) pushDown(ucell int, dcell int, fence int, col int) {
	if dcell != fence && g.board[dcell][col] != 0 {
		g.board[fence][col], g.board[dcell][col] = g.board[dcell][col], 0
		g.anyBlockMoved = true
	}

	if g.board[fence][col] == 0 || g.board[ucell][col] == 0 {
		if g.board[ucell][col] != 0 {
			g.anyBlockMoved = true
		}
		g.board[fence][col] += g.board[ucell][col]
		g.board[ucell][col] = 0
		return
	}

	if g.board[fence][col] == g.board[ucell][col] {
		g.board[fence][col] <<= 1
		g.board[ucell][col] = 0
		g.Score += g.board[fence][col]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if ucell != fence-1 {
		g.board[fence-1][col], g.board[ucell][col] = g.board[ucell][col], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) PushRight() Outcome {
	for _, row := range g.board {
		fence := g.Size - 1
		for fence > 0 {
			rcell := fence
			for rcell > 1 && row[rcell] == 0 {
				rcell--
			}
			lcell := rcell - 1
			for lcell > 0 && row[lcell] == 0 {
				lcell--
			}
			g.pushRight(lcell, rcell, fence, row)
			fence--
		}
	}

	if g.anyBlockMoved {
		g.spawn()
	}
	g.anyBlockMoved = false

	return g.calcOutcome()
}

func (g *Game) PushLeft() Outcome {
	for _, row := range g.board {
		fence := 0
		for fence < g.Size-1 {
			lcell := fence
			for lcell < g.Size-2 && row[lcell] == 0 {
				lcell++
			}
			rcell := lcell + 1
			for rcell < g.Size-1 && row[rcell] == 0 {
				rcell++
			}
			g.pushLeft(lcell, rcell, fence, row)
			fence++
		}
	}

	if g.anyBlockMoved {
		g.spawn()
	}
	g.anyBlockMoved = false

	return g.calcOutcome()
}

func (g *Game) PushUp() Outcome {
	for col := 0; col < g.Size; col++ {
		fence := 0
		for fence < g.Size-1 {
			ucell := fence
			for ucell < g.Size-2 && g.board[ucell][col] == 0 {
				ucell++
			}
			dcell := ucell + 1
			for dcell < g.Size-1 && g.board[dcell][col] == 0 {
				dcell++
			}
			g.pushUp(ucell, dcell, fence, col)
			fence++
		}
	}

	if g.anyBlockMoved {
		g.spawn()
	}
	g.anyBlockMoved = false

	return g.calcOutcome()
}

func (g *Game) PushDown() Outcome {
	for col := 0; col < g.Size; col++ {
		fence := g.Size - 1
		for fence > 0 {
			dcell := fence
			for dcell > 1 && g.board[dcell][col] == 0 {
				dcell--
			}
			ucell := dcell - 1
			for ucell > 0 && g.board[ucell][col] == 0 {
				ucell--
			}
			g.pushDown(ucell, dcell, fence, col)
			fence--
		}
	}

	if g.anyBlockMoved {
		g.spawn()
	}
	g.anyBlockMoved = false

	return g.calcOutcome()
}

func (g *Game) Cell(i int, j int) int {
	if i < 0 || i >= g.Size || j < 0 || j >= g.Size {
		return -1
	}
	return g.board[i][j]
}
