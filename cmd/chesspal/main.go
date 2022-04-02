package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
	"github.com/windler/chesspal/pkg/eval"
	"github.com/windler/chesspal/pkg/game"
	"github.com/windler/chesspal/pkg/player"
	"github.com/windler/chesspal/pkg/ui"
	"github.com/windler/chesspal/pkg/util"
	"gopkg.in/yaml.v3"
)

var upgrader = websocket.Upgrader{}

const (
	MSG_START        string = "start"
	MSG_UNDO_N_MOVES string = "undo"
	MSG_SET_RESULT   string = "result"
)

type Message struct {
	Action    string       `json:"action"`
	Options   StartOptions `json:"startOptions"`
	UndoMoves int          `json:"undoMoves"`
	Result    string       `json:"result"`
}

type StartOptions struct {
	White      Player `json:"white"`
	Black      Player `json:"black"`
	EvalMode   int    `json:"evalMode"`
	UpsideDown bool   `json:"upsideDown"`
}

type Player struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Config struct {
	Address       string              `yaml:"address"`
	Web           string              `yaml:"web"`
	GamesFolder   string              `yaml:"gamesFolder"`
	ArchiveFolder string              `yaml:"archiveFolder"`
	DgtPort       string              `yaml:"dgtPort"`
	Engines       map[string]string   `yaml:"engines"`
	Bots          []player.BotOptions `yaml:"bots"`
	Eval          Eval                `yaml:"eval"`
}

type Eval struct {
	Engine     string   `yaml:"engine"`
	Depth      int      `yaml:"depth"`
	Threads    int      `yaml:"threads"`
	MoveTimeMs int      `yaml:"moveTimeMs"`
	Options    []string `yaml:"options"`
}

var started = false
var g *game.Game
var engine *player.DGTEngine
var currentBoard chess.Board

type WSResponse struct {
	Bots []player.BotOptions `json:"bots"`
}

type GameHistory struct {
	Games []Game `json:"games"`
}

type Game struct {
	ID       string `json:"id"`
	PGN      string `json:"pgn"`
	SVG      string `json:"svg"`
	White    string `json:"white"`
	Black    string `json:"black"`
	Date     string `json:"date"`
	Result   string `json:"result"`
	Archived bool   `json:"archived"`
	Botgame  bool   `json:"botgame"`
}

