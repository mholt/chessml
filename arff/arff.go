package arff

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mholt/chessml/analysis"
	"github.com/mholt/chessml/chess"
)

func GenerateARFF(games []chess.Game, pctMoves float64) {
	f, err := os.Create("data/chess2.arff")
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
	f.WriteString("@attribute class          REAL\n\n")
	f.WriteString("@data\n%%\n%% " + strconv.Itoa(len(games)) + " instances\n%%\n")

	for i := 0; i < len(games); i++ {
		numMoves := int(float64(len(games[i].Moves)) * pctMoves)
		games[i].Execute(numMoves)

		material := float64((analysis.Material(games[i], chess.WhiteTeam) + 1)) / float64((analysis.Material(games[i], chess.BlackTeam) + 1))
		attackValue := float64((analysis.AttackValue(games[i], chess.WhiteTeam) + 1)) / float64((analysis.AttackValue(games[i], chess.BlackTeam) + 1))
		mobility := float64((analysis.Mobility(games[i], chess.WhiteTeam) + 1)) / float64((analysis.Mobility(games[i], chess.BlackTeam) + 1))
		space := float64((analysis.Space(games[i], chess.WhiteTeam) + 1)) / float64((analysis.Space(games[i], chess.BlackTeam) + 1))

		// For now, we assume that we are training to predict WHITE's move (ie. it's white's turn)
		var outcome float64
		switch games[i].Tags["Result"] {
		case chess.WhiteWin:
			outcome = float64(analysis.Material(games[i], chess.WhiteTeam)) / float64(analysis.Material(games[i], chess.BlackTeam))
		case chess.BlackWin:
			outcome = -float64(analysis.Material(games[i], chess.BlackTeam)) / float64(analysis.Material(games[i], chess.WhiteTeam))
		default:
			outcome = 0
		}

		f.WriteString(fmt.Sprintf("%f,%f,%f,%f,%f\n", material, attackValue, mobility, space, outcome))
	}

	f.WriteString("%%\n%%\n%%\n")

	f.Sync()
}
