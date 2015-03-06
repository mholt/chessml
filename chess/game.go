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
		// Movetext has only the destination spot
		// This should happen only when a single piece is capable of moving there,
		// and there is no capture.

		coord := NotationToCoord(m.Text)

		var movedPiece bool
		for rank := 0; rank < Size && !movedPiece; rank++ {
			for file := 0; file < Size && !movedPiece; file++ {
				piece := g.Board.Spaces[rank][file]
				if piece.Rank == Empty || piece.Color != m.PlayerColor {
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

	}

	return nil
}
