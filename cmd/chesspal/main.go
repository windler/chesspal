package main

import (
	"bytes"
	"errors"
	"image/color"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
	"github.com/windler/chesspal/pkg/eval"
	"github.com/windler/chesspal/pkg/game"
	"github.com/windler/chesspal/pkg/player"
	"github.com/windler/chesspal/pkg/ui"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Action string `json:"action"`
}

type StartMessage struct {
	Message
	Options StartOptions `json:"options"`
}

type StartOptions struct {
	White    Player `json:"white"`
	Black    Player `json:"black"`
	EvalMode int    `json:"evalMode"`
}

type Player struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

var started = false

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "./web/frontend/public/")

	e.GET("/ws", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if !errors.Is(err, nil) {
			log.Println(err)
		}
		defer ws.Close()

		for {
			startMsg := &StartMessage{}
			err := ws.ReadJSON(&startMsg)
			if !errors.Is(err, nil) {
				log.Printf("error occurred: %v", err)
				break
			}

			if !started {
				go func() {
					var engine *player.DGTEngine
					evals := []game.EvalEngine{}

					if startMsg.Options.Black.Type == 0 || startMsg.Options.White.Type == 0 {
						engine = player.NewDGTEngine()
						engine.Start()
					}
					var white, black game.Player
					if startMsg.Options.Black.Type == 0 {
						black = player.NewDGTPlayer(engine)
					} else {
						black = player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish", startMsg.Options.Black.Type)
					}
					if startMsg.Options.White.Type == 0 {
						white = player.NewDGTPlayer(engine)
					} else {
						white = player.NewUCIPlayer("/home/windler/projects/chess/chesspal/bin/stockfish", startMsg.Options.White.Type)
					}

					game := game.NewGame(black, white, &WSUI{
						ws:    ws,
						mutex: &sync.Mutex{},
					}, ui.NewConsoleUI())

					if startMsg.Options.EvalMode == 1 {
						evals = append(evals, eval.NewLastMoveEval("/home/windler/projects/chess/chesspal/bin/stockfish"))
					}

					game.Start(evals...)

				}()
				started = true
			}
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))

}

type WSUI struct {
	ws    *websocket.Conn
	mutex *sync.Mutex
}

type GameState struct {
	SVGPosition string  `json:"svgPosition"`
	Pawn        float64 `json:"pawn"`
	Accuracy    string  `json:"accuracy"`
	LastMove    string  `json:"lastMove"`
	Turn        string  `json:"turn"`
	PGN         string  `json:"pgn"`
}

var yellow = color.RGBA{255, 255, 0, 1}

var currentState = &GameState{}

func (u *WSUI) Render(game chess.Game, action game.UIAction) {
	u.mutex.Lock()
	if currentState.SVGPosition == "" {
		buf := bytes.NewBufferString("")
		if err := image.SVG(buf, game.Position().Board()); err != nil {
			log.Printf("error occurred: %v", err)
		}
		currentState.SVGPosition = buf.String()
	}
	if action.Move != nil {
		buf := bytes.NewBufferString("")
		mark := image.MarkSquares(yellow, action.Move.S1(), action.Move.S2())
		if err := image.SVG(buf, game.Position().Board(), mark); err != nil {
			log.Printf("error occurred: %v", err)
		}
		currentState.SVGPosition = buf.String()
		currentState.LastMove = action.Move.String()
		currentState.Turn = game.Position().Turn().String()
		currentState.PGN = game.String()
	}

	if action.Evaluation != nil {
		currentState.Pawn = action.Evaluation.Pawn + 50
		currentState.Accuracy = string(action.Evaluation.Accuracy)
		if action.Evaluation.IsForcedMate {
			currentState.Pawn = 100
			if action.Evaluation.ForcedMateIn < 0 {
				currentState.Pawn = 0
			}
		}
	}

	if err := u.ws.WriteJSON(currentState); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
	u.mutex.Unlock()
}
