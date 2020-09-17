package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Outcome int

const (
	Continue Outcome = iota
	GameOver
	GameOverWin
)

const (
	MinSize   = 4
	MaxSize   = 6
	MinTarget = 2048
	MaxTarget = 8192

	blockFourProbability float64 = 0.15
)

var (
	rng               *rand.Rand
	textiHorizLine    string
	textiScoreLineFmt string
)

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
	if ptrue(blockFourProbability) {
		return 4
	}
	return 2
}

// assumes board size is in limits
func buildTextiParts(playerName string, boardSize int) {
	textiHorizLine = "\n+" + strings.Repeat("------+", boardSize) + "\n"
	textiScoreLineFmt = playerName + "'s score: %d"
}

type Game struct {
	Player string // player's name
	Score  int    // player's score
	Target int    // end-game block
	Size   int    // size of the board

	board    [][]int // matrix of Size x Size cells
	blockCnt int     // number of free cells
}

func NewGame(player string, size int, target int) (*Game, error) {
	if player == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}

	if size < MinSize || size > MaxSize {
		return nil, fmt.Errorf("invalid size: %d, allowed range [%d, %d]",
			size, MinSize, MaxSize)
	}

	if target < MinTarget || target > MaxTarget || target&(target-1) != 0 {
		return nil,
			fmt.Errorf("invalid target: %d, allowed values: 2048,4096,8192",
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

		Size:     size,
		board:    board,
		blockCnt: 0,
	}

	game.spawn()
	game.spawn()

	return game, nil
}

func (g *Game) String() string {
	tostring := func(i int) string {
		if i == 0 {
			return " "
		}
		return strconv.Itoa(i)
	}

	var str strings.Builder

	str.WriteString(fmt.Sprintf(textiScoreLineFmt, g.Score))
	str.WriteString(textiHorizLine)
	for _, row := range g.board {
		for _, blkval := range row {
			str.WriteString(fmt.Sprintf("| %-4s ", tostring(blkval)))
		}
		str.WriteString("|" + textiHorizLine)
	}

	return str.String()
}

func (g *Game) spawn() {
	if g.blockCnt == g.Size*g.Size {
		return
	}

	for {
		n := rng.Intn(g.Size * g.Size)
		i, j := n/g.Size, n%g.Size
		if g.board[i][j] == 0 {
			g.board[i][j] = randBlockval()
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
		g.blockCnt--
		return
	}
	if lblock != fence-1 {
		row[fence-1], row[lblock] = row[lblock], 0
	}
}

func (g *Game) PushRight() Outcome {
	for _, row := range g.board {
		fence := g.Size - 1
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

	g.spawn()
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
		g.blockCnt--
		return
	}
	if rblock != fence+1 {
		row[fence+1], row[rblock] = row[rblock], 0
	}
}

func (g *Game) PushLeft() Outcome {
	for _, row := range g.board {
		fence := 0
		for fence < g.Size-1 {
			lblock := fence
			for lblock < g.Size-2 && row[lblock] == 0 {
				lblock++
			}
			rblock := lblock + 1
			for rblock < g.Size-1 && row[rblock] == 0 {
				rblock++
			}
			g.pushLeft(lblock, rblock, fence, row)
			fence++
		}
	}

	g.spawn()
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
		g.blockCnt--
		return
	}
	if dblock != fence+1 {
		g.board[fence+1][col], g.board[dblock][col] = g.board[dblock][col], 0
	}
}

func (g *Game) PushUp() Outcome {
	for col := 0; col < g.Size; col++ {
		fence := 0
		for fence < g.Size-1 {
			ublock := fence
			for ublock < g.Size-2 && g.board[ublock][col] == 0 {
				ublock++
			}
			dblock := ublock + 1
			for dblock < g.Size-1 && g.board[dblock][col] == 0 {
				dblock++
			}
			g.pushUp(ublock, dblock, fence, col)
			fence++
		}
	}

	g.spawn()
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
		g.blockCnt--
		return
	}
	if ublock != fence-1 {
		g.board[fence-1][col], g.board[ublock][col] = g.board[ublock][col], 0
	}
}

func (g *Game) PushDown() Outcome {
	for col := 0; col < g.Size; col++ {
		fence := g.Size - 1
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

	g.spawn()
	return g.calcOutcome()
}

var (
	texti  bool   // text interface
	player string // player name
	size   int    // board size
	target int    // target
)

func init() {
	flag.IntVar(&size, "size", 4, "Board size: 4 (classic), 5 or 6")
	flag.IntVar(&target, "target", 2048, "End-game block: 2048, 4096 or 8192")
	flag.BoolVar(&texti, "t", true, "Text interface")
	flag.StringVar(&player, "player", "Player", "Player's name")

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	flag.Parse()
	if texti {
		TextGame(player, size, target)
		os.Exit(0)
	}
}

func TextGame(player string, size int, target int) {
	game, err := NewGame(player, size, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in game initialization: %v\n", err)
		os.Exit(1)
	}

	buildTextiParts(player, size)

	reader := bufio.NewReader(os.Stdin)
	outcome := Continue

	fmt.Println("Controls: 'w' (Up), 'a' (Left), 'd' (Right), 's' (Down), 'q' (Quit)")
	fmt.Print(game.String())
	for outcome == Continue {
		// read a command
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		// execute the command
		switch char {
		case 'd', 'D', 'l', 'L':
			outcome = game.PushRight()
			fmt.Printf("%s", game.String())
		case 'a', 'A', 'h', 'H':
			outcome = game.PushLeft()
			fmt.Printf("%s", game.String())
		case 'w', 'W', 'k', 'K':
			outcome = game.PushUp()
			fmt.Printf("%s", game.String())
		case 's', 'S', 'j', 'J':
			outcome = game.PushDown()
			fmt.Printf("%s", game.String())
		case 'e', 'E', 'q', 'Q':
			outcome = GameOver
		}
	}

	if outcome == GameOverWin {
		fmt.Printf("\n%s WINS!\nScore: %d\n", game.Player, game.Score)
	} else {
		fmt.Printf("\nScore: 0\nGAME OVER!\n")
	}
}