var yellow = color.RGBA{255, 255, 0, 1}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "/home/windler/projects/chess/chesspal/configs/chesspal.yaml", "Path to config")

	filename, _ := filepath.Abs(configPath)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}

	wsUI := ui.NewWS()

	engine = player.NewDGTEngine()
	engine.Start(config.DgtPort)

	go func() {
		for board := range engine.PostionChannel() {
			if !started {
				currentBoard = board
				wsUI.SendBoard(board)
			}
		}
	}()
	engine.ReadCurrentPosition()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Static("/", config.Web)

	e.DELETE("/history/:id", func(c echo.Context) error {
		file := c.Param("id")

		err := os.Remove(fmt.Sprintf("%s%s", config.GamesFolder, file))
		if err != nil {
			err := os.Remove(fmt.Sprintf("%s%s", config.ArchiveFolder, file))
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	e.POST("/history/:id/archive", func(c echo.Context) error {
		file := c.Param("id")
		fileName := fmt.Sprintf("%s%s", config.GamesFolder, file)
		archive := fmt.Sprintf("%s%s", config.ArchiveFolder, file)

		err := os.Rename(fileName, archive)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	e.GET("/history", func(c echo.Context) error {
		files, err := ioutil.ReadDir(config.GamesFolder)
		if err != nil {
			panic(err)
		}

		games := []Game{}
		for _, file := range files {
			g := getGame(file, *config, false)
			if g != nil {
				games = append(games, *g)
			}
		}

		filesArchive, err := ioutil.ReadDir(config.ArchiveFolder)
		if err != nil {
			panic(err)
		}

		for _, file := range filesArchive {
			g := getGame(file, *config, true)
			if g != nil {
				games = append(games, *g)
			}
		}

		c.JSON(http.StatusOK, GameHistory{Games: games})

		return nil
	})

	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer ws.Close()

		if err := ws.WriteJSON(WSResponse{Bots: config.Bots}); !errors.Is(err, nil) {
			log.Printf("error occurred: %v", err)
		}

		wsUI.AddWebsocket(ws)
		if started {
			sendGameStarted(ws)
		}

		for {
			msg := &Message{}
			err := ws.ReadJSON(&msg)
			if !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
				break
			}

			switch msg.Action {
			case MSG_START:
				if !started {
					go startGame(msg, wsUI, *config)
					started = true
					sendGameStarted(ws)
				}

			case MSG_UNDO_N_MOVES:
				if msg.UndoMoves > 0 {
					g.UndoMoves(msg.UndoMoves)
				}
			case MSG_SET_RESULT:
				switch msg.Result {
				case "draw":
					g.Draw()
				case "resign":
					g.Resign()
				}
			}

		}

		return nil
	})

	e.Logger.Fatal(e.Start(config.Address))

}

func getGame(f fs.FileInfo, config Config, archive bool) *Game {
	if !strings.HasSuffix(f.Name(), "pgn") {
		return nil
	}

	folder := config.GamesFolder
	if archive {
		folder = config.ArchiveFolder
	}
	file := fmt.Sprintf("%s%s", folder, f.Name())
	contents, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer contents.Close()

	pgn, err := chess.PGN(contents)
	if err != nil {
		panic(err)
	}
	g := chess.NewGame(pgn)

	lastMove := g.Moves()[len(g.Moves())-1]
	mark := image.MarkSquares(yellow, lastMove.S1(), lastMove.S2())
	svg := util.GetSVG(*g.Position().Board(), mark)

	botGame := false

	if g.GetTagPair("Botgame") != nil {
		if g.GetTagPair("Botgame").Value == "true" {
			botGame = true
		}
	}

	return &Game{
		ID:       f.Name(),
		PGN:      g.String(),
		SVG:      svg,
		White:    g.GetTagPair("White").Value,
		Black:    g.GetTagPair("Black").Value,
		Result:   string(g.Outcome()),
		Date:     g.GetTagPair("Date").Value,
		Archived: archive,
		Botgame:  botGame,
	}
}

type Started struct {
	Started bool `json:"started"`
}

func sendGameStarted(ws *websocket.Conn) {
	if err := ws.WriteJSON(&Started{Started: true}); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
}

func startGame(msg *Message, ui game.UI, cfg Config) {
	evals := []game.EvalEngine{}

	engine.SetUpsideDown(msg.Options.UpsideDown)

	var white, black game.Player
	if msg.Options.Black.Type == 0 {
		black = player.NewDGTPlayer(msg.Options.Black.Name, engine)
	} else {
		i := msg.Options.Black.Type - 1
		options := cfg.Bots[i]
		options.Path = cfg.Engines[options.Engine]
		black = player.NewUCIPlayer(options)
	}
	if msg.Options.White.Type == 0 {
		white = player.NewDGTPlayer(msg.Options.White.Name, engine)
	} else {
		i := msg.Options.White.Type - 1
		options := cfg.Bots[i]
		options.Path = cfg.Engines[options.Engine]
		white = player.NewUCIPlayer(options)
	}

	g = game.NewGame(black, white, ui)

	if msg.Options.EvalMode == 1 {
		evals = append(evals, eval.NewLastMoveEval(
			cfg.Engines[cfg.Eval.Engine],
			cfg.Eval.Options,
			cfg.Eval.Threads,
			cfg.Eval.Depth,
			cfg.Eval.MoveTimeMs,
		))
	}

	g.Start(currentBoard.String(), evals...)
	g.Save(cfg.GamesFolder)
}
