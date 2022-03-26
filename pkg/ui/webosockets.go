package ui

import (
	"bytes"
	"errors"
	"image/color"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
	"github.com/windler/chesspal/pkg/game"
)

type WSUI struct {
	sockets []*websocket.Conn
	mutex   *sync.Mutex
}

func NewWS() *WSUI {
	return &WSUI{
		mutex: &sync.Mutex{},
	}
}

func (u *WSUI) AddWebsocket(ws *websocket.Conn) {
	u.sendCurentState(ws)
	u.sockets = append(u.sockets, ws)
}

type GameState struct {
	SVGPosition     string  `json:"svgPosition"`
	SVGNextBestMove string  `json:"svgNextBestMove"`
	Pawn            float64 `json:"pawn"`
	Moves           []Move  `json:"moves"`
	Turn            string  `json:"turn"`
	PGN             string  `json:"pgn"`
	FEN             string  `json:"fen"`
	Outcome         string  `json:"outcome"`
}

type Move struct {
	Move     string `json:"move"`
	Accuracy string `json:"accuracy"`
	Color    string `json:"color"`
}

var yellow = color.RGBA{255, 255, 0, 1}
var green = color.RGBA{0, 90, 0, 1}

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

		if len(action.Evaluation.BestMoves) > 0 {
			buf := bytes.NewBufferString("")

			nextBestMove := action.Evaluation.BestMoves[0]
			move := g.Moves()[len(g.Moves())-1]

			markLast := image.MarkSquares(yellow, move.S1(), move.S2())
			markBest := image.MarkSquares(green, nextBestMove.S1(), nextBestMove.S2())

			if err := image.SVG(buf, g.Position().Board(), markLast, markBest, colors); err != nil {
				log.Printf("error occurred: %v", err)
			}
			currentState.SVGNextBestMove = buf.String()
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

	for _, ws := range u.sockets {
		u.sendCurentState(ws)
	}
	u.mutex.Unlock()
}

func (u *WSUI) sendCurentState(ws *websocket.Conn) {
	if err := ws.WriteJSON(currentState); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
}
