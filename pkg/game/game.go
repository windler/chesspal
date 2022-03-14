package game

import (
	"context"
	"sync"

	"github.com/notnil/chess"
)

type Game struct {
	black Player
	white Player
	uis   []UI
	game  *chess.Game
}

type EvalEngine interface {
	Eval(chess.Game) EvalResult
}

type EvalResult struct {
	Pawn         float64
	Accuracy     EvalAccuracy
	BestMove     string
	IsForcedMate bool
	ForcedMateIn int
}

type EvalAccuracy string

const (
	EVAL_ACC_GOOD       EvalAccuracy = "Good"
	EVAL_ACC_INACCURATE EvalAccuracy = "Inaccuracy"
	EVAL_ACC_MISTAKE    EvalAccuracy = "Mistake"
	EVAL_ACC_BLUNDER    EvalAccuracy = "Blunder"
	EVAL_ACC_BEST       EvalAccuracy = "Best"
)

type Player interface {
	MakeMove(*chess.Game)
	End()
}

type UI interface {
	Start(context.Context, chess.Game)
	GetActionChannel() chan UIAction
}

type UIAction struct {
	Move       *chess.Move
	Evaluation *EvalResult
}

func NewGame(black, white Player, uis ...UI) *Game {
	return &Game{
		black: black,
		white: white,
		uis:   uis,
	}
}

func (g *Game) Start(evalEngines ...EvalEngine) {
	g.game = chess.NewGame()

	wg := &sync.WaitGroup{}

	ctx, cancelFunc := context.WithCancel(context.Background())

	for _, ui := range g.uis {
		ui.Start(ctx, *g.game)
	}

	for g.game.Outcome() == chess.NoOutcome {
		if g.game.Position().Turn() == chess.Black {
			g.black.MakeMove(g.game)
		} else {
			g.white.MakeMove(g.game)
		}

		move := g.game.Moves()[len(g.game.Moves())-1]
		g.callEvalEngines(evalEngines, wg)
		g.callUIs(UIAction{
			Move: move,
		})
	}

	wg.Wait()
	cancelFunc()

	g.black.End()
	g.white.End()
}

func (g *Game) callUIs(action UIAction) {
	for _, ui := range g.uis {
		ui.GetActionChannel() <- action
	}
}

func (g *Game) callEvalEngines(engines []EvalEngine, wg *sync.WaitGroup) {
	for _, engine := range engines {
		wg.Add(1)

		go func(engine EvalEngine, game chess.Game) {
			defer wg.Done()
			evaluation := engine.Eval(game)

			g.callUIs(UIAction{
				Evaluation: &evaluation,
			})
		}(engine, *g.game)
	}
}
