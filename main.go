package main

import (
	"fmt"
	"os"

	"github.com/mholt/chessml/pgn"

	"github.com/mholt/chessml/chess"
)

func main() {
	// TODO: Parse command line flags, etc...

	f, err := os.Open("Adams.pgn")
	if err != nil {
		panic(err)
	}

	g, err := pgn.Parse(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d games loaded!\n", len(g))

	board := chess.Board{}
	board.Setup()
}
