package texti

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cicovic-andrija/2048/core"
)

var (
	textiHorizLine    string
	textiScoreLineFmt string
)

// assumes board size is in limits
func buildTextiParts(playerName string, boardSize int) {
	textiHorizLine = "\n+" + strings.Repeat("------+", boardSize) + "\n"
	textiScoreLineFmt = playerName + "'s score: %d, undos left: %d"
}

func drawBoard(g *core.Game) {
	tostring := func(v int) string {
		if v == 0 {
			return " "
		}
		return strconv.Itoa(v)
	}

	var str strings.Builder
	str.WriteString(fmt.Sprintf(textiScoreLineFmt, g.Score(), g.UndosLeft()))
	str.WriteString(textiHorizLine)
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			str.WriteString(fmt.Sprintf("| %-4s ", tostring(g.Block(i, j))))
		}
		str.WriteString("|" + textiHorizLine)
	}

	fmt.Print(str.String())
}

func NewTextGame(player string, size int, target int, undos int) error {
	game, err := core.NewGame(player, size, target, undos)
	if err != nil {
		return fmt.Errorf("error in game initialization: %v", err)
	}

	buildTextiParts(player, size)

	reader := bufio.NewReader(os.Stdin)
	outcome := core.Continue

	fmt.Println("Controls: 'w' (Up) / 'a' (Left) / 'd' (Right) / 's' (Down) / 'u' (Undo) / 'q' (Quit)")
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
			outcome = game.Push(core.Right)
			drawBoard(game)
		case 'a', 'A', 'h', 'H':
			outcome = game.Push(core.Left)
			drawBoard(game)
		case 'w', 'W', 'k', 'K':
			outcome = game.Push(core.Up)
			drawBoard(game)
		case 's', 'S', 'j', 'J':
			outcome = game.Push(core.Down)
			drawBoard(game)
		case 'u', 'U':
			if ok := game.Undo(); !ok {
				fmt.Println("Can't undo: no undos left / second undo in a row / first move.")
				break
			}
			drawBoard(game)
		case 'q', 'Q':
			outcome = core.GameOver
		case '\n', '\r':
			// ignore
		default:
			fmt.Println("Invalid command.")
		}
	}

	if outcome == core.GameOverWin {
		fmt.Printf("===\n%s WINS! Score: %d\n===\n", game.Player, game.Score())
	} else {
		fmt.Printf("===\nGAME OVER! Score: 0\n===\n")
	}

	return nil
}
