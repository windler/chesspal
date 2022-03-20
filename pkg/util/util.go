package util

import (
	"fmt"

	"github.com/notnil/chess/uci"
)

type EngineOptions struct {
	SkillLevel int
	ELO        int
}

func CreateUCIEngine(engine string, opts EngineOptions, threads int) (*uci.Engine, error) {
	eng, err := uci.New(engine)
	if err != nil {
		return nil, err
	}

	cmds := []uci.Cmd{
		uci.CmdUCI,
		uci.CmdIsReady,
		uci.CmdUCINewGame,
		uci.CmdSetOption{Name: "Threads", Value: fmt.Sprintf("%d", threads)},
	}

	if opts.ELO != 0 {
		cmds = append(cmds, uci.CmdSetOption{Name: "UCI_LimitStrength", Value: "true"})
		cmds = append(cmds, uci.CmdSetOption{Name: "UCI_Elo", Value: fmt.Sprintf("%d", opts.ELO)})
	} else {
		cmds = append(cmds, uci.CmdSetOption{Name: "Skill Level", Value: fmt.Sprintf("%d", opts.SkillLevel)})
	}

	err = eng.Run(
		cmds...,
	)

	if err != nil {
		return nil, err
	}

	return eng, nil
}
