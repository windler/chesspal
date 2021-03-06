package ui

import (
	"errors"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
	"github.com/windler/chesspal/pkg/game"
	"github.com/windler/chesspal/pkg/util"
)

type WSUI struct {
	sockets      map[*websocket.Conn]*sync.Mutex
	mutex        *sync.Mutex
	currentState *GameState
}

func NewWS() *WSUI {
	game := chess.NewGame()
	return &WSUI{
		mutex:   &sync.Mutex{},
		sockets: make(map[*websocket.Conn]*sync.Mutex),
		currentState: &GameState{
			SVGPosition: util.GetSVG(*game.Position().Board()),
		},
	}
}

func (u *WSUI) AddWebsocket(ws *websocket.Conn) {
	u.sockets[ws] = &sync.Mutex{}
	u.sendCurentState(ws, u.sockets[ws])
}

func (u *WSUI) RemoveWebsocket(ws *websocket.Conn) {
	delete(u.sockets, ws)
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
		u.currentState.SVGPosition = util.GetSVG(*g.Position().Board())
	} else {
		move := g.Moves()[len(g.Moves())-1]
		mark := image.MarkSquares(yellow, move.S1(), move.S2())

		u.currentState.SVGPosition = util.GetSVG(*g.Position().Board(), mark)

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

			u.currentState.SVGNextBestMove = util.GetSVG(*g.Position().Board(), markLast, markBest)
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

	for ws, mutex := range u.sockets {
		go u.sendCurentState(ws, mutex)
	}
	u.mutex.Unlock()
}

func (u *WSUI) sendCurentState(ws *websocket.Conn, mutex *sync.Mutex) {
	mutex.Lock()
	ws.SetWriteDeadline(time.Now().Add(time.Second * 5))

	if err := ws.WriteJSON(u.currentState); !errors.Is(err, nil) {
		log.Printf("error occurred: %v", err)
	}
	mutex.Unlock()
}

func (u *WSUI) Reset() {
	u.currentState = &GameState{}
}

func (u *WSUI) SendBoard(board chess.Board) {
	u.currentState = &GameState{
		SVGPosition: util.GetSVG(board),
	}
	for ws, mutex := range u.sockets {
		go u.sendCurentState(ws, mutex)
	}
}
