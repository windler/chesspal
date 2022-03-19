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
}

func NewUCIPlayer(engine string, skillLevel int) *UCI {
	eng, err := util.CreateUCIEngine(engine, skillLevel, 4)
	if err != nil {
		log.Fatal(err)
	}

	return &UCI{
		engine: eng,
	}
}

func (p *UCI) MakeMove(game *chess.Game) {
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: 500 * time.Millisecond}

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
