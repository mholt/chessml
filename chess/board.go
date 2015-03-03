package chess

type Color int

const (
	WhiteTeam Color = iota
	BlackTeam
)

type Rank int

const (
	Empty Rank = iota
	King
	Queen
	Bishop
	Knight
	Rook
	Pawn
)

type Piece struct {
	Color
	Rank
}

const Size = 8

// A Board represents a chess board.
type Board struct {
	Spaces [Size][Size]Piece
}


func (b Board) Setup() {
	// Wipe everything off.
	for c := 0; c < Size; c++ {
		for r := 0; r < Size; r++ {
			b.Spaces[r][c].Rank = Empty
			b.Spaces[r][c].Color = WhiteTeam
		}
	}

	// Place all of the Pawns.
	for c := 0; c < Size; c++ {
		b.Spaces[1][c].Rank = Pawn

		b.Spaces[6][c].Rank = Pawn
		b.Spaces[6][c].Color = BlackTeam
	}

	for c := 0; c < Size; c++ {
		b.Spaces[6][c].Color = BlackTeam
		b.Spaces[7][c].Color = BlackTeam
	}

	// Place the Kings.
	b.Spaces[0][4].Rank = King
	b.Spaces[7][4].Rank = King

	// Place the Queens.
	b.Spaces[0][3].Rank = Queen
	b.Spaces[7][3].Rank = Queen

	// Place the Bishops.
	b.Spaces[0][5].Rank = Bishop
	b.Spaces[0][2].Rank = Bishop
	b.Spaces[7][5].Rank = Bishop
	b.Spaces[7][2].Rank = Bishop

	// Place the Knights
	b.Spaces[0][1].Rank = Knight
	b.Spaces[0][6].Rank = Knight
	b.Spaces[7][1].Rank = Knight
	b.Spaces[7][6].Rank = Knight

	// Place the Rooks
	b.Spaces[0][0].Rank = Rook
	b.Spaces[0][7].Rank = Rook
	b.Spaces[7][0].Rank = Rook
	b.Spaces[7][7].Rank = Rook
}
