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

	game := games[18] // In game 18, white's move on turn 46 'fxg6' is en passant
	err = game.Execute(-1)
	if err != nil {
		log.Fatal(err)
	}

	/*
		for i, game := range games {
			fmt.Println("Playing game", i)
			err = game.Execute(-1)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Played all games!")
	*/
}
