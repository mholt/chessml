package chess

import "fmt"

// PossibleMoves returns the possible moves of piece p from
// the position row,col.
func PossibleMoves(b Board, p Piece, row, col int, recurse bool) []ValidMove {
	switch p.Rank {
	case King:
		return KingMoves(b, row, col, recurse)
	case Queen:
		return QueenMoves(b, row, col, recurse)
	case Bishop:
		return BishopMoves(b, row, col, recurse)
	case Knight:
		return KnightMoves(b, row, col, recurse)
	case Rook:
		return RookMoves(b, row, col, recurse)
	case Pawn:
		return PawnMoves(b, row, col, recurse)
	default:
		panic(fmt.Sprintf("Invalid piece: bad Rank value %d", p.Rank))
	}
}

// RookMoves computes possible moves for a rook on board b at the row and col position.
func RookMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, 0, recurse)...) // up
	possible = append(possible, lineMove(b, row, col, 1, 0, recurse)...)  // down
	possible = append(possible, lineMove(b, row, col, 0, -1, recurse)...) // left
	possible = append(possible, lineMove(b, row, col, 0, 1, recurse)...)  // right

	return
}

// BishopMoves computes possible moves for a bishop on board b at the row and col position.
func BishopMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, -1, recurse)...) // up-left
	possible = append(possible, lineMove(b, row, col, -1, 1, recurse)...)  // up-right
	possible = append(possible, lineMove(b, row, col, 1, -1, recurse)...)  // down-left
	possible = append(possible, lineMove(b, row, col, 1, 1, recurse)...)   // down-right

	return
}

// QueenMoves computes possible moves for a queen on board b at the row and col position.
func QueenMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	possible = append(possible, lineMove(b, row, col, -1, 0, recurse)...)  // up
	possible = append(possible, lineMove(b, row, col, 1, 0, recurse)...)   // down
	possible = append(possible, lineMove(b, row, col, 0, -1, recurse)...)  // left
	possible = append(possible, lineMove(b, row, col, 0, 1, recurse)...)   // right
	possible = append(possible, lineMove(b, row, col, -1, -1, recurse)...) // up-left
	possible = append(possible, lineMove(b, row, col, -1, 1, recurse)...)  // up-right
	possible = append(possible, lineMove(b, row, col, 1, -1, recurse)...)  // down-left
	possible = append(possible, lineMove(b, row, col, 1, 1, recurse)...)   // down-right

	return
}

// KnightMoves computes possible moves for a knight on board b at the row and col position.
func KnightMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	possible = tryAndAppend(possible, b, row, col, -2, -1, recurse) // up-left
	possible = tryAndAppend(possible, b, row, col, -2, 1, recurse)  // up-right
	possible = tryAndAppend(possible, b, row, col, 2, -1, recurse)  // down-left
	possible = tryAndAppend(possible, b, row, col, 2, 1, recurse)   // down-right

	possible = tryAndAppend(possible, b, row, col, -1, -2, recurse) // left-up
	possible = tryAndAppend(possible, b, row, col, 1, -2, recurse)  // left-down
	possible = tryAndAppend(possible, b, row, col, -1, 2, recurse)  // right-up
	possible = tryAndAppend(possible, b, row, col, 1, 2, recurse)   // right-down

	return
}

// KingMoves computes possible moves for a king on board b at the row and col position.
func KingMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	possible = tryAndAppend(possible, b, row, col, -1, 0, recurse) // up
	possible = tryAndAppend(possible, b, row, col, 1, 0, recurse)  // down
	possible = tryAndAppend(possible, b, row, col, 0, -1, recurse) // left
	possible = tryAndAppend(possible, b, row, col, 0, 1, recurse)  // right

	possible = tryAndAppend(possible, b, row, col, -1, -1, recurse) // up-left
	possible = tryAndAppend(possible, b, row, col, -1, 1, recurse)  // up-right
	possible = tryAndAppend(possible, b, row, col, 1, -1, recurse)  // down-left
	possible = tryAndAppend(possible, b, row, col, 1, 1, recurse)   // down-right

	return
}

