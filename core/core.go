package core

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Terminology:
//   cell - a position inside a matrix
//   block - a value of the cell
//     - 0 is a special value meaning "no block"

const (
	MinSize = 4
	MaxSize = 5

	MaxBlock       = 8192
	MaxBlockDigits = 4

	MinTarget = 2048
	MaxTarget = MaxBlock

	blockFourProbability float64 = 0.15
)

type Direction int

const (
	Right = iota
	Left
	Up
	Down
)

type Outcome int

const (
	Continue Outcome = iota
	GameOver
	GameOverWin
)

type stableState struct {
	board    [][]int // matrix of cells
	score    int     // player's score
	blockCnt int     // number of blocks on the board
}

// assumes size is in limits
func newInitialState(size int) *stableState {
	// all cells are initially empty (aka all blocks are 0)
	board := make([][]int, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size)
	}

	return &stableState{
		board:    board,
		score:    0,
		blockCnt: 0,
	}
}

// assumes board size of both states are equal, since this
// function is only used inside the package
func (s *stableState) deepCopyFrom(other *stableState) {
	for i, row := range other.board {
		for j, block := range row {
			s.board[i][j] = block
		}
	}
	s.score, s.blockCnt = other.score, other.blockCnt
}

type Game struct {
	Player string // player's name
	Target int    // end-game block
	Size   int    // size of the board

	stableState                  // embedded current state
	prevStableState *stableState // previous stable state
	scratchState    *stableState // scratch state

	canUndo       bool
	UndosLeft     int
	anyBlockMoved bool
	rng           *rand.Rand // random number generator
}

func NewGame(player string, size int, target int, undos int) (*Game, error) {
	// param validation
	//
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

	if undos < 0 {
		undos = 0
	}

	// create a new game
	game := &Game{
		Player:          player,
		Target:          target,
		Size:            size,
		stableState:     *newInitialState(size),
		prevStableState: newInitialState(size),
		scratchState:    newInitialState(size),
		canUndo:         false,
		UndosLeft:       undos,
		anyBlockMoved:   false,
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
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

func (g *Game) copyToPrevState() {
	if g.UndosLeft == 0 {
		return
	}
	g.prevStableState.deepCopyFrom(&g.stableState)
}

func (g *Game) rollbackPrevState() {
	if g.UndosLeft == 0 {
		return
	}
	g.prevStableState.deepCopyFrom(g.scratchState)
}

func (g *Game) commitPrevState() {
	if g.UndosLeft == 0 {
		return
	}
	g.scratchState.deepCopyFrom(g.prevStableState)
}

func (g *Game) _pushRight(lcell int, rcell int, fence int, row []int) {
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
		g.score += row[fence]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if lcell != fence-1 {
		row[fence-1], row[lcell] = row[lcell], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) _pushLeft(lcell int, rcell int, fence int, row []int) {
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
		g.score += row[fence]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if rcell != fence+1 {
		row[fence+1], row[rcell] = row[rcell], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) _pushUp(ucell int, dcell int, fence int, col int) {
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
		g.score += g.board[fence][col]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if dcell != fence+1 {
		g.board[fence+1][col], g.board[dcell][col] = g.board[dcell][col], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) _pushDown(ucell int, dcell int, fence int, col int) {
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
		g.score += g.board[fence][col]
		g.blockCnt--
		g.anyBlockMoved = true
		return
	}

	if ucell != fence-1 {
		g.board[fence-1][col], g.board[ucell][col] = g.board[ucell][col], 0
		g.anyBlockMoved = true
	}
}

func (g *Game) pushRight() {
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
			g._pushRight(lcell, rcell, fence, row)
			fence--
		}
	}
}

func (g *Game) pushLeft() {
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
			g._pushLeft(lcell, rcell, fence, row)
			fence++
		}
	}
}

func (g *Game) pushUp() {
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
			g._pushUp(ucell, dcell, fence, col)
			fence++
		}
	}
}

func (g *Game) pushDown() {
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
			g._pushDown(ucell, dcell, fence, col)
			fence--
		}
	}
}

func (g *Game) Push(dir Direction) Outcome {
	g.copyToPrevState()

	switch dir {
	case Right:
		g.pushRight()
	case Left:
		g.pushLeft()
	case Up:
		g.pushUp()
	case Down:
		g.pushDown()
	}

	if !g.anyBlockMoved {
		g.rollbackPrevState()
		return Continue
	}

	g.anyBlockMoved = false
	g.canUndo = true
	g.spawn()
	g.commitPrevState()
	return g.calcOutcome()
}

func (g *Game) Undo() bool {
	if g.UndosLeft == 0 || !g.canUndo {
		return false
	}
	g.stableState.deepCopyFrom(g.prevStableState)
	g.UndosLeft--
	g.canUndo = false // two consecutive undos are impossible
	return true
}

func (g *Game) Score() int {
	return g.score
}

func (g *Game) Block(i int, j int) int {
	if i < 0 || i >= g.Size || j < 0 || j >= g.Size {
		return -1
	}
	return g.board[i][j]
}
