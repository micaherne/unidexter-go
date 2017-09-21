package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/micaherne/unidexter-go/board"
)

var (
	debug = false
)

func main() {
	var b *board.Board
	reader := bufio.NewReader(os.Stdin)

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			text := scanner.Text()
			// log.Printf("Reading from subprocess: %s", text)

			commandParts := strings.SplitN(text, " ", 2)

			switch commandParts[0] {
			case "uci":
				if b != nil {
					fmt.Println("uci must be first command")
					return
				}
				fmt.Println("id name Unidexter 0.0.1")
				fmt.Println("id author Michael Aherne")

				// TODO: option commands
				fmt.Println("uciok")
			case "debug":
				if commandParts[1] == "on" {
					debug = true
				} else if commandParts[1] == "off" {
					debug = false
				}
			case "isready":
				fmt.Println("readyok")
			case "setoption":
				// TODO: Support options
			case "register":
				// Not required
			case "ucinewgame":
				// TODO: Flush caches etc.
			case "position":
				// TODO: Error checking
				var fen string
				var moves []string
				positionParts := strings.Split(commandParts[1], " ")
				if positionParts[0] == "startpos" {
					fen = board.InitialPositionFEN
					if len(positionParts) > 1 && positionParts[1] == "moves" {
						moves = positionParts[2:]
					}
				}
				for i, part := range positionParts {
					if part == "moves" {
						if fen == "" {
							fen = strings.Join(positionParts[:i-1], " ")
						}
						moves = positionParts[i+1:]
						break
					}
				}
				if fen == "" {
					fen = strings.Join(positionParts, " ")
				}

				b = board.FromFEN(fen)
				for _, move := range moves {
					board.MakeMoveFromNotation(b, move)
				}

			case "go":
				// TODO: Start calculating
				bestMove := board.NegamaxAlphaBeta(b, 5)
				fmt.Printf("bestmove %s\n", bestMove)
				// TODO: What if we don't find a decent move?
			case "stop":
				// TODO: Stop calculating
			case "ponderhit":
				// TODO: Start calculating expected move
			case "quit":
				os.Exit(0)
			}
		}
	}(reader)

	for {
		time.Sleep(time.Second)
	}

}
