package texti

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicovic-andrija/term2048/core"
)

var (
	textiHorizLine    string
	textiScoreLineFmt string
)

// assumes board size is in limits
func buildTextiParts(playerName string, boardSize int) {
	textiHorizLine = "\n+" + strings.Repeat("------+", boardSize) + "\n"
	textiScoreLineFmt = playerName + "'s score: %d"
}

func drawBoard(g *core.Game) {
	tostring := func(v int) string {
		if v == 0 {
			return " "
		}
		return strconv.Itoa(v)
	}

	var str strings.Builder
	str.WriteString(fmt.Sprintf(textiScoreLineFmt, g.Score))
	str.WriteString(textiHorizLine)
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			str.WriteString(fmt.Sprintf("| %-4s ", tostring(g.Cell(i, j))))
		}
		str.WriteString("|" + textiHorizLine)
	}

	fmt.Print(str.String())
}

func NewTextGame(player string, size int, target int) {
	game, err := core.NewGame(player, size, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in game initialization: %v\n", err)
		os.Exit(1)
	}

	buildTextiParts(player, size)

	reader := bufio.NewReader(os.Stdin)
	outcome := core.Continue

	fmt.Println("Controls: 'w' (Up), 'a' (Left), 'd' (Right), 's' (Down), 'q' (Quit)")
	drawBoard(game)
	for outcome == core.Continue {

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
			drawBoard(game)
		case 'a', 'A', 'h', 'H':
			outcome = game.PushLeft()
			drawBoard(game)
		case 'w', 'W', 'k', 'K':
			outcome = game.PushUp()
			drawBoard(game)
		case 's', 'S', 'j', 'J':
			outcome = game.PushDown()
			drawBoard(game)
		case 'e', 'E', 'q', 'Q':
			outcome = core.GameOver
		}
	}

	if outcome == core.GameOverWin {
		fmt.Printf("\n%s WINS!\nScore: %d\n", game.Player, game.Score)
	} else {
		fmt.Printf("\nScore: 0\nGAME OVER!\n")
	}
}