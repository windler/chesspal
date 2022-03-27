package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notnil/chess"
	"github.com/windler/chesspal/pkg/eval"
	"github.com/windler/chesspal/pkg/game"
	"github.com/windler/chesspal/pkg/player"
	"github.com/windler/chesspal/pkg/ui"
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

var started = false
var g *game.Game
var engine *player.DGTEngine
var currentBoard chess.Board

func main() {
	wsUI := ui.NewWS()

	engine = player.NewDGTEngine()
	engine.Start()

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

	e.Static("/", "./web/vue-frontend/dist/")

	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer ws.Close()

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
					go startGame(msg, wsUI)
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

	e.Logger.Fatal(e.Start(":8080"))

}

type Started struct {
	Started bool `json:"started"`
}

func sendGameStarted(ws *websocket.Conn) {
	if err := ws.WriteJSON(&Started{Started: true}); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
}

func startGame(msg *Message, ui game.UI) {
	evals := []game.EvalEngine{}

	engine.SetUpsideDown(msg.Options.UpsideDown)

	var white, black game.Player
	if msg.Options.Black.Type == 0 {
		black = player.NewDGTPlayer(engine)
	} else {
		black = player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish_14", msg.Options.Black.Type)
	}
	if msg.Options.White.Type == 0 {
		white = player.NewDGTPlayer(engine)
	} else {
		white = player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish_14", msg.Options.White.Type)
	}

	g = game.NewGame(black, white, ui)

	if msg.Options.EvalMode == 1 {
		evals = append(evals, eval.NewLastMoveEval("/home/windler/projects/chess/chesspal/bin/stockfish_14"))
	}

	g.Start(currentBoard.String(), evals...)
}
