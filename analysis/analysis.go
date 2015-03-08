// Package analysis provides functions for examining a chess game
// and returning values related to the current state/progress of
// the game. The functions in this package should not modify a
// chess game; they should be read-only.
package analysis

import "github.com/mholt/chessml/chess"

// Material computes the point value of pieces on the board
// for either WhiteTeam or BlackTeam.
func Material(game chess.Game, player chess.Color) int {
	// TODO
	return 0
}

// AttackValue computes the point value of all the opponent's pieces that
// may be captured right now by either WhiteTeam or BlackTeam attacking.
func AttackValue(game chess.Game, attacking chess.Color) int {
	// TODO
	return 0
}

// Mobility computes the number of moves possible right now
// for either WhiteTeam or BlackTeam.
func Mobility(game chess.Game, player chess.Color) int {
	// TODO
	return 0
}

// Space computes the number of spaces controlled/protected by
// WhiteTeam or BlackTeam.
func Space(game chess.Game, player chess.Color) int {
	// TODO
	return 0
}

// TODO: Functions for any other features we want to use for our learning algorithm
