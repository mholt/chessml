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
	f.WriteString("@attribute class          {won, lost, draw, other}\n\n")
	f.WriteString("@data\n%%\n%% " + FormatInt(len(games), 10) + " instances\n%%\n")

	f.WriteString("%%\n%%\n%%\n")

	f.Sync()
}
