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
			return fmt.Errorf("Move %d (%s) (%s's turn): %s", g.moveIdx, move.Text, move.Player, err)
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
	err = g.Board.MovePiece(from, to)
	if err != nil {
		return err
	}

	return nil
}

// findPiece looks on the board to find a piece can satisfy the move
// specified. It returns the first piece that satisfies the move;
// there should only be one (otherwise the movetext was ambiguous).
// This method only filters by criteria that are specified
// (non-zero values). If the piece was found, it returns that piece
// and its row,col position, and true. Otherwise, returns false.
func (g *Game) findPiece(pm ParsedMove) (Piece, int, int, bool) {
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
			// TODO: Check to make sure the castling is possible
			if pm.Castle != "" {
				return piece, row, col, true
			}

			if g.movePossible(piece, row, col, destRow, destCol) {
				return piece, row, col, true
			}
		}
	}

	return Piece{}, -1, -1, false
}

// movePossible returns whether piece can move from row,col to destRow,destCol.
func (g *Game) movePossible(piece Piece, row, col, destRow, destCol int) bool {
	possible := PossibleMoves(g.Board, piece, row, col)
	for _, move := range possible {
		if move.To.Row == destRow && move.To.Col == destCol {
			return true
		}
	}
	return false
}
