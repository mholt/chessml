// Package analysis provides functions for examining a chess game
// and returning values related to the current state/progress of
// the game. The functions in this package should not modify a
// chess game; they should be read-only.
package analysis

import "github.com/mholt/chessml/chess"

func PointValue(p Piece) string {
	switch p.Rank {
	case King:
		return 0
	case Queen:
		return 9
	case Bishop:
		return 3
	case Knight:
		return 3
	case Rook:
		return 5
	case Pawn:
		return 1
	default:
		return 1
	}
}

// Material computes the point value of pieces on the board
// for either WhiteTeam or BlackTeam.
func Material(game chess.Game, player chess.Color) int {
	// TODO
	total := 0

	b * Board = game.Board
	for c := 0; c < Size; c++ {
		for r := 0; r < Size; r++ {
			if b.Spaces[r][c].Rank != Empty && b.Spaces[r][c].Color == player {
				total += PointValue(b.Spaces[r][c])
			}
		}
	}

	return total
}

// AttackValue computes the point value of all the opponent's pieces that
// may be captured right now by either WhiteTeam or BlackTeam attacking.
func AttackValue(game chess.Game, attacking chess.Color) int {
	// TODO
	total := 0

	b * Board = game.Board
	for c := 0; c < Size; c++ {
		for r := 0; r < Size; r++ {
			if b.Spaces[r][c].Rank != Empty && b.Spaces[r][c].Color == attacking {
				possibleMoves := PossibleMoves(b, b.Spaces[r][c], r, c)
				for m := 0; m < len(possibleMoves); m++ {
					if possibleMoves[m].Capture {
						total += PointValue(b.Spaces[possibleMoves[m].To.Row][possibleMoves[m].To.Col])
					}
				}
			}
		}
	}

	return total
}

// Mobility computes the number of moves possible right now
// for either WhiteTeam or BlackTeam.
func Mobility(game chess.Game, player chess.Color) int {

	// TODO
	total := 0

	b * Board = game.Board
	for c := 0; c < Size; c++ {
		for r := 0; r < Size; r++ {
			if b.Spaces[r][c].Rank != Empty && b.Spaces[r][c].Color == player {
				possibleMoves := PossibleMoves(b, b.Spaces[r][c], r, c)
				total += len(possibleMoves)
			}
		}
	}

	return total
}

// Space computes the number of spaces controlled/protected by
// WhiteTeam or BlackTeam.
func Space(game chess.Game, player chess.Color) int {
	// TODO
	return 0
}

// TODO: Functions for any other features we want to use for our learning algorithm
