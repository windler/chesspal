package ui

import (
	"fmt"
	"sync"

	"github.com/gosuri/uilive"
	"github.com/notnil/chess"
	"github.com/windler/chesspal/pkg/game"
)

type Console struct {
	actionChannel  chan game.UIAction
	lastPosition   string
	lastEvaluation game.EvalResult
	lastResult     string
	lastStatus     string
	writer         *uilive.Writer
	mutex          *sync.Mutex
}

func NewConsoleUI() *Console {
	console := &Console{
		actionChannel: make(chan game.UIAction),
		writer:        uilive.New(),
		mutex:         &sync.Mutex{},
	}

	return console
}

func (c *Console) GetActionChannel() chan game.UIAction {
	return c.actionChannel
}

func (c *Console) Render(g chess.Game, action game.UIAction) {
	c.mutex.Lock()

	c.writer.Start()
	if action.Move != nil {
		g.Move(action.Move)
		c.lastPosition = fmt.Sprintf("%s\n%s was played", g.Position().Board().Draw(), action.Move.String())
	}

	if g.Position().Status() != chess.NoMethod {
		c.lastStatus = fmt.Sprintf("%s\n", g.Position().Status())
	}
	if g.Outcome() != chess.NoOutcome {
		c.lastResult = fmt.Sprintf("result: %s\n", g.Outcome())
	}

	if action.Evaluation != nil {
		c.lastEvaluation = *action.Evaluation

		acc := ""
		// switch c.lastEvaluation.Accuracy {
		// case game.EVAL_ACC_BLUNDER:
		// 	acc = color.RedString(fmt.Sprintf("%v", c.lastEvaluation.Accuracy))
		// case game.EVAL_ACC_INACCURATE:
		// 	acc = fmt.Sprintf("%v", c.lastEvaluation.Accuracy)
		// case game.EVAL_ACC_MISTAKE:
		// 	acc = color.YellowString(fmt.Sprintf("%v", c.lastEvaluation.Accuracy))
		// }
		if c.lastEvaluation.IsForcedMate {
			fmt.Fprintf(c.writer, "%s\nForced mate in %d\n%s%s\n", c.lastPosition, action.Evaluation.ForcedMateIn, c.lastStatus, c.lastResult)
		} else {
			fmt.Fprintf(c.writer, "%s\nEvaluation: %.2f\n%v\n\n%s%s\n", c.lastPosition, action.Evaluation.Pawn, acc, c.lastStatus, c.lastResult)
		}
	} else {
		fmt.Fprintf(c.writer, "%s\n%s%s\n", c.lastPosition, c.lastStatus, c.lastResult)
	}

	c.writer.Stop()
	c.mutex.Unlock()
}
