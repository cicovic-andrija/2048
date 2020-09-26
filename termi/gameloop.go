package termi

import (
	"fmt"
	"os"

	"github.com/cicovic-andrija/2048/core"
)

func asserterr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func NewTerminalGraphicsGame(player string, size int, target int, undos int) {
	game, err := core.NewGame(player, size, target, undos)
	asserterr(err)
	termGame, err := NewTermGame(game, 0, 0)
	asserterr(err)
	termGame.Run()
}
