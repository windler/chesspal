package game

import (
	"github.com/notnil/chess"
)

type Game struct {
	black Player
	white Player
	uis   []UI
	game  *chess.Game
}

type EvalEngine interface {
	Eval(*chess.Game) EvalResult
}

type EvalResult struct {
	Pawn         float64
	BestMoves    []chess.Move
	IsForcedMate bool
	ForcedMateIn int
}

type EvalAccuracy string

const (
	EVAL_ACC_INACCURATE EvalAccuracy = "Inaccuracy"
	EVAL_ACC_MISTAKE    EvalAccuracy = "Mistake"
	EVAL_ACC_BLUNDER    EvalAccuracy = "Blunder"
)

type Player interface {
	MakeMove(*chess.Game)
	SetColor(chess.Color)
	End()
}

type UI interface {
	Render(chess.Game, UIAction)
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
	g.black.SetColor(chess.Black)
	g.white.SetColor(chess.White)

	g.callUIs(UIAction{})

	for g.game.Outcome() == chess.NoOutcome {
		if g.game.Position().Turn() == chess.Black {
			g.black.MakeMove(g.game)
		} else {
			g.white.MakeMove(g.game)
		}

		move := g.game.Moves()[len(g.game.Moves())-1]
		g.callEvalEngines(evalEngines)
		g.callUIs(UIAction{
			Move: move,
		})
	}

	g.black.End()
	g.white.End()
}

func (g *Game) callUIs(action UIAction) {
	for _, ui := range g.uis {
		ui.Render(*g.game, action)
	}
}

func (g *Game) callEvalEngines(engines []EvalEngine) {
	for _, engine := range engines {

		func(engine EvalEngine, game *chess.Game) {
			evaluation := engine.Eval(game)

			g.callUIs(UIAction{
				Evaluation: &evaluation,
			})
		}(engine, g.game)
	}
}

func (g *Game) UndoMoves(n int) error {
	err := g.game.UndoMoves(n)
	move := g.game.Moves()[len(g.game.Moves())-1]
	g.callUIs(UIAction{
		Move: move,
	})

	return err
}

func (g *Game) Draw() {
	g.game.Draw(chess.DrawOffer)
	move := g.game.Moves()[len(g.game.Moves())-1]
	g.callUIs(UIAction{
		Move: move,
	})
}

func (g *Game) Resign() {
	g.game.Resign(g.game.Position().Turn())
	move := g.game.Moves()[len(g.game.Moves())-1]
	g.callUIs(UIAction{
		Move: move,
	})
}
