package player

import (
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/windler/chesspal/pkg/util"
)

type UCI struct {
	engine *uci.Engine
	depth  int
	ms     int
}

func NewUCIPlayer(engine string, level int) *UCI {

	ms := 50
	depth := 5
	skill := -9

	// based on lichess
	switch level {
	case 2:
		ms = 100
		depth = 5
		skill = -5
	case 3:
		ms = 150
		depth = 5
		skill = -1
	case 4:
		ms = 200
		depth = 5
		skill = 3
	case 5:
		ms = 300
		depth = 5
		skill = 7
	case 6:
		ms = 400
		depth = 8
		skill = 7
	case 7:
		ms = 500
		depth = 13
		skill = 16
	case 8:
		ms = 1000
		depth = 22
		skill = 20
	}
	eng, err := util.CreateUCIEngine(engine, skill, 4)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UCI player created with engine %s and level %d", engine, skill)

	return &UCI{
		engine: eng,
		depth:  depth,
		ms:     ms,
	}
}

func (p *UCI) MakeMove(game *chess.Game) {
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Duration(p.ms) * time.Millisecond, Depth: p.depth}

	if err := p.engine.Run(cmdPos, cmdGo); err != nil {
		log.Fatal(err)
	}
	move := p.engine.SearchResults().BestMove
	if err := game.Move(move); err != nil {
		log.Fatal(err)
	}

}

func (p *UCI) SetColor(color chess.Color) {

}

func (p *UCI) End() {
	p.engine.Close()
}
