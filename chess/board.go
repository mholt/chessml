package chess

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
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

	placePiece := func(row, col int, piece Rank, team Color) {
		b.Spaces[row][col].Rank = piece
		b.Spaces[row][col].Color = team
	}

	// Place all of the Pawns.
	for c := 0; c < Size; c++ {
		placePiece(1, c, Pawn, WhiteTeam)
		placePiece(6, c, Pawn, BlackTeam)
	}

	// Place the Kings.
	placePiece(0, 4, King, WhiteTeam)
	placePiece(7, 4, King, BlackTeam)

	// Place the Queens.
	placePiece(0, 3, Queen, WhiteTeam)
	placePiece(7, 3, Queen, BlackTeam)

	// Place the Bishops.
	placePiece(0, 5, Bishop, WhiteTeam)
	placePiece(0, 2, Bishop, WhiteTeam)
	placePiece(7, 5, Bishop, BlackTeam)
	placePiece(7, 2, Bishop, BlackTeam)

	// Place the Knights
	placePiece(0, 1, Knight, WhiteTeam)
	placePiece(0, 6, Knight, WhiteTeam)
	placePiece(7, 1, Knight, BlackTeam)
	placePiece(7, 6, Knight, BlackTeam)

	// Place the Rooks
	placePiece(0, 0, Rook, WhiteTeam)
	placePiece(0, 7, Rook, WhiteTeam)
	placePiece(7, 0, Rook, BlackTeam)
	placePiece(7, 7, Rook, BlackTeam)
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

// MovePiece moves a piece. It performs captures by replacing
// any existing piece on the to coordinate.
func (b *Board) MovePiece(from, to Coord) error {
	if from.Col < 0 || from.Col >= Size || to.Col < 0 || to.Col >= Size {
		return errors.New("Coordinate out of bounds")
	}

	if b.Spaces[from.Row][from.Col].Rank == Empty {
		return fmt.Errorf("No piece to move at row,col (%d,%d)", from.Row, from.Col)
	}

	b.Spaces[to.Row][to.Col] = b.Spaces[from.Row][from.Col]
	b.Spaces[from.Row][from.Col].Rank = Empty

	// Reset the en passant flags for all pieces (pawns) of this color
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if b.Spaces[i][j].Color == b.Spaces[to.Row][to.Col].Color {
				b.Spaces[i][j].EnPassantable = false
			}
		}
	}

	// If this piece is a pawn, see if the opponent could use
	// en passant on their next turn and set the flag.
	if b.Spaces[to.Row][to.Col].Rank == Pawn &&
		(to.Row-from.Row == 2 || to.Row-from.Row == -2) {
		b.Spaces[to.Row][to.Col].EnPassantable = true
	}

	return nil
}

// PieceSymbol returns the unicode chess symbol for p.
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

// NotationToCoord takes a two-character algebraic notation
// like "E4" and converts it to a coordinate.
func NotationToCoord(algebra string) Coord {
	if len(algebra) != 2 {
		panic("Algebraic notation must be 2 characters precisely; got: '" + algebra + "'")
	}
	algebra = strings.ToUpper(algebra)

	var c Coord
	file := algebra[0]
	rank := algebra[1]

	// Remember, these are ASCII code points, not numbers
	if file < 65 || file > 72 || rank < 48 || rank > 57 {
		panic("Bad position (" + algebra + ")")
	}

	c.Row = int(rank - 48 - 1)
	c.Col = int(file - 65)

	return c
}

type (
	// Color represents white or black
	Color int

	// Rank represents a kind of piece (King, Pawn, etc...)
	Rank int

	// Piece is a piece on the board
	Piece struct {
		EnPassantable bool
		Color
		Rank
	}

	// Coord is a chessboard coordinate.
	Coord struct {
		Row, Col int
	}
)

// Player colors
const (
	NoColor Color = iota
	WhiteTeam
	BlackTeam
)

// Values that represent different kinds of pieces
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
	// Map of rank (kind of piece) to algabraic symbol
	RankToSymbol = map[Rank]string{
		Empty:  " ",
		King:   "K",
		Queen:  "Q",
		Bishop: "B",
		Knight: "N",
		Rook:   "R",
		Pawn:   "",
	}

	// The reverse of RankToSymbol
	SymbolToRank = map[string]Rank{
		" ": Empty,
		"K": King,
		"Q": Queen,
		"B": Bishop,
		"N": Knight,
		"R": Rook,
		"":  Pawn,
	}

	// Map of color to its string representation
	ColorToSymbol = map[Color]string{
		WhiteTeam: "W",
		BlackTeam: "B",
	}

	// Converts rank strings (1-8) to row index (0-7)
	rankToRow = map[string]int{"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7}

	// Converts file strings (a-f) to col index (0-7) - case-sensitive!
	fileToCol = map[string]int{"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7}
)
