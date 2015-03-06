package pgn

import (
	"bufio"
	"io"

	"github.com/mholt/chessml/chess"
)

// Parse parses the PGN file format from input into games.
// A file may contain zero or more chess games.
func Parse(input io.Reader) (games []chess.Game, err error) {
	var game chess.Game
	var done bool

	// The scanner reads the input
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanRunes)

	// The parser assembles the games
	parser := newGameParser(scanner)

	for {
		game, done, err = parser.parseGame()
		if err != nil {
			return games, err
		} else if done {
			break
		}

		game.Board.Setup()
		games = append(games, game)
	}

	return games, nil
}
