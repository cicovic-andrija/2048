package main

import (
	"flag"
	"fmt"

	"github.com/cicovic-andrija/2048/termi"
	"github.com/cicovic-andrija/2048/texti"
)

var (
	// used in main
	local         bool // local game
	hosted        bool // hosted game
	textinterface bool // text interface
	terminterface bool // terminal interface

	// passed and validated later in other packages
	player string // player name
	size   int    // board size
	target int    // end-game block
	undos  int    // number of undos
)

func init() {
	flag.BoolVar(&local, "local", true, "Local game")
	flag.BoolVar(&hosted, "hosted", false, "Hosted game (overrides --local)")
	flag.BoolVar(&terminterface, "termi", true, "Terminal graphics")
	flag.BoolVar(&textinterface, "texti", false, "Text interface (overrides --terminterface)")
	flag.StringVar(&player, "player", "Player", "Player's `name`")
	flag.IntVar(&size, "size", 4, "Board size: 4 (classic), 5 or 6")
	flag.IntVar(&target, "target", 2048, "End-game `block`: 2048, 4096 or 8192")
	flag.IntVar(&undos, "undos", 3, "Number of undos")
}

func parseCmdline() {
	flag.Parse()

	if textinterface {
		terminterface = false
	}

	if hosted {
		local = false
	}
}

func main() {
	parseCmdline()

	if local && textinterface {
		texti.NewTextGame(player, size, target, undos)
	}

	if local && terminterface {
		termi.NewTerminalGraphicsGame()
	}

	if hosted {
		fmt.Printf("Hosted games are not yet supported.\n")
	}
}
