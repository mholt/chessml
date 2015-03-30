package arff

import (
	"github.com/mholt/chessml/analysis"
	"github.com/mholt/chessml/chess"
	"os"
	"strconv"
)

func GenerateARFF(games []chess.Game, numMoves int) {
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
	f.WriteString("@attribute class          {1-0, 0-1, 1/2-1/2, *}\n\n")
	f.WriteString("@data\n%%\n%% " + strconv.Itoa(len(games)) + " instances\n%%\n")

	for i := 0; i < len(games); i++ {
		games[i].Execute(numMoves)
		material := analysis.Material(games[i], chess.WhiteTeam)
		attackValue := analysis.AttackValue(games[i], chess.WhiteTeam)
		mobility := analysis.Mobility(games[i], chess.WhiteTeam)
		space := analysis.Space(games[i], chess.WhiteTeam)
		outcome := games[i].Tags["Result"]

		f.WriteString(strconv.Itoa(material) + "," + strconv.Itoa(attackValue) + "," + strconv.Itoa(mobility) + "," + strconv.Itoa(space) + "," + outcome + "\n")
	}

	f.WriteString("%%\n%%\n%%\n")

	f.Sync()
}
