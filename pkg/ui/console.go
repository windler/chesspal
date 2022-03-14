package ui

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/notnil/chess"
	"github.com/windler/chesspal/pkg/game"
)

type Console struct {
	actionChannel  chan game.UIAction
	lastPosition   string
	lastEvaluation game.EvalResult
	lastResult     string
	writer         *uilive.Writer
}

func NewConsoleUI() *Console {
	return &Console{
		actionChannel: make(chan game.UIAction),
		writer:        uilive.New(),
	}
}

func (c *Console) GetActionChannel() chan game.UIAction {
	return c.actionChannel
}

func (c *Console) Start(ctx context.Context, g chess.Game) {
	go func(c *Console, g chess.Game) {
		c.writer.Start()
		for {
			select {
			case action := <-c.actionChannel:
				if action.Move != nil {
					g.Move(action.Move)
					c.lastPosition = fmt.Sprintf("%s\n%s was played", g.Position().Board().Draw(), action.Move.String())
				}

				if g.Position().Status() != chess.NoMethod {
					c.lastResult = fmt.Sprintf("status: %s\n", g.Position().Status())
				}
				if g.Outcome() != chess.NoOutcome {
					c.lastResult = fmt.Sprintf("result: %s\n", g.Outcome())
				}

				if action.Evaluation != nil {
					c.lastEvaluation = *action.Evaluation

					acc := ""
					switch c.lastEvaluation.Accuracy {
					case game.EVAL_ACC_BLUNDER:
						acc = color.RedString(fmt.Sprintf("%v. %s was best", c.lastEvaluation.Accuracy, c.lastEvaluation.BestMove))
					case game.EVAL_ACC_INACCURATE:
						acc = fmt.Sprintf("%v. %s was best", c.lastEvaluation.Accuracy, c.lastEvaluation.BestMove)
					case game.EVAL_ACC_MISTAKE:
						acc = color.YellowString(fmt.Sprintf("%v. %s was best", c.lastEvaluation.Accuracy, c.lastEvaluation.BestMove))
					}
					if c.lastEvaluation.IsForcedMate {
						fmt.Fprintf(c.writer, "%s\nForced mate in %d\n%s\n", c.lastPosition, action.Evaluation.ForcedMateIn, c.lastResult)
					} else {
						fmt.Fprintf(c.writer, "%s\nEvaluation: %.2f\n%v\n\n%s", c.lastPosition, action.Evaluation.Pawn, acc, c.lastResult)
					}
				} else {
					fmt.Fprintf(c.writer, "%s\n%s\n", c.lastPosition, c.lastResult)
				}
			case <-ctx.Done():
				c.writer.Stop()
				close(c.actionChannel)

				return
			}
		}
	}(c, g)
}
