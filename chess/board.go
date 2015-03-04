package chess

import (
	"bytes"
	"fmt"
)

type Color int

const (
	WhiteTeam Color = iota
	BlackTeam
)

type Rank int

const (
	Empty Rank = iota
	King
	Queen
	Bishop
	Knight
	Rook
	Pawn
)

type Piece struct {
	Color
	Rank
}

const Size = 8

var RankToSymbol map[Rank]string = map[Rank]string{
	Empty:  " ",
	King:   "K",
	Queen:  "Q",
	Bishop: "B",
	Knight: "N",
	Rook:   "R",
	Pawn:   "P",
}

var ColorToSymbol map[Color]string = map[Color]string{
	WhiteTeam: "W",
	BlackTeam: "B",
}

// A Board represents a chess board.
type Board struct {
	Spaces [Size][Size]Piece
}

type Coord struct {
	Row, Col int
}

type ValidMove struct {
	From, To Coord
	Capture  bool
}

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

func (b Board) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("+----+----+----+----+----+----+----+----+\n")
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			piece := b.Spaces[r][c]

			if piece.Rank == Empty {
				buffer.WriteString("|    ")
			} else {
				buffer.WriteString(fmt.Sprintf("| %s%s ", ColorToSymbol[piece.Color], RankToSymbol[piece.Rank]))
			}
		}

		buffer.WriteString("|\n+----+----+----+----+----+----+----+----+\n")
	}

	return buffer.String()
}

func RookMove(b *Board, row, col int) (possible []ValidMove) {
	possible = extend(possible, lineMove(b, row, col, -1, 0)) // up
	possible = extend(possible, lineMove(b, row, col, 1, 0))  // down
	possible = extend(possible, lineMove(b, row, col, 0, -1)) // left
	possible = extend(possible, lineMove(b, row, col, 0, 1))  // right

	return
}

func BishopMove(b *Board, row, col int) (possible []ValidMove) {
	possible = extend(possible, lineMove(b, row, col, -1, -1)) // up-left
	possible = extend(possible, lineMove(b, row, col, -1, 1))  // up-right
	possible = extend(possible, lineMove(b, row, col, 1, -1))  // down-left
	possible = extend(possible, lineMove(b, row, col, 1, 1))   // down-right

	return
}

func QueenMove(b *Board, row, col int) (possible []ValidMove) {
	possible = extend(possible, lineMove(b, row, col, -1, 0))  // up
	possible = extend(possible, lineMove(b, row, col, 1, 0))   // down
	possible = extend(possible, lineMove(b, row, col, 0, -1))  // left
	possible = extend(possible, lineMove(b, row, col, 0, 1))   // right
	possible = extend(possible, lineMove(b, row, col, -1, -1)) // up-left
	possible = extend(possible, lineMove(b, row, col, -1, 1))  // up-right
	possible = extend(possible, lineMove(b, row, col, 1, -1))  // down-left
	possible = extend(possible, lineMove(b, row, col, 1, 1))   // down-right

	return
}

func extend(vm, toAdd []ValidMove) []ValidMove {
	for _, ta := range toAdd {
		vm = append(vm, ta)
	}

	return vm
}

func KnightMove(b *Board, row, col int) (possible []ValidMove) {
	possible = tryAndAppend(possible, b, row, col, -2, -1) // up-left
	possible = tryAndAppend(possible, b, row, col, -2, 1)  // up-right
	possible = tryAndAppend(possible, b, row, col, 2, -1)  // down-left
	possible = tryAndAppend(possible, b, row, col, 2, 1)   // down-right

	possible = tryAndAppend(possible, b, row, col, -1, -2) // left-up
	possible = tryAndAppend(possible, b, row, col, 1, -2)  // left-down
	possible = tryAndAppend(possible, b, row, col, -1, 2)  // right-up
	possible = tryAndAppend(possible, b, row, col, 1, 2)   // right-down

	return
}

func KingMove(b *Board, row, col int) (possible []ValidMove) {
	possible = tryAndAppend(possible, b, row, col, -1, 0) // up
	possible = tryAndAppend(possible, b, row, col, 1, 0)  // down
	possible = tryAndAppend(possible, b, row, col, 0, -1) // left
	possible = tryAndAppend(possible, b, row, col, 0, 1)  // right

	possible = tryAndAppend(possible, b, row, col, -1, -1) // up-left
	possible = tryAndAppend(possible, b, row, col, -1, 1)  // up-right
	possible = tryAndAppend(possible, b, row, col, 1, -1)  // down-left
	possible = tryAndAppend(possible, b, row, col, 1, 1)   // down-right

	return
}

func tryAndAppend(vm []ValidMove, b *Board, row, col, rowDiff, colDiff int) []ValidMove {
	color := b.Spaces[row][col].Color

	valid, capture := tryMove(b, color, row+rowDiff, col+colDiff) // down-right
	if valid {
		return append(vm, ValidMove{To: Coord{row, col}, From: Coord{row + 1, col + 1}, Capture: capture})
	}

	return vm
}

func PawnMove(b *Board, row, col int) (possible []ValidMove) {
	color := b.Spaces[row][col].Color

	if color == WhiteTeam {
		// move down (+)
		valid, _ := tryMove(b, color, row+1, col)

		if valid {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row + 1, col}, Capture: false})
		}

		valid, capture := tryMove(b, color, row+1, col-1) // left capture

		if valid && capture {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row + 1, col - 1}, Capture: true})
		}

		valid, capture = tryMove(b, color, row+1, col+1) // right capture

		if valid && capture {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row + 1, col + 1}, Capture: true})
		}

		if row == 1 { // double move from starting row
			valid, capture = tryMove(b, color, row+2, col)

			if valid && !capture {
				possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row + 2, col}, Capture: true})
			}
		}
	} else {
		// move up (-)
		valid, _ := tryMove(b, color, row-1, col)

		if valid {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row - 1, col}, Capture: false})
		}

		valid, capture := tryMove(b, color, row-1, col-1) // left capture

		if valid && capture {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row - 1, col - 1}, Capture: true})
		}

		valid, capture = tryMove(b, color, row-1, col+1) // right capture

		if valid && capture {
			possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row - 1, col + 1}, Capture: true})
		}

		if row == 6 { // double move from starting row
			valid, capture = tryMove(b, color, row-2, col)

			if valid && !capture {
				possible = append(possible, ValidMove{To: Coord{row, col}, From: Coord{row - 2, col}, Capture: true})
			}
		}
	}

	return
}

func tryMove(b *Board, pieceColor Color, row, col int) (valid, capture bool) {
	if row < 0 || row >= Size || col < 0 || col >= Size {
		return false, false
	}

	target := b.Spaces[row][col]

	if target.Rank == Empty {
		// Valid move to empty square.
		return true, false
	} else if target.Color != pieceColor {
		// Enemy Piece captured.
		return true, true
	}

	return false, false
}

func lineMove(b *Board, row, col, rowDiff, colDiff int) (possible []ValidMove) {
	color := b.Spaces[row][col].Color
	toRow, toCol := row, col

	for {
		toRow += rowDiff
		toCol += colDiff

		valid, capture := tryMove(b, color, toRow, toCol)

		if !valid {
			break
		} else if capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{toRow, toCol}, Capture: capture})
			break
		}

		possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{toRow, toCol}, Capture: capture})
	}

	return
}
