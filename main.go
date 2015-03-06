package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mholt/chessml/pgn"
)

func main() {
	// TODO: Parse command line flags, etc...

	f, err := os.Open("Adams.pgn")
	if err != nil {
		panic(err)
	}

	games, err := pgn.Parse(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d games loaded\n", len(games))

	game := games[0]

	fmt.Println("Playing some moves")
	err = game.Execute(7)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(game.Board)
}
