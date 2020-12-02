package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	count   int = 0
	correct int = 0
)

func main() {
	fmt.Println("Quiz game")

	f, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewReader(os.Stdin)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		count++
		fmt.Printf("Problem #%d: %s = ", count, record[0])
		response, err := sc.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if strings.TrimSpace(response) == strings.TrimSpace(record[1]) {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, count)
}
