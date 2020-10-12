package termi

import (
	"github.com/cicovic-andrija/2048/core"
)

func NewTerminalGraphicsGame(player string, size int, target int, undos int) error {
	game, err := core.NewGame(player, size, target, undos)
	if err != nil {
		return err
	}

	termGame, err := NewTermGame(game, 0, 0)
	if err != nil {
		return err
	}

	return termGame.Run()
}
