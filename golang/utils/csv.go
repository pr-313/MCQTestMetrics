package utils

import (
	"encoding/csv"
	// "fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func EvalResponses(args CmdlineArgs) {
	// var eval = make([]questionData, numQuestions)
	responses := CsvToQuestionData(ReadCSV("responses.csv"), args)
	key := CsvToQuestionData(ReadCSV("key.csv"), args)

	for i, question := range responses {
	keyIdxSearch:
		for _, ans := range key {
			if question.QuestionNum == ans.QuestionNum {
				if question.Response == ans.Response {
					responses[i].Correct = "Correct"
				} else {
					responses[i].Correct = "Wrong"
				}
				responses[i].Answer = ans.Response
				break keyIdxSearch
			}
		}
	}
	PrintCSV(PushQuestionBankToCsv(responses, "results.csv")[1:])
}

func ReadCSV(filename string) [][]string {
	// Read a CSV file into a 2D string array
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records[1:]
}

func WriteCSV(filename string, data [][]string) {
	// Write a 2D string array to a CSV file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		panic(err)
	}
}

func PrintCSV(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Question", "Response Time (s)", "Cumulative Time (s)", "Response", "Correct Answer", "Result"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
