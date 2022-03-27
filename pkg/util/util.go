package util

import (
	"fmt"
	"strings"

	"github.com/notnil/chess/uci"
)

func CreateUCIEngine(engine string, opts []string, threads int) (*uci.Engine, error) {
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

	for _, o := range opts {
		split := strings.Split(o, "=")
		cmds = append(cmds, uci.CmdSetOption{Name: split[0], Value: split[1]})
	}

	err = eng.Run(
		cmds...,
	)

	if err != nil {
		return nil, err
	}

	return eng, nil
}
