package main

import (
	"github.com/windler/chesspal/pkg/eval"
	"github.com/windler/chesspal/pkg/game"
	"github.com/windler/chesspal/pkg/player"
	"github.com/windler/chesspal/pkg/ui"
)

func main() {
	black := player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish", 20)
	white := player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish", 1)

	game := game.NewGame(black, white, ui.NewConsoleUI())

	game.Start(eval.NewLastMoveEval("/home/windler/projects/chess/chesspal/bin/stockfish"))
}
