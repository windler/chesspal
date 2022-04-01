package player

import (
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/windler/chesspal/pkg/util"
)

type BotOptions struct {
	Name       string   `yaml:"name" json:"name"`
	Engine     string   `yaml:"engine" json:"-"`
	Path       string   `yaml:"-" json:"-"`
	Depth      int      `yaml:"depth" json:"-"`
	MoveTimeMs int      `yaml:"moveTimeMs" json:"-"`
	Threads    int      `yaml:"threads" json:"-"`
	Options    []string `yaml:"options" json:"-"`
}

type UCI struct {
	engine *uci.Engine
	depth  int
	ms     int
	name   string
}

func (p *UCI) Name() string {
	return p.name
}

func NewUCIPlayer(options BotOptions) *UCI {
	eng, err := util.CreateUCIEngine(options.Path, options.Options, options.Threads)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UCI player created with name %s", options.Name)

	return &UCI{
		engine: eng,
		depth:  options.Depth,
		ms:     options.MoveTimeMs,
		name:   options.Name,
	}
}

func (p *UCI) MakeMove(game *chess.Game) {
	cmds := []uci.Cmd{uci.CmdPosition{Position: game.Position()}, uci.CmdGo{MoveTime: time.Duration(p.ms) * time.Millisecond, Depth: p.depth}}

	if err := p.engine.Run(cmds...); err != nil {
		log.Fatal(err)
	}
	move := p.engine.SearchResults().BestMove
	if err := game.Move(move); err != nil {
		log.Fatal(err)
	}

}

func (p *UCI) SetColor(color chess.Color) {

}

func (p *UCI) End() {
	p.engine.Close()
}
