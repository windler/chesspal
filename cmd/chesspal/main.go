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

func main() {
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
					go func() {
						var engine *player.DGTEngine
						evals := []game.EvalEngine{}

						if msg.Options.Black.Type == 0 || msg.Options.White.Type == 0 {
							engine = player.NewDGTEngine()
							engine.SetUpsideDown(msg.Options.UpsideDown)
							engine.Start()
						}
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

						g = game.NewGame(black, white, &WSUI{
							ws:    ws,
							mutex: &sync.Mutex{},
						})

						if msg.Options.EvalMode == 1 {
							evals = append(evals, eval.NewLastMoveEval("/home/windler/projects/chess/chesspal/bin/stockfish_14"))
						}

						g.Start(evals...)

					}()
					started = true
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

type WSUI struct {
	ws    *websocket.Conn
	mutex *sync.Mutex
}

type GameState struct {
	SVGPosition string  `json:"svgPosition"`
	Pawn        float64 `json:"pawn"`
	Moves       []Move  `json:"moves"`
	Turn        string  `json:"turn"`
	PGN         string  `json:"pgn"`
	FEN         string  `json:"fen"`
	Outcome     string  `json:"outcome"`
}

type Move struct {
	Move     string `json:"move"`
	Accuracy string `json:"accuracy"`
	Color    string `json:"color"`
}

var yellow = color.RGBA{255, 255, 0, 1}

var currentState = &GameState{}
var moveEncoder = chess.AlgebraicNotation{}

func (u *WSUI) Render(g chess.Game, action game.UIAction) {
	u.mutex.Lock()
	colors := image.SquareColors(color.RGBA{R: 0xC7, G: 0xC6, B: 0xC1}, color.RGBA{R: 0x82, G: 0x82, B: 0x82})
	if len(g.Moves()) == 0 {
		buf := bytes.NewBufferString("")
		if err := image.SVG(buf, g.Position().Board(), colors); err != nil {
			log.Printf("error occurred: %v", err)
		}
		currentState.SVGPosition = buf.String()
	} else {
		move := g.Moves()[len(g.Moves())-1]
		buf := bytes.NewBufferString("")
		mark := image.MarkSquares(yellow, move.S1(), move.S2())
		if err := image.SVG(buf, g.Position().Board(), mark, colors); err != nil {
			log.Printf("error occurred: %v", err)
		}
		currentState.SVGPosition = buf.String()

		currentState.Turn = g.Position().Turn().String()
		currentState.PGN = g.String()
	}

	moves := []Move{}
	//TODO add a GetComments function
	for moveIndex, m := range g.Moves() {
		pos := g.Positions()[moveIndex]
		moveEncoded := moveEncoder.Encode(pos, m)

		accuracy := ""
		if len(g.Comments()) >= moveIndex+1 && len(g.Comments()[moveIndex]) > 0 {
			accuracy = g.Comments()[moveIndex][0]
		}

		moves = append(moves, Move{
			Move:     moveEncoded,
			Color:    pos.Turn().String(),
			Accuracy: accuracy,
		})
	}

	currentState.Moves = moves

	if action.Evaluation != nil {
		currentState.Pawn = action.Evaluation.Pawn + 50
		if action.Evaluation.IsForcedMate {
			currentState.Pawn = 100
			if action.Evaluation.ForcedMateIn < 0 {
				currentState.Pawn = 0
			}
		}
	}

	switch g.Outcome() {
	case chess.BlackWon:
		currentState.Pawn = 0
	case chess.WhiteWon:
		currentState.Pawn = 100
	case chess.Draw:
		currentState.Pawn = 50
	}

	currentState.FEN = g.Position().Board().String()

	currentState.Outcome = g.Outcome().String()

	if err := u.ws.WriteJSON(currentState); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
	u.mutex.Unlock()
}
