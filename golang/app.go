package main

import (
	"encoding/csv"
	"fmt"
	"os"
	//"strconv"
	// "flag"
	"github.com/jroimartin/gocui"
	"github.com/integrii/flaggy"
	"strings"
	"time"
)

type questionData struct {
	questionNum  int
	answer       string
	response     string
	responseTime float64
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
)

func main() {
	setupCmdlineArgs()
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
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
    flaggy.Int(&startIdx, "s", "startIdx", "Start from Question Index")
    flaggy.Int(&stopIdx, "e", "stopIdx", "Stop at Question Index")
    flaggy.Int(&testDuration_m, "t", "dur", "Duration of test")
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
				fmt.Fprintf(v, "Time spent on this question: %ds\n", int(time.Since(lastAnsTime).Seconds()))
				fmt.Fprintf(v, "Total Time remaining: %ds\n", (testDuration_s - int(time.Since(startTime).Seconds())))
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

		// Write response times and responses to CSV
		headers := []string{"Question", "Response Time (s)", "Cumulative Time (s)", "Response"}
		data := make([][]string, numQuestions)
		cumulative_time := 0.0
		for i := startIdx; i < stopIdx+1; i++ {
			cumulative_time = cumulative_time + questionBank[i-startIdx].responseTime
			data[i-startIdx] = []string{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%f", questionBank[i-startIdx].responseTime),
				fmt.Sprintf("%f", cumulative_time),
				questionBank[i-startIdx].response}
		}
		writeCSV("responses.csv", append([][]string{headers}, data...))

		// Compare responses to answer key
		// if len(answers) > 0 {
		// 	numCorrect := 0
		// 	for i := 0; i < numQuestions; i++ {
		// 		if responses[i] == answers[i] {
		// 			numCorrect++
		// 		}
		// 	}

		// 	// Write score to CSV
		// 	score := float64(numCorrect) / float64(numQuestions) * 100
		// 	writeCSV("score.csv", [][]string{{fmt.Sprintf("%.2f%%", score)}})
		// } else {
		// 	g.Update(func(g *gocui.Gui) error {
		// 		v, _ := g.View("question")
		// 		v.Clear()
		// 		v.Title = "Test complete"
		// 		fmt.Fprintf(v, "You have completed the test.\n")
		// 		return nil
		// 	})
		// }

		return os.ErrProcessDone
	}
	return nil
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

// func setAnswerKey(g *gocui.Gui, v *gocui.View) error {
// 	answerKey := strings.TrimSpace(v.Buffer())
// 	if answerKey == "" {
// 		return nil
// 	}
// 	answers = strings.Split(answerKey, "\n")
// 	v.Clear()

// 	// Write answer key to CSV
// 	headers := []string{"Question", "Answer"}
// 	data := make([][]string, numQuestions)
// 	for i := 0; i < numQuestions; i++ {
// 		data[i] = []string{fmt.Sprintf("%d", i+1), answers[i]}
// 	}
// 	writeCSV("answer_key.csv", append([][]string{headers}, data...))

// 	g.Update(func(g *gocui.Gui) error {
// 		v, _ := g.View("question")
// 		v.Clear()
// 		v.Title = "Answer key set"
// 		fmt.Fprintf(v, "You have set the answer key.\n")
// 		return nil
// 	})

// 	return nil
// }

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

	return records
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
