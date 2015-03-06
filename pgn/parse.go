package pgn

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/mholt/chessml/chess"
)

// newGameParser constructs, well, a new game parser which
// consumes tokens from the scanner. Use one gameParser per
// file so that line and character counts will be accurate.
func newGameParser(scanner *bufio.Scanner) *gameParser {
	return &gameParser{
		scanner: scanner,
		line:    1,
		char:    1,
	}
}

// gameParser is capable of parsing a single game at a time
// from a PGN file being read by its scanner. The scanner
// must not be nil. Use the same gameParser for the entire
// file.
type gameParser struct {
	scanner *bufio.Scanner
	game    chess.Game
	line    int
	char    int
	mv      string
}

// parseGame will parse the input until an entire game is parsed.
// parseGame returns the game, whether there is any more input
// available in the file, and any error that may have occured. If an
// the EOF was encountered (i.e. the middle return is true) or an
// error occured, parsing was aborted abruptly and did not finish.
func (gp *gameParser) parseGame() (chess.Game, bool, error) {
	gp.game = chess.Game{}

	done, err := gp.parseTags()
	if err != nil {
		return gp.game, false, err
	} else if done {
		return gp.game, true, nil
	}

	err = gp.parseMoves()
	if err != nil {
		return gp.game, false, err
	}

	return gp.game, false, nil
}

