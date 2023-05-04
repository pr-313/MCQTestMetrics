package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/integrii/flaggy"
	"github.com/jroimartin/gocui"
    "github.com/olekukonko/tablewriter"
)

type questionData struct {
	questionNum  int
	answer       string
	response     string
	responseTime float64
	correct      string
}

var (
	startIdx       int
	stopIdx        int
	numQuestions   int
	testDuration_m int
	testDuration_s int
	startTime      time.Time
	lastAnsTime    time.Time
	currIdx        int
	questionBank   = make([]questionData, 1)
	key_mode       bool
	checkResponses bool
)

func main() {
	setupCmdlineArgs()

	if checkResponses {
		evalResponses()
		os.Exit(0)
	}
	// Initialize gocui
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	// Set the view dimensions
	g.SetManagerFunc(layout)

	// Set the keybindings
	err = keybindings(g)
	if err != nil {
		panic(err)
	}

	// Start the test
	startTime = time.Now()
	lastAnsTime = startTime
	currIdx = startIdx
	go runTimer(g)
	g.MainLoop()
}

func setupCmdlineArgs() {
	startIdx = 1
	stopIdx = 10
	testDuration_m = 10
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.Int(&startIdx, "s", "startIdx", "Start from Question Index")
	flaggy.Int(&stopIdx, "e", "stopIdx", "Stop at Question Index")
	flaggy.Int(&testDuration_m, "t", "dur", "Duration of test")
	flaggy.Bool(&key_mode, "k", "key", "Answer Key mode")
	flaggy.Bool(&checkResponses, "c", "check", "Check Responses Against Answer Key")
	flaggy.Parse()
	numQuestions = stopIdx - startIdx + 1
	questionBank = make([]questionData, (stopIdx - startIdx + 1))
	testDuration_s = testDuration_m * 60
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("question", 1, 1, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("Question %d", startIdx)
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "This is question", startIdx)
	}
	if v, err := g.SetView("response", 1, 6, maxX/2-1, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Responses"
		// fmt.Fprintf(v, "Enter your response for question %d here:\n", questionNum)
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = true
	}
	if v, err := g.SetView("textBox", 1, maxY-5, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Type Here"
		// fmt.Fprintf(v, "Enter your response for question %d here:\n", questionNum)
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		g.SetCurrentView("textBox")
	}
	if v, err := g.SetView("timer", maxX/2, 6, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Timer"
		// fmt.Fprintf(v, "Time elapsed: %v", time.Since(startTime).Round(time.Second))
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("textBox", gocui.KeyEnter, gocui.ModNone, nextQuestion); err != nil {
		return err
	}
	return nil
}

func runTimer(g *gocui.Gui) {
	timerChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-timerChan:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("timer")
				if err != nil {
					return err
				}
				v.Clear()
				if key_mode {
					fmt.Fprintf(v, "Answer Key Mode")
				} else {
					fmt.Fprintf(v, "Time spent on this question: %ds\n", int(time.Since(lastAnsTime).Seconds()))
					fmt.Fprintf(v, "Total Time remaining: %ds\n", (testDuration_s - int(time.Since(startTime).Seconds())))
				}
				return nil
			})
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func checkValidResponse(response string) error {
	valid_responses := []string{"a", "b", "c", "d", "e"}
	for _, validOpt := range valid_responses {
		if response == validOpt {
			return nil
		}
	}
	return os.ErrInvalid
}

func checkEOT(g *gocui.Gui, v *gocui.View) error {
	if currIdx > stopIdx {
		// End test
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("question")
			v.Clear()
			v.Title = "Test complete"
			fmt.Fprintf(v, "\nYou have completed the test.\n")
			return nil
		})
		pushQuestionBankToCsv(questionBank)
		// Write response times and responses to CSV
		return os.ErrProcessDone
	}
	return nil
}

func pushQuestionBankToCsv(qData []questionData) [][]string {
	headers := []string{"Question", "Response Time (s)", "Cumulative Time (s)", "Response", "Correct Answer", "Result"}
	data := make([][]string, numQuestions)
	cumulative_time := 0.0
	for i, q := range qData {
		cumulative_time = cumulative_time + q.responseTime
		data[i] = []string{
			fmt.Sprintf("%d", i+startIdx),
			fmt.Sprintf("%f", q.responseTime),
			fmt.Sprintf("%f", cumulative_time),
			q.response,
			q.answer,
			q.correct}
	}
	if key_mode {
		writeCSV("key.csv", append([][]string{headers}, data...))
	} else {
		writeCSV("responses.csv", append([][]string{headers}, data...))
	}
	return append([][]string{headers}, data...)
}

func nextQuestion(g *gocui.Gui, v *gocui.View) error {

	if err := checkEOT(g, v); err != nil {
		if err == os.ErrProcessDone {
			return nil
		}
	}
	// Record response time for previous question
	questionBank[currIdx-startIdx].responseTime = time.Since(lastAnsTime).Seconds()

	// Record response for current question
	response := strings.TrimSpace(v.Buffer())
	if response == "" {
		response = "Not answered"
	} else if err := checkValidResponse(response); err != nil {
		if err == os.ErrInvalid {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("response")
				fmt.Fprintf(v, "\nInvalid Response")
				return nil
			})
			v.Clear()
			v.EditDelete(true)
			return nil
		}
	}
	responseView, _ := g.View("response")
	fmt.Fprintf(responseView, "\nQuestion %d  Answer: %s  Resp Time: %0f", currIdx, response, questionBank[currIdx-startIdx].responseTime)
	questionBank[currIdx-startIdx].response = response
	v.Clear()
	v.EditDelete(true)

	lastAnsTime = time.Now()
	// Move to next question or end test
	currIdx++
	if err := checkEOT(g, v); err != nil {
		if err == os.ErrProcessDone {
			return nil
		}
	}
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("question")
		v.Title = fmt.Sprintf("Question %d", currIdx)
		fmt.Fprintf(v, "This is question %d\n", currIdx)
		return nil
	})
	return nil
}

func evalResponses() {
	// var eval = make([]questionData, numQuestions)
	responses := csvToQuestionData(readCSV("responses.csv"))
	key := csvToQuestionData(readCSV("key.csv"))

	for i, question := range responses {
		if question.response == key[i].response {
			responses[i].correct = "Correct"
		} else {
			responses[i].correct = "Wrong"
		}
		responses[i].answer = key[i].response
	}
    printCSV(pushQuestionBankToCsv(responses)[1:])

}

func csvToQuestionData(responses_raw [][]string) []questionData {
	var responses = make([]questionData, numQuestions)
	for i, validOpt := range responses_raw {
		responses[i].questionNum, _ = strconv.Atoi(validOpt[0])
		responses[i].responseTime, _ = strconv.ParseFloat(validOpt[1], 32)
		responses[i].response = validOpt[3]
		responses[i].answer = validOpt[4]
		responses[i].correct = validOpt[5]
	}
	return responses
}

func readCSV(filename string) [][]string {
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

func writeCSV(filename string, data [][]string) {
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

func printCSV(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Question", "Response Time (s)", "Cumulative Time (s)", "Response", "Correct Answer", "Result"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
