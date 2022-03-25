package eval

import (
	"math"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/windler/chesspal/pkg/game"

	"github.com/windler/chesspal/pkg/util"
)

type LastMove struct {
	engine        *uci.Engine
	cps           map[string]float64
	wasForcedMate bool
}

func NewLastMoveEval(engine string) *LastMove {
	eng, err := util.CreateUCIEngine(engine, util.EngineOptions{
		SkillLevel: 20,
	}, 8)
	if err != nil {
		panic(err)
	}
	return &LastMove{
		engine: eng,
		cps:    map[string]float64{},
	}
}

func (e *LastMove) Eval(g *chess.Game) game.EvalResult {
	result := game.EvalResult{}

	e.engine.Run(uci.CmdStop)

	move := g.Moves()[len(g.Moves())-1]

	e.engine.Run(uci.CmdPosition{Position: g.Position()}, uci.CmdGo{Depth: 22, MoveTime: 500 * time.Millisecond})
	cp := float64(e.engine.SearchResults().Info.Score.CP) / 100.0
	if g.Position().Turn() == chess.Black {
		cp = cp * -1
	}

	cpPreviousMove := float64(0)
	if len(g.Moves()) > 1 {
		cpPreviousMove = e.cps[g.Moves()[len(g.Moves())-2].String()]
	}
	if e.engine.SearchResults().Info.Score.Mate == 0 {
		centiPawnLoss := math.Abs(cpPreviousMove - cp)
		acc := ""
		if e.wasForcedMate {
			e.wasForcedMate = false
			acc = string(game.EVAL_ACC_INACCURATE)
		} else if math.Abs(cp) > 6 {
			// this is fine
		} else if centiPawnLoss >= 3 {
			acc = string(game.EVAL_ACC_BLUNDER)
		} else if centiPawnLoss >= 2 {
			acc = string(game.EVAL_ACC_MISTAKE)
		} else if centiPawnLoss >= 1 {
			acc = string(game.EVAL_ACC_INACCURATE)
		}

		if acc != "" {
			g.AddComment(move, acc)
		}

		result.Pawn = cp
	} else {
		result.IsForcedMate = true
		result.ForcedMateIn = e.engine.SearchResults().Info.Score.Mate
		if g.Position().Turn() == chess.Black {
			result.ForcedMateIn = result.ForcedMateIn * -1
		}
		e.wasForcedMate = true
	}
	e.cps[g.Moves()[len(g.Moves())-1].String()] = cp

	return result
}
