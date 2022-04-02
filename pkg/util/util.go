package util

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
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

func GetSVG(board chess.Board, opts ...func(*image.Encoder)) string {
	colors := image.SquareColors(color.RGBA{R: 0xC7, G: 0xC6, B: 0xC1}, color.RGBA{R: 0x82, G: 0x82, B: 0x82})
	o := []func(*image.Encoder){colors}

	for _, opt := range opts {
		o = append(o, opt)
	}

	buf := bytes.NewBufferString("")

	if err := image.SVG(buf, &board, o...); err != nil {
		log.Printf("error occurred: %v", err)
	}

	return buf.String()
}
