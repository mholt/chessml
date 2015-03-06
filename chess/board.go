package chess

import (
	"bytes"
	"fmt"
)

// A Board represents a chess board.
type Board struct {
	Spaces [Size][Size]Piece
}

// Setup resets the board state, placing pieces in their initial positions.
func (b *Board) Setup() {
	// Wipe everything off.
	for c := 0; c < Size; c++ {
		for r := 0; r < Size; r++ {
			b.Spaces[r][c].Rank = Empty
			b.Spaces[r][c].Color = WhiteTeam
		}
	}

	// Place all of the Pawns.
	for c := 0; c < Size; c++ {
		b.Spaces[1][c].Rank = Pawn

		b.Spaces[6][c].Rank = Pawn
		b.Spaces[6][c].Color = BlackTeam
	}

	for c := 0; c < Size; c++ {
		b.Spaces[6][c].Color = BlackTeam
		b.Spaces[7][c].Color = BlackTeam
	}

	// Place the Kings.
	b.Spaces[0][4].Rank = King
	b.Spaces[7][4].Rank = King

	// Place the Queens.
	b.Spaces[0][3].Rank = Queen
	b.Spaces[7][3].Rank = Queen

	// Place the Bishops.
	b.Spaces[0][5].Rank = Bishop
	b.Spaces[0][2].Rank = Bishop
	b.Spaces[7][5].Rank = Bishop
	b.Spaces[7][2].Rank = Bishop

	// Place the Knights
	b.Spaces[0][1].Rank = Knight
	b.Spaces[0][6].Rank = Knight
	b.Spaces[7][1].Rank = Knight
	b.Spaces[7][6].Rank = Knight

	// Place the Rooks
	b.Spaces[0][0].Rank = Rook
	b.Spaces[0][7].Rank = Rook
	b.Spaces[7][0].Rank = Rook
	b.Spaces[7][7].Rank = Rook
}

// String creates a string representation of the current state of the board.
func (b Board) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("    A    B    C    D    E    F    G    H   \n")
	buffer.WriteString("  ╔════╤════╤════╤════╤════╤════╤════╤════╗\n")
	for r := Size - 1; r >= 0; r-- {
		buffer.WriteString(fmt.Sprintf("%d ", r+1))

		for c := 0; c < Size; c++ {
			piece := b.Spaces[r][c]
			vertical := "│"
			if c == 0 {
				vertical = "║"
			}

			if piece.Rank == Empty {
				buffer.WriteString(vertical + "    ")
			} else {
				buffer.WriteString(fmt.Sprintf("%s %s  ", vertical, PieceSymbol(piece)))
			}
		}

		buffer.WriteString(fmt.Sprintf("║ %d\n", r+1))

		if r > 0 {
			buffer.WriteString("  ╟────┼────┼────┼────┼────┼────┼────┼────╢\n")
		} else {
			buffer.WriteString("  ╚════╧════╧════╧════╧════╧════╧════╧════╝\n")
		}
	}
	buffer.WriteString("    A    B    C    D    E    F    G    H   \n")
	return buffer.String()
}

func PieceSymbol(p Piece) string {
	if p.Color == WhiteTeam {
		switch p.Rank {
		case King:
			return "♔"
		case Queen:
			return "♕"
		case Bishop:
			return "♗"
		case Knight:
			return "♘"
		case Rook:
			return "♖"
		case Pawn:
			return "♙"
		default:
			return "◯"
		}
	} else if p.Color == BlackTeam {
		switch p.Rank {
		case King:
			return "♚"
		case Queen:
			return "♛"
		case Bishop:
			return "♝"
		case Knight:
			return "♞"
		case Rook:
			return "♜"
		case Pawn:
			return "♟"
		default:
			return "⬤"
		}
	}
	return "?"
}

type (
	// Color represents white or black
	Color int

	// Rank represents a kind of piece (King, Pawn, etc...)
	Rank int

	// Piece is a piece on the board
	Piece struct {
		Color
		Rank
	}

	// Coord is a chessboard coordinate.
	Coord struct {
		Row, Col int
	}
)

const (
	WhiteTeam Color = iota
	BlackTeam
)

const (
	Empty Rank = iota
	King
	Queen
	Bishop
	Knight
	Rook
	Pawn
)

// Number of spaces in one direction
const Size = 8

var (
	RankToSymbol map[Rank]string = map[Rank]string{
		Empty:  " ",
		King:   "K",
		Queen:  "Q",
		Bishop: "B",
		Knight: "N",
		Rook:   "R",
		Pawn:   "P",
	}

	ColorToSymbol map[Color]string = map[Color]string{
		WhiteTeam: "W",
		BlackTeam: "B",
	}
)
