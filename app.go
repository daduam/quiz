package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	sc = bufio.NewReader(os.Stdin)
)

func checkResponse(s string, answer string) bool {
	return strings.TrimSpace(s) == strings.TrimSpace(answer)
}

type quiz struct {
	questions [][]string
	correct   int
}

// LoadQuestions loads questions from a csv file
func (q *quiz) LoadQuestions(filename string, shuffle bool) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}
	q.questions = records
}

// Start handles gameplay
func (q *quiz) Start(seconds int) {
	fmt.Println("Quiz game")
	fmt.Print("Press Enter/Return to begin...")
	_, err := sc.ReadString('\n')
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(time.Duration(seconds) * time.Second)

	for index, question := range q.questions {
		fmt.Printf("Problem #%d: %s = ", index+1, question[0])
		c := make(chan string)
		go func() {
			response, err := sc.ReadString('\n')
			if err != nil {
				panic(err)
			}
			c <- response
		}()

		select {
		case <-timer.C:
			println()
			goto Announce
		case response := <-c:
			if checkResponse(response, question[1]) {
				q.correct++
			}
		}
	}

Announce:
	fmt.Printf("You scored %d out of %d.\n", q.correct, len(q.questions))
}

func main() {
	// flags
	csvPtr := flag.String("csv", "./problems.csv", "path to csv file in the format `question,answer` (defaults to problems.csv)")
	limitPtr := flag.Int("limit", 30, "time limit for the quiz in seconds")
	shufflePtr := flag.Bool("shuffle", false, "shuffle problems")
	flag.Parse()

	q := &quiz{}
	q.LoadQuestions(*csvPtr, *shufflePtr)
	q.Start(*limitPtr)
}
