package util

import (
	"fmt"

	"github.com/notnil/chess/uci"
)

func CreateUCIEngine(engine string, skill, threads int) (*uci.Engine, error) {
	eng, err := uci.New(engine)
	if err != nil {
		return nil, err
	}

	err = eng.Run(
		uci.CmdUCI,
		uci.CmdIsReady,
		uci.CmdUCINewGame,
		uci.CmdSetOption{Name: "Skill Level", Value: fmt.Sprintf("%d", skill)},
		uci.CmdSetOption{Name: "Threads", Value: fmt.Sprintf("%d", threads)},
	)

	if err != nil {
		return nil, err
	}

	return eng, nil
}
