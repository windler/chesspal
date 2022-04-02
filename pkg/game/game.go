package game

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

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
	Name() string
	IsBot() bool
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

func (g *Game) Start(fenString string, evalEngines ...EvalEngine) {
	// TODO check castling availability
	fen, err := chess.FEN(fmt.Sprintf("%s w KQkq - 0 1", fenString))
	if err != nil {
		panic(err)
	}
	g.game = chess.NewGame(fen)

	g.game.AddTagPair("White", g.white.Name())
	g.game.AddTagPair("Black", g.black.Name())
	g.game.AddTagPair("Date", time.Now().Format(time.RFC822))

	if g.white.IsBot() || g.black.IsBot() {
		g.game.AddTagPair("Botgame", "true")
	}

	g.black.SetColor(chess.Black)
	g.white.SetColor(chess.White)

	g.callUIs(UIAction{})

	wg := &sync.WaitGroup{}
	go func() {
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
	}()
	go func() {
		for g.game.Outcome() == chess.NoOutcome {
			time.Sleep(500 * time.Millisecond)
		}
		wg.Done()
	}()
	wg.Add(1)
	wg.Wait()

	g.game.AddTagPair("Result", g.game.Outcome().String())
	g.callUIs(UIAction{})

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

func (g *Game) Save(folder string) {
	file := fmt.Sprintf("%s%d_%s_vs_%s.pgn", folder, time.Now().UnixMilli(), g.game.GetTagPair("White").Value, g.game.GetTagPair("Black").Value)
	ioutil.WriteFile(file, []byte(g.game.String()), 0644)
}
