package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Line struct {
	Answer   string
	Question string
}

func main() {

	questions := loadQuestions("db.csv")

	for _, item := range questions {

		timer := time.NewTimer(time.Duration(30) * time.Second)

		ch := make(chan bool, 1)
		defer close(ch)
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("What %s, sir ? \n", item.Question)

		go func(t time.Timer, ch chan bool) {
			for {
				select {
				case <-timer.C:
					log.Fatalf("No answer")

				case v := <-ch:
					if v {
						return
					}
				case <-time.Tick(time.Duration(5) * time.Second):
					fmt.Println("(Every 5 seconds) Answer the question: ")
				}
			}
		}(*timer, ch)

		answer, _ := reader.ReadString('\n')
		ch <- true

		if strings.TrimSpace(item.Answer) != strings.TrimSpace(answer) {
			log.Fatalln("Test failed!")
			break
		}
	}

	fmt.Println("Succesfully")
}

func loadQuestions(csvFile string) []Line {
	db, err := os.Open(csvFile)

	if err != nil {
		log.Fatalf("File with question is abcent: %s", err)
	}

	defer db.Close()

	lines, err := csv.NewReader(db).ReadAll()

	if err != nil {
		log.Fatalf("Cannot read db, check csv file: %s", err)
	}

	result := []Line{}

	for _, line := range lines {
		result = append(result, Line{
			Answer:   line[1],
			Question: line[0],
		})
	}

	return result
}