// PawnMoves computes possible moves for a pawn on board b at the row and col position.
func PawnMoves(b Board, row, col int, recurse bool) (possible []ValidMove) {
	color := b.Spaces[row][col].Color

	if color == WhiteTeam {
		// move up (+)
		valid, capture, check := tryMove(b, color, row, col, row+1, col, recurse)
		if valid && !capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col}, Capture: false, Check: check})
		}

		valid, capture, check = tryMove(b, color, row, col, row+1, col-1, recurse) // left capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col - 1}, Capture: true, Check: check})
		}

		valid, capture, check = tryMove(b, color, row, col, row+1, col+1, recurse) // right capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col + 1}, Capture: true, Check: check})
		}

		if row == 1 { // double move from starting row
			valid, capture, check = tryMove(b, color, row, col, row+2, col, recurse)
			if valid && !capture && b.Spaces[row+1][col].Rank == Empty {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 2, col}, Capture: true, Check: check})
			}
		}

		if col+1 < Size && b.Spaces[row][col+1].Color != color && b.Spaces[row][col+1].EnPassantable {
			valid, capture, check = tryMove(b, color, row, col, row+1, col+1, recurse)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col + 1}, EnPassant: true, Check: check})
			}
		}

		if col-1 >= 0 && b.Spaces[row][col-1].Color != color && b.Spaces[row][col-1].EnPassantable {
			valid, capture, check = tryMove(b, color, row, col, row+1, col-1, recurse)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row + 1, col - 1}, EnPassant: true, Check: check})
			}
		}

	} else {
		// move down (-)
		valid, capture, check := tryMove(b, color, row, col, row-1, col, recurse)
		if valid && !capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col}, Capture: false, Check: check})
		}

		valid, capture, check = tryMove(b, color, row, col, row-1, col-1, recurse) // left capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col - 1}, Capture: true, Check: check})
		}

		valid, capture, check = tryMove(b, color, row, col, row-1, col+1, recurse) // right capture
		if valid && capture {
			possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col + 1}, Capture: true, Check: check})
		}

		if row == 6 { // double move from starting row
			valid, capture, check = tryMove(b, color, row, col, row-2, col, recurse)
			if valid && !capture && b.Spaces[row-1][col].Rank == Empty {
				possible = append(possible, ValidMove{
					From:    Coord{row, col},
					To:      Coord{row - 2, col},
					Capture: true,
					Check:   check,
				})
			}
		}

		if col+1 < Size && b.Spaces[row][col+1].Color != color && b.Spaces[row][col+1].EnPassantable {
			valid, capture, check = tryMove(b, color, row, col, row-1, col+1, recurse)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col + 1}, EnPassant: true, Check: check})
			}
		}

		if col-1 >= 0 && b.Spaces[row][col-1].Color != color && b.Spaces[row][col-1].EnPassantable {
			valid, capture, check = tryMove(b, color, row, col, row-1, col-1, recurse)
			if valid && !capture {
				possible = append(possible, ValidMove{From: Coord{row, col}, To: Coord{row - 1, col - 1}, EnPassant: true, Check: check})
			}
		}

	}

	return
}

func tryMove(b Board, pieceColor Color, fromRow, fromCol, toRow, toCol int, recurse bool) (valid, capture, check bool) {
	if toRow < 0 || toRow >= Size || toCol < 0 || toCol >= Size {
		return false, false, false
	}

	target := b.Spaces[toRow][toCol]

	canPutInCheck := false
	if recurse {
		possible := PossibleMoves(b, b.Spaces[fromRow][fromCol], fromRow, fromCol, false)
		for _, move := range possible {
			if move.Capture && b.Spaces[move.To.Row][move.To.Col].Rank == King {
				canPutInCheck = true
				break
			}
		}
	}

	if target.Rank == Empty {
		return true, false, canPutInCheck // Valid move to empty square.
	} else if target.Color != pieceColor {
		return true, true, canPutInCheck // Valid move, and enemy piece would be captured.
	}

	return false, false, false
}

func lineMove(b Board, row, col, rowDiff, colDiff int, recurse bool) (possible []ValidMove) {
	color := b.Spaces[row][col].Color
	toRow, toCol := row, col

	for {
		toRow += rowDiff
		toCol += colDiff

		valid, capture, check := tryMove(b, color, row, col, toRow, toCol, recurse)

		if !valid {
			break
		} else if capture {
			possible = append(possible, ValidMove{
				From:    Coord{row, col},
				To:      Coord{toRow, toCol},
				Capture: capture,
				Check:   check,
			})
			break
		}

		possible = append(possible, ValidMove{
			From:    Coord{row, col},
			To:      Coord{toRow, toCol},
			Capture: capture,
			Check:   check,
		})
	}

	return
}

func tryAndAppend(vm []ValidMove, b Board, row, col, rowDiff, colDiff int, recurse bool) []ValidMove {
	color := b.Spaces[row][col].Color

	valid, capture, check := tryMove(b, color, row, col, row+rowDiff, col+colDiff, recurse)
	if valid {
		return append(vm, ValidMove{
			From:    Coord{row, col},
			To:      Coord{row + rowDiff, col + colDiff},
			Capture: capture,
			Check:   check,
		})
	}

	return vm
}

// movePossible returns whether piece can move from row,col to destRow,destCol.
func movePossible(b Board, piece Piece, row, col, destRow, destCol int) bool {
	possible := PossibleMoves(b, piece, row, col, false)
	for _, move := range possible {
		if move.To.Row == destRow && move.To.Col == destCol {
			return true
		}
	}
	return false
}

// NumCheckingKing returns the number of pieces that are putting the king in check
// if count is true, otherwise it returns just 1 if the king is at all in check.
func NumCheckingKing(b Board, c Color, count bool) int {
	var kingPos Coord
	total := 0

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
					if !count {
						return 1
					} else {
						total += 1
					}
				}
			}
		}
	}

	return total
}

// ValidMove represents a possible move that has not necessarily been made.
type ValidMove struct {
	From, To  Coord
	Capture   bool
	EnPassant bool
	Check     bool
}
