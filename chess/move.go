package chess

import (
	"errors"
	"strings"
)

// A Move represents a move by a single player in
// a chess game.
type Move struct {
	Player      string
	PlayerColor Color
	Text        string
}

// A parsed move represents movetext that is more usable
// by the program. Only values that are directly inferred
// from the movetext will be filled out; the rest will be
// the zero-value. One exception is the Color, which is
// copied in from the Move struct for convenience.
type ParsedMove struct {
	PieceType       Rank
	Color           Color
	DepartureRank   string
	DepartureFile   string
	DestinationRank string
	DestinationFile string
	Destination     string // concatenated DestinationRank and DestinationFile if both are given
	Check           bool
	Capture         bool
	Castle          string // KingsideCastle or QueensideCastle
	PawnPromotion   Rank   // what the pawn is promoted to
	EnPassant       bool   // Note: Movetext may not specify this; might have to rely on game state
}

// Parse parses the movetext into usable values; those
// values are returned in ParsedMove
func (m Move) Parse() (*ParsedMove, error) {
	pm := &ParsedMove{Color: m.PlayerColor}
	t := m.Text

	if len(t) < 2 {
		return nil, errors.New("Movetext too short")
	}

	if t == "O-O" || t == "0-0" { // PGN uses capital Os, but SAN uses zeros
		// Kingside castle
		pm.PieceType = King
		pm.Castle = KingsideCastle
		return pm, nil
	}

	if t == "O-O-O" || t == "0-0-0" {
		// Queenside castle
		pm.PieceType = King
		pm.Castle = QueensideCastle
		return pm, nil
	}

	if strings.HasSuffix(t, "+") {
		pm.Check = true
		t = t[:len(t)-1]
	}

	switch len(t) {
	case 2:
		return parseTextLen2(t, pm)
	case 3:
		return parseTextLen3(t, pm)
	case 4:
		return parseTextLen4(t, pm)
	case 5:
		return parseTextLen5(t, pm)
	default:
		return pm, errors.New("Unable to parse movetext " + t)
	}
}

// parseTextLen2 parses movetext of length 2 (after stripping +)
func parseTextLen2(t string, pm *ParsedMove) (*ParsedMove, error) {
	pm.PieceType = Pawn

	if isFile[t[1]] {
		// Pawn capture
		// Example: ed
		pm.DepartureFile = t[:1]
		pm.DestinationFile = t[1:]
		pm.Capture = true
	} else if isRank[t[1]] {
		// Simple pawn movement
		// Examples: f4, e5, e6
		pm.DestinationFile = t[:1]
		pm.DestinationRank = t[1:]
		pm.Destination = t
	} else {
		return pm, errors.New("Invalid two-character movetext (" + t + ")")
	}

	return pm, nil
}

// parseTextLen3 parses movetext of length 3 (after stripping +).
func parseTextLen3(t string, pm *ParsedMove) (*ParsedMove, error) {
	if isPiece[t[0]] {
		// Piece type specified along with destination file and rank
		// Examples: Nd2, Qh4, Bf6
		pm.PieceType = SymbolToRank[t[:1]]
		pm.Destination = t[1:]
		pm.DestinationFile = t[1:2]
		pm.DestinationRank = t[2:3]
	} else if isFile[t[0]] && (t[1] == 'x' || t[1] == ':') && isFile[t[1]] {
		// Pawn capture, specified more verbosely using ':' or ':'
		// Examples: exd, e:d
		pm.PieceType = Pawn
		pm.DepartureFile = t[:1]
		pm.DestinationFile = t[2:3]
		pm.Capture = true
	} else {
		return pm, errors.New("Invalid 3-character movetext (" + t + ")")
	}

	return pm, nil
}

// parseTextLen4 parses movetext of length 4 (after stripping +).
func parseTextLen4(t string, pm *ParsedMove) (*ParsedMove, error) {
	if strings.Index(t, "=") == 2 {
		// Special case: pawn promotion!
		// Example: c1=Q
		pm.PieceType = Pawn
		pm.Destination = t[0:2]
		pm.DestinationFile = t[0:1]
		pm.DestinationRank = t[1:2]
		pm.PawnPromotion = SymbolToRank[t[3:4]]
		return pm, nil
	}

	if isPiece[t[0]] {
		// Piece type and rank/file/capture along with destination rank and file
		// Examples: Nfd7, Bxh7, Rde1, Kxd8, Qxd4, R7g5
		pm.PieceType = SymbolToRank[t[:1]]
	} else if isFile[t[0]] {
		// 4-char movetext without piece letter is a pawn capture, according to Wikipedia
		// Examples: fxe5, exf5
		pm.PieceType = Pawn
		pm.DepartureFile = t[:1]
	} else if isRank[t[0]] {
		pm.DepartureRank = t[:1]
	}

	if t[1] == 'x' || t[1] == ':' {
		// 'x' is the official way to indicate a capture, but sometimes ':' may be used
		pm.Capture = true
	} else if isFile[t[1]] {
		pm.DepartureFile = t[1:2]
	} else if isRank[t[1]] {
		pm.DepartureRank = t[1:2]
	} else {
		return pm, errors.New("Not sure how to parse " + t)
	}

	// Last two chars must be the destination file and rank
	pm.DestinationFile = t[2:3]
	pm.DestinationRank = t[3:4]
	pm.Destination = t[2:]

	return pm, nil
}

// parseTextLen5 parses movetext of length 5 (after stripping +).
func parseTextLen5(t string, pm *ParsedMove) (*ParsedMove, error) {
	// 5-char movetext is less common; usually a disambiguated non-pawn capture
	// Examples: N5xf3, Rdxd5, Nbxd4
	pm.PieceType = SymbolToRank[t[:1]]

	// File trumps rank when disambiguating moves, just FYI
	if isFile[t[1]] {
		pm.DepartureFile = t[1:2]
	} else if isRank[t[1]] {
		pm.DepartureRank = t[1:2]
	} else {
		return pm, errors.New("Unable to parse " + t)
	}

	if t[2] == 'x' || t[2] == ':' {
		pm.Capture = true
	} else {
		return pm, errors.New("Unable to parse " + t)
	}

	// Last two chars must be destination
	pm.DestinationFile = t[3:4]
	pm.DestinationRank = t[4:5]
	pm.Destination = t[3:]

	return pm, nil
}

const (
	White = "W"
	Black = "B"
)

const (
	WhiteWin = "1-0"
	BlackWin = "0-1"
	Draw     = "1/2-1/2"
	Other    = "*"
)

const (
	KingsideCastle  = "O-O"
	QueensideCastle = "O-O-O"
)

var (
	isRank  = map[uint8]bool{'1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true}
	isFile  = map[uint8]bool{'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true}
	isPiece = map[uint8]bool{'K': true, 'Q': true, 'B': true, 'N': true, 'R': true}
)
