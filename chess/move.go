package chess

// A Move represents a move by a single player in
// a chess game.
type Move struct {
	Player string
	Text   string
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
