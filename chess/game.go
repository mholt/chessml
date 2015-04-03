package chess

import (
	"errors"
	"fmt"
)

// A Game represents a chess game.
type Game struct {
	Tags    map[string]string
	Moves   []Move
	Board   Board
	moveIdx int
}

// Execute plays n moves of the game or until the game
// runs out of moves. It does nothing if the game has ended.
// Pass in -1 to play all the moves.
func (g *Game) Execute(n int) error {
	if g.moveIdx >= len(g.Moves) {
		return nil // no moves left
	}

	if n < 0 {
		n = len(g.Moves) - g.moveIdx // do all moves
	}

	for i := 0; i < n && g.moveIdx < len(g.Moves); i++ {
		move := g.Moves[g.moveIdx]
		err := g.move(move)
		if err != nil {
			return fmt.Errorf("Turn %d %s (%s) move %d - %s", g.moveIdx/2+int(move.PlayerColor)-1, move.Player, move.Text, g.moveIdx, err)
		}
		g.moveIdx++
	}

	return nil
}

// move executes the move m.
func (g *Game) move(m Move) error {
	pm, err := m.Parse()
	if err != nil {
		return err
	}

	// Find the piece that can satisfy the move
	_, row, col, found := g.findPiece(pm)
	if !found {
		fmt.Println(g.Board)
		fmt.Printf("Parsed Move: %#v\n", pm)
		return errors.New("Couldn't find any piece to satisfy the move '" + m.Text + "'")
	}

	// Departure coordinate
	from := Coord{Row: row, Col: col}

	// Handle castles a little differently
	if pm.Castle == KingsideCastle {
		g.Board.MovePiece(from, Coord{Row: row, Col: col + 2})              // King
		g.Board.MovePiece(Coord{Row: row, Col: 7}, Coord{Row: row, Col: 5}) // Rook
		return nil
	} else if pm.Castle == QueensideCastle {
		g.Board.MovePiece(from, Coord{Row: row, Col: col - 2})              // King
		g.Board.MovePiece(Coord{Row: row, Col: 0}, Coord{Row: row, Col: 3}) // Rook
		return nil
	}

	// Destination coordinate
	to := NotationToCoord(pm.Destination)

	// Execute the move
	_, err = g.Board.MovePiece(from, to)
	if err != nil {
		return err
	}

	// If it was a pawn promotion, promote it!
	if pm.PawnPromotion != Empty {
		g.Board.Spaces[to.Row][to.Col].Rank = pm.PawnPromotion
	}

	// If it was en passant, remove the captured piece
	if pm.EnPassant {
		g.Board.Spaces[from.Row][to.Col].Rank = Empty
	}

	return nil
}

// findPiece looks on the board to find a piece can satisfy the move
// specified. It returns the first piece that satisfies the move;
// there should only be one (otherwise the movetext was ambiguous).
// This method only filters by criteria that are specified
// (non-zero values). If the piece was found, it returns that piece
// and its row,col position, and true. Otherwise, returns false.
func (g *Game) findPiece(pm *ParsedMove) (Piece, int, int, bool) {
	departRow := rankToRow[pm.DepartureRank]
	departCol := fileToCol[pm.DepartureFile]
	destRow := rankToRow[pm.DestinationRank]
	destCol := fileToCol[pm.DestinationFile]

	for row := 0; row < Size; row++ {
		if pm.DepartureRank != "" && row != departRow {
			continue
		}

		for col := 0; col < Size; col++ {
			if pm.DepartureFile != "" && col != departCol {
				continue
			}

			piece := g.Board.Spaces[row][col]

			if (pm.PieceType > 0 && piece.Rank != pm.PieceType) ||
				piece.Rank == Empty || piece.Color != pm.Color {
				continue
			}

			// Handle castling moves a little differently
			// TODO: Check to make sure the castling is possible/allowed?
			if pm.Castle != "" {
				return piece, row, col, true
			}

			if movePossible(g.Board, piece, row, col, destRow, destCol) {
				// First, see if it's an en passant; the movetext probably won't specify
				if piece.Rank == Pawn && g.Board.Spaces[destRow][destCol].Rank == Empty && pm.Capture {
					pm.EnPassant = true
				}

				// First, make sure doing this move won't put the player in check
				boardCopy := g.Board.copy()
				boardCopy.MovePiece(Coord{Row: row, Col: col}, Coord{Row: destRow, Col: destCol})

				// En passant doesn't directly replace the captured piece which, to our
				// program, would not appear to remove the player out of check if they are
				// in it; so we have to simulate that before the check for check...
				if pm.EnPassant {
					boardCopy.Spaces[row][destCol].Rank = Empty
				}

				if NumCheckingKing(boardCopy, pm.Color, false) > 0 {
					continue // Not allowed; find another piece
				}

				return piece, row, col, true
			}
		}
	}

	return Piece{}, -1, -1, false
}
