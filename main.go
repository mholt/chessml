package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/mholt/chessml/arff"
	"github.com/mholt/chessml/chess"
	"github.com/mholt/chessml/pgn"
)

const numGames = 2000

func main() {
	fmt.Printf("Loading %d random games\n", numGames)
	games := loadRandomGames("pgnfiles/", numGames)

	fmt.Print("\nSnapshotting each game and writing ARFF file...")
	arff.GenerateARFF(games, []float64{0.75}, "data/chess75.arff")
	fmt.Println(" done!")
}

// loadRandomGames will load at most n games from
// any .pgn files in the directory dir. It will
// traverse the directory to child directories
// searching as well. The games are randomly
// chosen. This function is O(n) because it uses
// reservoir sampling.
func loadRandomGames(dir string, n int) []chess.Game {
	var games = make([]chess.Game, 0, n)
	var k int

	rand.Seed(time.Now().UnixNano())

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) != ".pgn" {
			return nil
		}

		fmt.Println(path)

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		fgames, err := pgn.Parse(f)
		if err != nil {
			log.Fatal(err)
		}

		for _, game := range fgames {
			k++

			if k <= n {
				// First fill up the reservoir
				games = append(games, game)
			} else {
				// Otherwise keep each new element with probability n/k
				rnd := rand.Intn(k)
				if rnd < n-1 {
					games[rnd] = game
				}
			}

		}

		return nil
	})

	return games
}
