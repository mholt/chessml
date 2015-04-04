package arff

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mholt/chessml/analysis"
	"github.com/mholt/chessml/chess"
)

// GenerateARFF makes an ARFF file by snapshotting the games at
// pctMoves of the way through the game. In other words, a pctMoves
// of 0.7 will play 70% of the moves and then snapshot the game,
// writing a line into the ARFF file.
func GenerateARFF(games []chess.Game, pctMoves float64) {
	f, err := os.Create("data/chess.arff")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString("%% Title: Database for predicting chess outcomes\n\n")
	f.WriteString("@relation chess\n\n")
	f.WriteString("@attribute material       REAL\n")
	f.WriteString("@attribute attack-value   REAL\n")
	f.WriteString("@attribute mobility       REAL\n")
	f.WriteString("@attribute space          REAL\n")
	f.WriteString("@attribute currentCheck   REAL\n")
	f.WriteString("@attribute putInCheck     REAL\n")
	f.WriteString("@attribute class          REAL\n\n")
	f.WriteString("@data\n%%\n%% " + strconv.Itoa(len(games)) + " instances\n%%\n")

	for _, game := range games {
		numMoves := int(float64(len(game.Moves)) * pctMoves)
		err := game.Execute(numMoves)
		if err != nil {
			log.Println(err)
			log.Println("^ Skipping that game")
			continue
		}

		material := (analysis.Material(game, chess.WhiteTeam) + 1) / (analysis.Material(game, chess.BlackTeam) + 1)
		attackValue := (analysis.AttackValue(game, chess.WhiteTeam) + 1) / (analysis.AttackValue(game, chess.BlackTeam) + 1)
		mobility := (analysis.Mobility(game, chess.WhiteTeam) + 1) / (analysis.Mobility(game, chess.BlackTeam) + 1)
		space := (analysis.Space(game, chess.WhiteTeam) + 1) / (analysis.Space(game, chess.BlackTeam) + 1)
		currentCheck := (analysis.CurrentCheck(game, chess.WhiteTeam) + 1) / (analysis.CurrentCheck(game, chess.BlackTeam) + 1)
		putInCheck := (analysis.PutInCheck(game, chess.WhiteTeam) + 1) / (analysis.PutInCheck(game, chess.BlackTeam) + 1)

		// Finish executing game to know how extreme the win/loss is
		err = game.Execute(-1)
		if err != nil {
			log.Println(err)
			log.Println("^ Skipping that game")
			continue
		}

		// For now, we assume that we are training to predict WHITE's move (ie. it's white's turn)
		var outcome float64
		switch game.Tags["Result"] {
		case chess.WhiteWin:
			outcome = float64(analysis.Material(game, chess.WhiteTeam)) / float64(analysis.Material(game, chess.BlackTeam))
		case chess.BlackWin:
			outcome = -float64(analysis.Material(game, chess.BlackTeam)) / float64(analysis.Material(game, chess.WhiteTeam))
		default:
			outcome = 0
		}

		f.WriteString(fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f\n", material, attackValue, mobility, space, currentCheck, putInCheck, outcome))
	}

	f.WriteString("%%\n%%\n%%\n")

	f.Sync()
}
