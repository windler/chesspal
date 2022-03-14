package player

import (
	"fmt"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

type UCI struct {
	engine     *uci.Engine
	skillLevel int
}

func NewUCIPlayer(engine string, skillLevel int) *UCI {
	eng, err := uci.New(engine)
	if err != nil {
		panic(err)
	}

	// uci.CmdSetOption{Name: "MultiPV", Value: "3"}
	err = eng.Run(
		uci.CmdUCI,
		uci.CmdIsReady,
		uci.CmdUCINewGame,
		uci.CmdSetOption{Name: "Skill Level", Value: fmt.Sprintf("%d", skillLevel)},
		uci.CmdSetOption{Name: "Threads", Value: "8"},
	)

	if err != nil {
		panic(err)
	}

	return &UCI{
		engine:     eng,
		skillLevel: skillLevel,
	}
}

func (p *UCI) MakeMove(game *chess.Game) {
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: 500 * time.Millisecond}

	if err := p.engine.Run(cmdPos, cmdGo); err != nil {
		//TODO
		panic(err)
	}
	move := p.engine.SearchResults().BestMove
	if err := game.Move(move); err != nil {
		//TODO
		panic(err)
	}

}

func (p *UCI) End() {
	p.engine.Close()
}
