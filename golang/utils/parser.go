package utils

import (
	"fmt"
	"github.com/integrii/flaggy"
	"os"
	"strconv"
)

func (args *CmdlineArgs) SetupCmdlineArgs() []QuestionData {
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.Int(&args.StartIdx, "s", "startIdx", "Start from Question Index")
	flaggy.Int(&args.StopIdx, "e", "stopIdx", "Stop at Question Index")
	flaggy.Int(&args.TestDuration_m, "t", "dur", "Duration of test")
	flaggy.Bool(&args.Key_mode, "k", "key", "Answer Key mode")
	flaggy.Bool(&args.CheckResponses, "c", "check", "Check Responses Against Answer Key")
	flaggy.Parse()
	args.NumQuestions = args.StopIdx - args.StartIdx + 1
	args.TestDuration_s = args.TestDuration_m * 60
	questionBank := make([]QuestionData, (args.StopIdx - args.StartIdx + 1))
	return questionBank
}

func PushQuestionBankToCsv(qData []QuestionData, filename string) [][]string {
	headers := []string{"Question", "Response Time (s)", "Cumulative Time (s)", "Response", "Correct Answer", "Result"}
	data := make([][]string, len(qData))
	cumulative_time := 0.0
	for i, q := range qData {
		cumulative_time = cumulative_time + q.ResponseTime
		data[i] = []string{
			fmt.Sprintf("%d", q.QuestionNum),
			fmt.Sprintf("%f", q.ResponseTime),
			fmt.Sprintf("%f", cumulative_time),
			q.Response,
			q.Answer,
			q.Correct}
	}
	WriteCSV(filename, append([][]string{headers}, data...))
	return append([][]string{headers}, data...)
}

func CsvToQuestionData(responses_raw [][]string, args CmdlineArgs) []QuestionData {
	var responses = make([]QuestionData, 0)
	for _, validOpt := range responses_raw {
		var localQ QuestionData
		idx, _ := strconv.Atoi(validOpt[0])
		if idx >= args.StartIdx && idx <= args.StopIdx {
			localQ.QuestionNum = idx
			localQ.ResponseTime, _ = strconv.ParseFloat(validOpt[1], 32)
			localQ.Response = validOpt[3]
			localQ.Answer = validOpt[4]
			localQ.Correct = validOpt[5]
			responses = append(responses, localQ)
		}
	}
	return responses
}

func CheckValidResponse(response string) error {
	valid_responses := []string{"a", "b", "c", "d", "e"}
	for _, validOpt := range valid_responses {
		if response == validOpt {
			return nil
		}
	}
	return os.ErrInvalid
}
