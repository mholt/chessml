package arff

import (
	"bufio"
	"fmt"
	"github.com/mholt/chessml/chess"
	"io/ioutil"
	"os"
	"strconv"
)

func GenerateARFF(games []Game, numMoves int) {
	f, err := os.Create("/data/chess.arff")
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
	f.WriteString("@data\n%%\n%% " + FormatInt(len(games), 10) + " instances\n%%\n")

	for i := 0; i < len(games); i++ {
		games[i].Execute(numMoves)

		material = Material(games[i], WhiteTeam)
		attackValue = AttackValue(games[i], WhiteTeam)
		mobility = Mobility(games[i], WhiteTeam)
		space = Space(games[i], WhiteTeam)
		outcome = games[i].Tags["Result"]

		f.WriteString(FormatInt(material, 10) + "," + FormatInt(attackValue, 10) + "," + FormatInt(mobility, 10) + "," + FormatInt(space, 10) + "," + outcome)
	}

	f.WriteString("%%\n%%\n%%\n")

	f.Sync()
}
