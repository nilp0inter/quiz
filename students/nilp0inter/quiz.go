package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

var csvfile = flag.String("csvfile", "problems.csv", "File with questions and answers")
var quiztime = flag.Int("quiztime", 10, "Total time for answering the quiz (in seconds)")

func main() {
	flag.Parse()

	file, err := os.Open(*csvfile)
	if err != nil {
		fmt.Println("Cannot find problems file.")
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	entries, err := reader.ReadAll()
	file.Close()
	if err != nil {
		fmt.Println("Malformed csv file.")
		os.Exit(1)
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(10 * time.Second)
		timeout <- true
	}()

	answers := make(chan bool, 1)
	go func() {
		for _, e := range entries {
			q, a := e[0], e[1]
			fmt.Println(q)
			var u string
			fmt.Scanf("%s", &u)
			answers <- u == a
		}
	}()

	var correct int
	var answered int
	var timedOut bool
	for answered < len(entries) && !timedOut {
		select {
		case <-timeout:
			timedOut = true
		case isCorrect := <-answers:
			answered++
			if isCorrect {
				correct++
			}
		}
	}
	fmt.Printf("\nCorrect answers: %d/%d\n", correct, len(entries))
}
