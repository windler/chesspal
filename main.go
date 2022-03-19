package main

import (
	"fmt"

	"github.com/jacobsa/go-serial/serial"
	"github.com/notnil/chess"
)

const DGT_SEND_BRD = 0x42
const DGT_SEND_RESET = 0x31
const DGT_SEND_UPDATE_BRD = 0x44

var game = chess.NewGame()

func main() {
	// board := godgt.NewDgtBoard("/dev/ttyACM0")
	// defer board.Close()
	// go board.ReadLoop()
	// board.WriteCommand(godgt.DGT_SEND_RESET)
	// board.WriteCommand(godgt.DGT_SEND_BRD)
	// board.WriteCommand(godgt.DGT_SEND_UPDATE_BRD)

	// for msg := range board.MessagesFromBoard {
	// 	if msg.FieldUpdate != nil {
	// 		fmt.Println(msg.FieldUpdate.ToString())
	// 	}
	// 	if msg.BoardUpdate != nil {
	// 		fmt.Println(msg.BoardUpdate.ToString())
	// 	}
	// 	if msg.InfoUpdate != nil {
	// 		fmt.Println(msg.InfoUpdate.ToString())
	// 	}
	// }

	options := serial.OpenOptions{
		PortName:        "/dev/ttyACM0",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	io, err := serial.Open(options)
	if err != nil {
		fmt.Println(err)
	}
	defer io.Close()

	io.Write([]byte{DGT_SEND_RESET})
	io.Write([]byte{DGT_SEND_BRD})
	io.Write([]byte{DGT_SEND_UPDATE_BRD})

	// var s1, s2 *chess.Square
	for {
		// log.Println("About to read bytes")
		buf := make([]byte, 1024)
		n, err := io.Read(buf)
		if err != nil {
			fmt.Println("error reading bytes.")
		}
		if n > 0 {
			msgType := getMessageType(buf[0:n])
			fmt.Println(msgType)
			if msgType == COMMAND_FIELD_UPDATE {
				io.Write([]byte{DGT_SEND_BRD})
			} else if msgType == COMMAND_BOARD_DUMP {
				pieces := getChessBoard(buf[0:n])
				moves := game.Position().ValidMoves()
				for _, move := range moves {
					pos := *game.Clone().Position()
					sMap := pos.Update(move).Board().SquareMap()

					valid := true
					piecesFound := 0
					fmt.Println(pieces)
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
						game.Move(move)
						fmt.Println(game.Position().Board().Draw())
						fmt.Println(move)
						break
					}

				}
			}

			// if p.Color() != game.Position().Turn() {
			// 	continue
			// }

			// if s1 == nil {
			// 	s1 = &s
			// 	continue
			// }
			// if s2 == nil {
			// 	s2 = &s
			// 	for _, move := range game.Position().ValidMoves() {
			// 		if *s2 == move.S2() && *s1 == move.S1() {
			// 			game.Move(move)
			// 			break
			// 		}

			// 	}
			// 	fmt.Println(game.Position().Board().Draw())
			// 	s1 = nil
			// 	s2 = nil
			// }
		}

	}

}

func getChessBoard(msg []byte) []PieceOnSqaure {
	result := []PieceOnSqaure{}
	if len(msg) != COMMAND_BOARD_DUMP_SIZE || msg[0] != COMMAND_BOARD_DUMP {
		fmt.Printf("code: %d length %d\n", msg[0], len(msg))
		return result
	}

	for i := 0; i < 64; i++ {
		result = append(result, PieceOnSqaure{
			Piece:  getPiece(int(msg[3+i])),
			Sqaure: getField(i),
		})
	}

	return result
}

func getMessageType(msg []byte) int {
	if len(msg) == 0 {
		return UNKNOWN
	}

	switch msg[0] {
	case COMMAND_FIELD_UPDATE:
		return COMMAND_FIELD_UPDATE
	case COMMAND_BOARD_DUMP:
		return COMMAND_BOARD_DUMP
	}

	return UNKNOWN
}

func decryptMessage(msg []byte) (chess.Square, chess.Piece) {
	if len(msg) == 0 {
		return 0, 0
	}

	switch msg[0] {
	case COMMAND_FIELD_UPDATE:
		if len(msg) != COMMAND_FIELD_UPDATE_SIZE {
			return 0, 0
		}
		field := getField(int(msg[3]))
		piece := getPiece(int(msg[4]))

		return field, piece

	}
	return 0, 0
}

type PieceOnSqaure struct {
	Piece  chess.Piece
	Sqaure chess.Square
}

func getField(i int) chess.Square {
	fileIndex := int(i & 0x07)
	rankIndex := 7 - int((i&0x38)>>3)

	return chess.Square((int(rankIndex) * 8) + int(fileIndex))
}

func getPiece(i int) chess.Piece {
	if i == EMPTY {
		return chess.NoPiece
	}

	if i > QUEEN {
		i = i - QUEEN
		switch i {
		case PAWN:
			return chess.BlackPawn
		case BISHOP:
			return chess.BlackBishop
		case KNIGHT:
			return chess.BlackKnight
		case KING:
			return chess.BlackKing
		case QUEEN:
			return chess.BlackQueen
		case ROOK:
			return chess.BlackRook
		}
	}

	switch i {
	case PAWN:
		return chess.WhitePawn
	case BISHOP:
		return chess.WhiteBishop
	case KNIGHT:
		return chess.WhiteKnight
	case KING:
		return chess.WhiteKing
	case QUEEN:
		return chess.WhiteQueen
	case ROOK:
		return chess.WhiteRook
	}

	return chess.NoPiece
}

const UNKNOWN = -1

const MESSAGE_BIT = 0x80
const COMMAND_FIELD_UPDATE = (MESSAGE_BIT | 0x0e)
const COMMAND_FIELD_UPDATE_SIZE = 5

const DGT_BOARD_DUMP = 0x06
const COMMAND_BOARD_DUMP = (MESSAGE_BIT | DGT_BOARD_DUMP)
const COMMAND_BOARD_DUMP_SIZE = 67

const EMPTY = 0x00
const PAWN = 0x01
const ROOK = 0x02
const KNIGHT = 0x03
const BISHOP = 0x04
const KING = 0x05
const QUEEN = 0x06
