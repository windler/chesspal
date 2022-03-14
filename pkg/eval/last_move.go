package eval

import (
	"math"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/windler/chesspal/pkg/game"
)

type LastMove struct {
	engine    *uci.Engine
	cpCurrent float64
}

func NewLastMoveEval(engine string) *LastMove {
	eng, err := uci.New(engine)
	if err != nil {
		panic(err)
	}

	err = eng.Run(
		uci.CmdUCI,
		uci.CmdIsReady,
		uci.CmdUCINewGame,
		uci.CmdSetOption{Name: "Skill Level", Value: "20"},
		uci.CmdSetOption{Name: "Threads", Value: "8"},
	)

	if err != nil {
		panic(err)
	}

	return &LastMove{
		engine: eng,
	}
}

func (e *LastMove) Eval(g chess.Game) game.EvalResult {
	result := game.EvalResult{}

	e.engine.Run(uci.CmdStop)

	e.engine.Run(uci.CmdPosition{Position: g.Position()}, uci.CmdGo{Depth: 17})
	cp := float64(e.engine.SearchResults().Info.Score.CP) / 100.0
	if g.Position().Turn() == chess.Black {
		cp = cp * -1
	}

	if e.engine.SearchResults().Info.Score.Mate == 0 {
		centiPawnLoss := math.Abs(math.Abs(e.cpCurrent) - math.Abs(cp))
		if centiPawnLoss >= 3 {
			result.Accuracy = game.EVAL_ACC_BLUNDER
		} else if centiPawnLoss >= 2 {
			result.Accuracy = game.EVAL_ACC_MISTAKE
		} else if centiPawnLoss >= 1 {
			result.Accuracy = game.EVAL_ACC_INACCURATE
		}

		result.BestMove = e.engine.SearchResults().BestMove.String()
		result.Pawn = cp
	} else {
		result.IsForcedMate = true
		result.ForcedMateIn = e.engine.SearchResults().Info.Score.Mate
	}
	e.cpCurrent = cp

	return result
}
