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
	sockets      []*websocket.Conn
	mutex        *sync.Mutex
	currentState *GameState
}

func NewWS() *WSUI {
	game := chess.NewGame()
	return &WSUI{
		mutex: &sync.Mutex{},
		currentState: &GameState{
			SVGPosition: getSVG(*game.Position().Board()),
		},
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

var moveEncoder = chess.AlgebraicNotation{}

func (u *WSUI) Render(g chess.Game, action game.UIAction) {
	u.mutex.Lock()
	if len(g.Moves()) == 0 {
		u.currentState.SVGPosition = getSVG(*g.Position().Board())
	} else {
		move := g.Moves()[len(g.Moves())-1]
		mark := image.MarkSquares(yellow, move.S1(), move.S2())

		u.currentState.SVGPosition = getSVG(*g.Position().Board(), mark)

		u.currentState.Turn = g.Position().Turn().String()
		u.currentState.PGN = g.String()
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

	u.currentState.Moves = moves

	if action.Evaluation != nil {
		u.currentState.Pawn = action.Evaluation.Pawn + 50
		if action.Evaluation.IsForcedMate {
			u.currentState.Pawn = 100
			if action.Evaluation.ForcedMateIn < 0 {
				u.currentState.Pawn = 0
			}
		}

		if len(action.Evaluation.BestMoves) > 0 {
			nextBestMove := action.Evaluation.BestMoves[0]
			move := g.Moves()[len(g.Moves())-1]

			markLast := image.MarkSquares(yellow, move.S1(), move.S2())
			markBest := image.MarkSquares(green, nextBestMove.S1(), nextBestMove.S2())

			u.currentState.SVGNextBestMove = getSVG(*g.Position().Board(), markLast, markBest)
		}
	}

	switch g.Outcome() {
	case chess.BlackWon:
		u.currentState.Pawn = 0
	case chess.WhiteWon:
		u.currentState.Pawn = 100
	case chess.Draw:
		u.currentState.Pawn = 50
	}

	u.currentState.FEN = g.Position().Board().String()

	u.currentState.Outcome = g.Outcome().String()

	for _, ws := range u.sockets {
		u.sendCurentState(ws)
	}
	u.mutex.Unlock()
}

func getSVG(board chess.Board, opts ...func(*image.Encoder)) string {
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

func (u *WSUI) sendCurentState(ws *websocket.Conn) {
	if err := ws.WriteJSON(u.currentState); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
}
