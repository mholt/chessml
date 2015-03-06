package chess

import "fmt"

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
			return fmt.Errorf("Move %d (%s - %s): %s", g.moveIdx, move.Player, move.Text, err)
		}
		g.moveIdx++
	}

	return nil
}

func (g *Game) move(m Move) error {
	// Movetext comes in a variety of forms...

	if len(m.Text) == 2 {
		// Movetext has only the destination spot; it is a pawn movement
		// EXAMPLES: e4, e6, b5

		coord := NotationToCoord(m.Text)

		var movedPiece bool
		for rank := 0; rank < Size && !movedPiece; rank++ {
			for file := 0; file < Size && !movedPiece; file++ {
				piece := g.Board.Spaces[rank][file]
				if piece.Rank != Pawn || piece.Color != m.PlayerColor {
					continue
				}

				possible := PossibleMoves(g.Board, piece, rank, file)
				for _, possMove := range possible {
					if possMove.To.Row == coord.Row && possMove.To.Col == coord.Col {
						err := g.Board.MovePiece(Coord{rank, file}, possMove.To)
						if err != nil {
							return err
						}
						movedPiece = true
						break
					}
				}

			}
		}

		return nil
	}

	if len(m.Text) == 3 {
		// Movetext indicates piece being moved; not a pawn
		// EXAMPLES: Nd2, Nc6, Kb1, Qc2

		symbol := m.Text[:1]
		coord := NotationToCoord(m.Text[1:])

		// Look for a piece of that type that can move there
		var movedPiece bool
		for rank := 0; rank < Size && !movedPiece; rank++ {
			for file := 0; file < Size && !movedPiece; file++ {
				piece := g.Board.Spaces[rank][file]
				if piece.Rank != SymbolToRank[symbol] || piece.Color != m.PlayerColor {
					continue
				}

				// TODO: This is the same as above; perhaps should be its own function
				possible := PossibleMoves(g.Board, piece, rank, file)
				for _, possMove := range possible {
					if possMove.To.Row == coord.Row && possMove.To.Col == coord.Col {
						err := g.Board.MovePiece(Coord{rank, file}, possMove.To)
						if err != nil {
							return err
						}
						movedPiece = true
						break
					}
				}

			}
		}
	}

	return nil
}