// parseTags parses the tag pairs portion of a PGN file and
// loads them into the struct. The PGN file format requires
// that games are preceded by certain required tags, but
// this parser can detect if they are missing and skip
// them gracefully. It returns true if there is no more
// input available, i.e. no more tokens to scan.
func (gp *gameParser) parseTags() (bool, error) {
	gp.game.Tags = make(map[string]string)

	for gp.scan() {
		ch := gp.getch()

		if unicode.IsSpace(ch) {
			continue
		}

		if ch == ';' {
			// skip commented line
			for gp.scan() {
				ch = gp.getch()
				if ch == '\n' {
					break // EOL
				}
			}
			continue
		}

		if ch == '1' {
			return false, nil // beginning of movetext
		}

		err := gp.parseTag()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// parseTag parses a tag from the tag pairs portion
// of a PGN file. It expects the currently-loaded token
// to be an opening tag character.
func (gp *gameParser) parseTag() error {
	if gp.getch() != openTag {
		return gp.err("Expecting open tag '" +
			string(openTag) + "' or beginning of movetext")
	}

	tagName, err := gp.parseTagName()
	if err != nil {
		return err
	}

	tagVal, err := gp.parseTagValue()
	if err != nil {
		return err
	}

	gp.game.Tags[tagName] = tagVal

	if gp.getch() != closeTag {
		return gp.err("Expecting close tag '" +
			string(closeTag) + "' token")
	}

	return nil
}

// parseTagName parses the name of a tag from the tag
// pairs portion of a PGN file. It returns the tag name.
// Tag names may not be quoted.
func (gp *gameParser) parseTagName() (string, error) {
	var tagName string

	for gp.scan() {
		ch := gp.getch()

		if unicode.IsSpace(ch) {
			if len(tagName) == 0 {
				continue
			} else {
				break
			}
		}

		tagName += string(ch)
	}

	return tagName, nil
}

// parseTagValue parses the value of a tag from the tag
// pairs portion of a PGN file. It returns the tag value.
// Tag values MUST be quoted. Quotes should be escaped as
// \" and backslashes as \\.
func (gp *gameParser) parseTagValue() (string, error) {
	var tagValue string
	var escaped, quoted bool
	var quotes rune

	for gp.scan() {
		ch := gp.getch()

		if !quoted {
			if ch == '"' || ch == '\'' {
				quoted = true
				quotes = ch
			}
			continue
		}

		if quoted {
			if !escaped {
				if ch == '\\' {
					escaped = true
					continue
				}
				if ch == quotes {
					quoted = false
					break
				}
			}
			escaped = false
		}

		tagValue += string(ch)
	}

	gp.scan() // skip closing quote

	return tagValue, nil
}

// parseMoves parses the movetext of a PGN file.
// It terminates immediately after parsing the result
// of the match. It expects the currently-loaded
// token to be the number '1', which is what indicates
// the beginning of the movetext. Movetext MUST end
// with the result, one of 1-0, 0-1, 1/2-1/2, or *.
func (gp *gameParser) parseMoves() error {
	if gp.getch() != '1' {
		return gp.err("Expecting beginning of movetext")
	}

	gp.scan() // advance past the '1'

	for {
		end, err := gp.parseTurn()
		if err != nil {
			return err
		}

		if end {
			break
		}
	}

	return nil
}

// parseTurn parses a turn, which is two moves (presumably
// one by each player). It expects the currently-loaded token
// to be a period '.' which indicates the beginning of a turn,
// as it immediately follows the turn number. Parsed turns are
// automatically loaded into the game struct. parseTurn returns
// true if all the turns have been parsed for this game.
func (gp *gameParser) parseTurn() (bool, error) {
	var err error

	// Turns start with a '.' which must come after the turn
	// number, but only if we haven't already consumed it.
	if gp.mv == "" && gp.getch() != moveStart {
		return false, gp.err("Expected turn starting with a dot '.'")
	}

	// White's move is the first move of a turn, but we may
	// have already consumed it from the last turn, so use it.
	var whiteMove string
	if gp.mv != "" {
		whiteMove = gp.mv
		gp.mv = ""
	} else {
		whiteMove, err = gp.parseMove()
		if err != nil {
			return false, err
		}
	}

	// Save white's move
	gp.game.Moves = append(gp.game.Moves, chess.Move{
		Player:      chess.White,
		PlayerColor: chess.WhiteTeam,
		Text:        whiteMove,
	})

	// There is always at least one move in a turn, so we
	// don't check for end-of-game yet. The second one,
	// however, might be end-of-game...
	blackMove, err := gp.parseMove()
	if err != nil {
		return false, err
	}

	if blackMove == chess.WhiteWin || blackMove == chess.BlackWin ||
		blackMove == chess.Draw || blackMove == chess.Other {
		// Game over; turns out we just parsed the result instead
		return true, nil
	}

	// Save black's move
	gp.game.Moves = append(gp.game.Moves, chess.Move{
		Player:      chess.Black,
		PlayerColor: chess.BlackTeam,
		Text:        blackMove,
	})

	// This file format is terrible. There might even be a third
	// 'move' in a turn, which is actually just a result. This third
	// 'move' appears when black wins the game or causes a draw.
	// But since it's probably not the end result, we'll have over-consumed
	// so we have to store the actual move, which is the white move of the
	// next turn, in order to make use of it.
	mv, err := gp.parseMove()
	if err != nil {
		return false, err
	}
	if dotIdx := strings.Index(mv, "."); dotIdx > -1 {
		gp.mv = mv[dotIdx+1:]
	} else {
		return true, nil // result indicates end of game; won't have '.' in it
	}

	return false, nil
}

// parseMove parses a single move in a turn. It basically
// collects tokens until a space character is encountered.
func (gp *gameParser) parseMove() (string, error) {
	var mv string

	for gp.scan() {
		ch := gp.getch()

		if unicode.IsSpace(ch) {
			if len(mv) > 0 {
				break
			}
			continue
		}

		mv += string(ch)
	}

	if strings.HasPrefix(mv, "[") {
		return mv, gp.err("Expected a move or end-of-game result; not a tag")
	}

	return mv, nil
}

// scan loads the next token into the scanner; it also
// ups the character position and line number if necessary.
func (gp *gameParser) scan() bool {
	ok := gp.scanner.Scan()
	if ok {
		gp.char++
		if gp.scanner.Text() == "\n" {
			gp.line++
			gp.char = 0
		}
	}
	return ok
}

// err makes an error message with line and character information
// that's useful for debugging.
func (gp *gameParser) err(msg string) error {
	if gp.getch() == 0 {
		return errors.New(fmt.Sprintf("Parse error - line %d, char %d: "+
			"unexpected EOF", gp.line, gp.char))
	}
	return errors.New(fmt.Sprintf("Parse error - line %d, char %d: %s (have '%s')",
		gp.line, gp.char, msg, string(gp.getch())))
}

// getch gets the currently-loaded token (character)
// as a rune. If more than one rune is loaded
// then only the first will be returned.
func (gp *gameParser) getch() rune {
	txt := gp.scanner.Text()
	if len(txt) == 0 {
		return 0
	}
	return []rune(txt)[0]
}

const (
	openTag   = '['
	closeTag  = ']'
	moveStart = '.'
)
