// Package analysis provides functions for examining a chess game
// and returning values related to the current state/progress of
// the game. The functions in this package should not modify a
// chess game; they should be read-only.
package analysis

import "github.com/mholt/chessml/chess"

// Material computes the point value of pieces on the board
// for either WhiteTeam or BlackTeam.
func Material(game chess.Game, player chess.Color) float64 {
	total := 0.0

	for c := 0; c < chess.Size; c++ {
		for r := 0; r < chess.Size; r++ {
			if game.Board.Spaces[r][c].Rank != chess.Empty && game.Board.Spaces[r][c].Color == player {
				total += PointValue(game.Board.Spaces[r][c])
			}
		}
	}

	return total
}

// AttackValue computes the point value of all the opponent's pieces that
// may be captured right now by either WhiteTeam or BlackTeam attacking.
func AttackValue(game chess.Game, attacking chess.Color) float64 {
	total := 0.0

	for c := 0; c < chess.Size; c++ {
		for r := 0; r < chess.Size; r++ {
			if game.Board.Spaces[r][c].Rank != chess.Empty && game.Board.Spaces[r][c].Color == attacking {
				possibleMoves := chess.PossibleMoves(game.Board, game.Board.Spaces[r][c], r, c, false)
				for _, move := range possibleMoves {
					if move.Capture {
						total += PointValue(game.Board.Spaces[move.To.Row][move.To.Col])
					}
				}
			}
		}
	}

	return total
}

// Mobility computes the number of moves possible right now
// for either WhiteTeam or BlackTeam.
func Mobility(game chess.Game, player chess.Color) float64 {
	total := 0.0

	for c := 0; c < chess.Size; c++ {
		for r := 0; r < chess.Size; r++ {
			if game.Board.Spaces[r][c].Rank != chess.Empty && game.Board.Spaces[r][c].Color == player {
				possibleMoves := chess.PossibleMoves(game.Board, game.Board.Spaces[r][c], r, c, false)
				total += float64(len(possibleMoves))
			}
		}
	}

	return total
}

// Space computes the number of spaces controlled/protected by
// WhiteTeam or BlackTeam.
func Space(game chess.Game, player chess.Color) float64 {
	// the spaces that the player controls on the other player's half of the board
	total := 0.0

	for c := 0; c < chess.Size; c++ {
		for r := 0; r < chess.Size; r++ {
			if game.Board.Spaces[r][c].Rank != chess.Empty && game.Board.Spaces[r][c].Color == player {
				possibleMoves := chess.PossibleMoves(game.Board, game.Board.Spaces[r][c], r, c, false)
				for _, move := range possibleMoves {
					// the space is only controlled when the piece can make a move that is attacking, so for pawns, it must not be the same column
					if BoardHalfColor(move.To.Row) != player && (game.Board.Spaces[r][c].Rank != chess.Pawn || move.To.Col != c) {
						total += PointValue(game.Board.Spaces[move.To.Row][move.To.Col])
					}
				}
			}
		}
	}

	return total
}

// CurrentCheck computes the number of pieces that one team currently has the other in check with
func CurrentCheck(game chess.Game, player chess.Color) float64 {
	return float64(chess.NumCheckingKing(game.Board, player, true))
}

// PutInCheck computes the number of moves that one team could currently make to put the other in check
func PutInCheck(game chess.Game, player chess.Color) float64 {
	total := 0.0

	for c := 0; c < chess.Size; c++ {
		for r := 0; r < chess.Size; r++ {
			if game.Board.Spaces[r][c].Rank != chess.Empty && game.Board.Spaces[r][c].Color == player {
				possibleMoves := chess.PossibleMoves(game.Board, game.Board.Spaces[r][c], r, c, true)
				for _, move := range possibleMoves {
					if move.Check {
						total += 1.0
					}
				}
			}
		}
	}

	return total
}

// TODO: Functions for any other features we want to use for our learning algorithm

// helper functions

// white's half is row indices 0-3
// black's half is row indices 4-7
func BoardHalfColor(row int) chess.Color {
	if row <= 3 {
		return chess.WhiteTeam
	}
	return chess.BlackTeam
}

func PointValue(p chess.Piece) float64 {
	switch p.Rank {
	case chess.King:
		return 0.1
	case chess.Queen:
		return 8.8
	case chess.Bishop:
		return 3.33
	case chess.Knight:
		return 3.2
	case chess.Rook:
		return 5.1
	case chess.Pawn:
		return 1.0
	default:
		return 1.0
	}
}
