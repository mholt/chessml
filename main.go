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

	for i, game := range games {
		fmt.Println("Playing game", i)
		err = game.Execute(-1)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Played all games!")

}
