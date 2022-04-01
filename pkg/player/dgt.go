package player

import (
	"io"
	"log"
	"math"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/notnil/chess"
)

const DGT_SEND_BRD = 0x42
const DGT_SEND_RESET = 0x40
const DGT_SEND_UPDATE_BRD = 0x44

const MESSAGE_BIT = 0x80

const DGT_MSG_TYPE_UNKNOWN = -1

const DGT_MSG_TYPE_FIELD_UPDATE = (MESSAGE_BIT | 0x0e)
const DGT_MSG_TYPE_FIELD_UPDATE_SIZE = 5

const DGT_MSG_TYPE_BOARD_DUMP = (MESSAGE_BIT | 0x06)
const DGT_MSG_TYPE_BOARD_DUMP_SIZE = 67

const DGT_EMPTY = 0x00
const DGT_PAWN = 0x01
const DGT_ROOK = 0x02
const DGT_KNIGHT = 0x03
const DGT_BISHOP = 0x04
const DGT_KING = 0x05
const DGT_QUEEN = 0x06

type DGT struct {
	engine *DGTEngine
	name   string
}

type DGTEngine struct {
	io             io.ReadWriteCloser
	wg             *sync.WaitGroup
	colors         []chess.Color
	game           *chess.Game
	upsideDown     bool
	lastUpdateTime time.Time
	positionChan   chan chess.Board
}

func NewDGTPlayer(name string, engine *DGTEngine) *DGT {
	log.Printf("DGT player created ")
	return &DGT{
		engine: engine,
		name:   name,
	}
}

func (p *DGT) Name() string {
	return p.name
}

func NewDGTEngine() *DGTEngine {
	return &DGTEngine{
		wg:           &sync.WaitGroup{},
		colors:       []chess.Color{},
		positionChan: make(chan chess.Board),
	}
}

func (p *DGTEngine) PostionChannel() chan chess.Board {
	return p.positionChan
}

func (p *DGTEngine) AddColor(color chess.Color) {
	p.colors = append(p.colors, color)
}
func (p *DGTEngine) SetUpsideDown(ud bool) {
	p.upsideDown = ud
}

func (p *DGTEngine) MakeMove(game *chess.Game) {
	p.game = game
	p.wg.Add(1)
	p.wg.Wait()
}

func (p *DGTEngine) Start(port string) {
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	io, err := serial.Open(options)
	if err != nil {
		panic(err)
	}

	io.Write([]byte{DGT_SEND_RESET})
	io.Write([]byte{DGT_SEND_BRD})
	io.Write([]byte{DGT_SEND_UPDATE_BRD})

	p.io = io

	go p.readLoop()
}

func (p *DGTEngine) ReadCurrentPosition() {
	p.io.Write([]byte{DGT_SEND_BRD})
}

const delayMs = 200

func (p *DGTEngine) readLoop() {
	for {
		buf := make([]byte, 1024)
		n, err := p.io.Read(buf)
		if err != nil {
			log.Printf("error reading bytes from serial port: %s\n", err)
		}
		p.lastUpdateTime = time.Now()

		// delay move a bit to allow "sliding" of pieces on board
		time.AfterFunc(delayMs*time.Millisecond, func() {
			if p.lastUpdateTime.Add(delayMs * time.Millisecond).After(time.Now()) {
				return
			}

			msgType := getMessageType(buf[0:n])
			if msgType == DGT_MSG_TYPE_FIELD_UPDATE {
				p.io.Write([]byte{DGT_SEND_BRD})
			}

			if msgType == DGT_MSG_TYPE_BOARD_DUMP {
				pieces := p.getChessBoard(buf[0:n])
				p.positionChan <- getBoard(pieces)

				if p.game != nil && n > 0 && p.game.Outcome() == chess.NoOutcome {
					moves := p.game.Position().ValidMoves()
					for _, move := range moves {
						if p.isValidMoveInPostion(move, *p.game.Clone().Position(), pieces) {
							p.game.Move(move)
							p.wg.Done()
							break
						}
					}
				}
			}

		})
	}
}

func getBoard(pieces []PieceOnSqaure) chess.Board {
	boardSquares := map[chess.Square]chess.Piece{}

	for _, poc := range pieces {
		boardSquares[poc.Sqaure] = poc.Piece
	}

	return *chess.NewBoard(boardSquares)
}

func (p *DGTEngine) isValidMoveInPostion(move *chess.Move, pos chess.Position, pieces []PieceOnSqaure) bool {
	sMap := pos.Update(move).Board().SquareMap()

	valid := true
	piecesFound := 0

	for _, p := range pieces {
		if p.Piece.Type() == chess.NoPieceType {
			continue
		}
		if sMap[p.Sqaure] != p.Piece {
			valid = false
			break
		}
		piecesFound = piecesFound + 1
	}

	if valid && piecesFound == len(sMap) {
		return true
	}
	return false
}

func (p *DGT) SetColor(color chess.Color) {
	p.engine.AddColor(color)
}

func (p *DGT) MakeMove(game *chess.Game) {
	p.engine.MakeMove(game)
}

func (p *DGT) End() {

}

func (p *DGTEngine) getChessBoard(msg []byte) []PieceOnSqaure {
	result := []PieceOnSqaure{}
	if len(msg) != DGT_MSG_TYPE_BOARD_DUMP_SIZE || msg[0] != DGT_MSG_TYPE_BOARD_DUMP {
		return result
	}

	for i := 0; i < 64; i++ {
		result = append(result, PieceOnSqaure{
			Piece:  p.getPiece(int(msg[3+i])),
			Sqaure: p.getSquare(i),
		})
	}

	return result
}

func getMessageType(msg []byte) int {
	if len(msg) == 0 {
		return DGT_MSG_TYPE_UNKNOWN
	}

	return int(msg[0])
}

type PieceOnSqaure struct {
	Piece  chess.Piece
	Sqaure chess.Square
}

func (p *DGTEngine) getSquare(i int) chess.Square {
	fileIndex := int(i & 0x07)
	rankIndex := 7 - int((i&0x38)>>3)

	if p.upsideDown {
		rankIndex = int(math.Abs(float64(rankIndex) - 7))
		fileIndex = int(math.Abs(float64(fileIndex) - 7))
	}

	return chess.Square((int(rankIndex) * 8) + int(fileIndex))
}

func (p *DGTEngine) getPiece(i int) chess.Piece {
	if i == DGT_EMPTY {
		return chess.NoPiece
	}

	if i > DGT_QUEEN {
		i = i - DGT_QUEEN
		switch i {
		case DGT_PAWN:
			return chess.BlackPawn
		case DGT_BISHOP:
			return chess.BlackBishop
		case DGT_KNIGHT:
			return chess.BlackKnight
		case DGT_KING:
			return chess.BlackKing
		case DGT_QUEEN:
			return chess.BlackQueen
		case DGT_ROOK:
			return chess.BlackRook
		}
	}

	switch i {
	case DGT_PAWN:
		return chess.WhitePawn
	case DGT_BISHOP:
		return chess.WhiteBishop
	case DGT_KNIGHT:
		return chess.WhiteKnight
	case DGT_KING:
		return chess.WhiteKing
	case DGT_QUEEN:
		return chess.WhiteQueen
	case DGT_ROOK:
		return chess.WhiteRook
	}

	return chess.NoPiece
}
