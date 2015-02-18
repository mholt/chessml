package chess

// A Game represents a chess game.
type Game struct {
	Tags  map[string]string
	Moves []Move
	Board Board
}
