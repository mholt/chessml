package chess

import "fmt"

// PossibleMoves returns the possible moves of piece p from
// the position row,col.
func PossibleMoves(b Board, p Piece, row, col int) []ValidMove {
	switch p.Rank {
	case King:
		return KingMoves(b, row, col)
	case Queen:
		return QueenMoves(b, row, col)
	case Bishop:
		return BishopMoves(b, row, col)
	case Knight:
		return KnightMoves(b, row, col)
	case Rook:
		return RookMoves(b, row, col)
	case Pawn:
		return PawnMoves(b, row, col)
	default:
		panic(fmt.Sprintf("Invalid piece: bad Rank value %d", p.Rank))
	}
}

// RookMoves computes possible moves for a rook on board b at the row and col position.
func RookMoves(b Board, row, col int) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, 0)...) // up
	possible = append(possible, lineMove(b, row, col, 1, 0)...)  // down
	possible = append(possible, lineMove(b, row, col, 0, -1)...) // left
	possible = append(possible, lineMove(b, row, col, 0, 1)...)  // right

	return
}

// BishopMoves computes possible moves for a bishop on board b at the row and col position.
func BishopMoves(b Board, row, col int) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, -1)...) // up-left
	possible = append(possible, lineMove(b, row, col, -1, 1)...)  // up-right
	possible = append(possible, lineMove(b, row, col, 1, -1)...)  // down-left
	possible = append(possible, lineMove(b, row, col, 1, 1)...)   // down-right

	return
}

// QueenMoves computes possible moves for a queen on board b at the row and col position.
func QueenMoves(b Board, row, col int) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, 0)...)  // up
	possible = append(possible, lineMove(b, row, col, 1, 0)...)   // down
	possible = append(possible, lineMove(b, row, col, 0, -1)...)  // left
	possible = append(possible, lineMove(b, row, col, 0, 1)...)   // right
	possible = append(possible, lineMove(b, row, col, -1, -1)...) // up-left
	possible = append(possible, lineMove(b, row, col, -1, 1)...)  // up-right
	possible = append(possible, lineMove(b, row, col, 1, -1)...)  // down-left
	possible = append(possible, lineMove(b, row, col, 1, 1)...)   // down-right

	return
}

// KnightMoves computes possible moves for a knight on board b at the row and col position.
func KnightMoves(b Board, row, col int) (possible []ValidMove) {
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

// KingMoves computes possible moves for a king on board b at the row and col position.
func KingMoves(b Board, row, col int) (possible []ValidMove) {
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

// PawnMoves computes possible moves for a pawn on board b at the row and col position.
func PawnMoves(b Board, row, col int) (possible []ValidMove) {
	color := b.Spaces[row][col].Color

	if color == WhiteTeam {
		// move up (+)
		valid, capture := tryMove(b, color, row+1, col)
		if valid && !capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col}, Capture: false})
		}

		valid, capture = tryMove(b, color, row+1, col-1) // left capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col - 1}, Capture: true})
		}

		valid, capture = tryMove(b, color, row+1, col+1) // right capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col + 1}, Capture: true})
		}

		if row == 1 { // double move from starting row
			valid, capture = tryMove(b, color, row+2, col)
			if valid && !capture && b.Spaces[row+1][col].Rank == Empty {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 2, col}, Capture: true})
			}
		}

		if col+1 < Size && b.Spaces[row][col+1].Color != color && b.Spaces[row][col+1].EnPassantable {
			valid, capture = tryMove(b, color, row+1, col+1)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col + 1}, EnPassant: true})
			}
		}

		if col-1 >= 0 && b.Spaces[row][col-1].Color != color && b.Spaces[row][col-1].EnPassantable {
			valid, capture = tryMove(b, color, row+1, col-1)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col - 1}, EnPassant: true})
			}
		}

	} else {
		// move down (-)
		valid, capture := tryMove(b, color, row-1, col)
		if valid && !capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col}, Capture: false})
		}

		valid, capture = tryMove(b, color, row-1, col-1) // left capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col - 1}, Capture: true})
		}

		valid, capture = tryMove(b, color, row-1, col+1) // right capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col + 1}, Capture: true})
		}

		if row == 6 { // double move from starting row
			valid, capture = tryMove(b, color, row-2, col)
			if valid && !capture && b.Spaces[row-1][col].Rank == Empty {
				possible = append(possible, ValidMove{
					From:    Coord{row, col},
					To:      Coord{row - 2, col},
					Capture: true,
				})
			}
		}

		if col+1 < Size && b.Spaces[row][col+1].Color != color && b.Spaces[row][col+1].EnPassantable {
			valid, capture = tryMove(b, color, row-1, col+1)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col + 1}, EnPassant: true})
			}
		}

		if col-1 >= 0 && b.Spaces[row][col-1].Color != color && b.Spaces[row][col-1].EnPassantable {
			valid, capture = tryMove(b, color, row-1, col-1)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col - 1}, EnPassant: true})
			}
		}

	}

	return
}

func tryMove(b Board, pieceColor Color, toRow, toCol int) (valid, capture bool) {
	if toRow < 0 || toRow >= Size || toCol < 0 || toCol >= Size {
		return false, false
	}

	target := b.Spaces[toRow][toCol]

	if target.Rank == Empty {
		return true, false // Valid move to empty square.
	} else if target.Color != pieceColor {
		return true, true // Valid move, and enemy piece would be captured.
	}

	return false, false
}

func lineMove(b Board, row, col, rowDiff, colDiff int) (possible []ValidMove) {
	color := b.Spaces[row][col].Color
	toRow, toCol := row, col

	for {
		toRow += rowDiff
		toCol += colDiff

		valid, capture := tryMove(b, color, toRow, toCol)

		if !valid {
			break
		} else if capture {
			possible = append(possible, ValidMove{
				From:    Coord{row, col},
				To:      Coord{toRow, toCol},
				Capture: capture,
			})
			break
		}

		possible = append(possible, ValidMove{
			From:    Coord{row, col},
			To:      Coord{toRow, toCol},
			Capture: capture,
		})
	}

	return
}

func tryAndAppend(vm []ValidMove, b Board, row, col, rowDiff, colDiff int) []ValidMove {
	color := b.Spaces[row][col].Color

	valid, capture := tryMove(b, color, row+rowDiff, col+colDiff)
	if valid {
		return append(vm, ValidMove{
			From:    Coord{row, col},
			To:      Coord{row + rowDiff, col + colDiff},
			Capture: capture,
		})
	}

	return vm
}

// movePossible returns whether piece can move from row,col to destRow,destCol.
func movePossible(b Board, piece Piece, row, col, destRow, destCol int) bool {
	possible := PossibleMoves(b, piece, row, col)
	for _, move := range possible {
		if move.To.Row == destRow && move.To.Col == destCol {
			return true
		}
	}
	return false
}

// isCheck returns true if c's King is in check, false otherwise.
func isCheck(b Board, c Color) bool {
	var kingPos Coord

	// Find c's king
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if b.Spaces[i][j].Color == c && b.Spaces[i][j].Rank == King {
				kingPos = Coord{Row: i, Col: j}
			}
		}
	}

	// See if anything of the other color can move to its spot
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			if b.Spaces[i][j].Color != c && b.Spaces[i][j].Rank != Empty {
				if movePossible(b, b.Spaces[i][j], i, j, kingPos.Row, kingPos.Col) {
					return true
				}
			}
		}
	}

	return false
}

// ValidMove represents a possible move that has not necessarily been made.
type ValidMove struct {
	From, To  Coord
	Capture   bool
	EnPassant bool
}
