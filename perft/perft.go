package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/micaherne/unidexter-go/board"
)

func mainBim() {
	b := board.FromFEN("8/8/8/8/8/8/1k6/R3K3 b Q - 0 1")
	d := divide(b, 2)
	fmt.Println(d)
}

func main() {

	suite := GetPerftsuite()
	pass := 0
	fail := 0
	max := 6
	for fen, counts := range suite {
		start := time.Now()
		fmt.Println("\n" + fen)
		b := board.FromFEN(fen)
		for i := 1; i <= max && i <= len(counts); i++ {
			p := perft(b, i)

			if p == counts[i-1] {
				pass++
				fmt.Print("PASS. ")
			} else {
				fail++
				fmt.Print("FAIL. ")
			}
			fmt.Printf("Perft(%d). Expected %d, got %d\n", i, counts[i-1], p)
		}
		fmt.Printf("Time taken: %s", time.Since(start))
	}

	fmt.Printf("\nPass: %d, Fail: %d", pass, fail)

}

type Divide map[string]int

func (p Divide) String() string {
	result := ""
	nodes := 0
	for move, count := range p {
		result += fmt.Sprintf("%s %d\n", move, count)
		nodes += count
	}
	result += fmt.Sprintf("Moves: %d\n", len(p))
	result += fmt.Sprintf("Nodes: %d\n", nodes)
	return result
}

func perft(b *board.Board, depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0

	moves := board.GenerateMoves(b)
	for _, move := range moves {
		if !board.LegalMove(b, move) {
			continue
		}
		board.MakeMove(b, move)
		nodes += perft(b, depth-1)
		board.UndoMove(b)
	}

	return nodes
}

func divide(b *board.Board, depth int) Divide {
	result := make(Divide)

	if depth == 0 {
		return result
	}

	moves := board.GenerateMoves(b)
	for _, move := range moves {
		if !board.LegalMove(b, move) {
			continue
		}
		board.MakeMove(b, move)
		nodes := perft(b, depth-1)
		board.UndoMove(b)

		result[move.String()] = nodes
	}

	return result
}

func GetPerftsuite() map[string][]int {
	file := "C:\\dev\\projects\\unidexter-go\\src\\github.com\\micaherne\\unidexter-go\\perft\\perftsuite.epd"
	var content []byte
	content, err := ioutil.ReadFile(file)
	if err != nil {
		// TODO: Do something
		log.Println(err)
		return make(map[string][]int, 0)
	}
	textLines := strings.Split(string(content[:]), "\n")
	perfts := make(map[string][]int)
	for _, line := range textLines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, ";")
		key := parts[0]
		values := make([]int, len(parts)-1)
		for j, result := range parts[1:] {
			resParts := strings.Split(result, " ")
			values[j], err = strconv.Atoi(resParts[1])
			if err != nil {
				// TODO: Report error
			}
		}
		perfts[key] = values
	}

	return perfts
}
