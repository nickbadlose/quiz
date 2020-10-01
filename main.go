package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	records := readCsvFile("./problems.csv")

	var answerTime int
	flag.IntVar(&answerTime, "limit", 10, "number of lines to read from the file")
	flag.Parse()

	questionsCorrect := 0

	totalQuestions := len(records)

	for _, v := range records {
		fmt.Println("What is", v[0], "?")

		answerChan := make(chan string)
		timeoutChan := make(chan bool)

		go func() {
			time.Sleep(time.Duration(answerTime) * time.Second)
			timeoutChan <- true
		}()

		go func() {
			var answer string

			_, err := fmt.Scan(&answer)
			if err != nil {
				fmt.Println("Something went wrong")
			}

			answerChan <- string(answer)
		}()

		select {
		case answer := <-answerChan:
			if answer == v[1] {
				questionsCorrect++
			}
		case <-timeoutChan:
			fmt.Println("You got", questionsCorrect, "questions right. Then you ran out of time!")
			close(timeoutChan)
			return
		}
		close(answerChan)
	}

	fmt.Println("You got", questionsCorrect, "questions right out of a possible", totalQuestions, "well done!")

	return
}
