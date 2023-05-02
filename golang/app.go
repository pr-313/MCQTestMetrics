package main

import (
	"encoding/csv"
	"fmt"
	"os"
	//"strconv"
	"time"
	"github.com/jroimartin/gocui"
)

const (
	numQuestions = 10 // Change this to the number of questions in your test
)

var (
	answers     = make([]string, numQuestions)
	responses   = make([]string, numQuestions)
	startTime   time.Time
	questionNum int
)

func main() {
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

	// Set the answer key (if available)
	// Uncomment the following line to manually set the answer key
	// setAnswerKey()

	// Start the test
	startTime = time.Now()
	questionNum = 1
	g.MainLoop()
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("question", 1, 1, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("Question %d", questionNum)
		fmt.Fprintln(v, "This is question", questionNum)
	}
	if v, err := g.SetView("response", 1, 6, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Response"
		fmt.Fprintf(v, "Enter your response for question %d here:\n", questionNum)
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		g.SetCurrentView("response")
	}
	if v, err := g.SetView("timer", maxX/2, 6, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Timer"
		fmt.Fprintf(v, "Time elapsed: %v", time.Since(startTime).Round(time.Second))
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("response", gocui.KeyEnter, gocui.ModNone, nextQuestion)
	if err != nil {
		return err
	}
	return nil
}

func nextQuestion(g *gocui.Gui, keyEvt *gocui.KeyEvent) error {
	// Record response time for previous question
	if questionNum > 1 {
		responseTime := time.Since(startTime).Seconds()
		responseTimes[questionNum-2] = fmt.Sprintf("%.2f", responseTime)
	}

	// Record response for current question
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("response")
		response := strings.TrimSpace(v.Buffer())
		if response == "" {
			response = "Not answered"
		}
		responses[questionNum-1] = response
		v.Clear()
		return nil
	})

	// Move to next question or end test
	questionNum++
	if questionNum > numQuestions {
		// End test
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("question")
			v.Clear()
			v.Title = "Test complete"
			fmt.Fprintf(v, "You have completed the test.\n")
			return nil
		})

		// Write response times and responses to CSV
		headers := []string{"Question", "Response Time (s)", "Response"}
		data := make([][]string, numQuestions)
		for i := 0; i < numQuestions; i++ {
			data[i] = []string{fmt.Sprintf("%d", i+1), responseTimes[i], responses[i]}
		}
		writeCSV("responses.csv", append([][]string{headers}, data...))

		// Compare responses to answer key
		if len(answers) > 0 {
			numCorrect := 0
			for i := 0; i < numQuestions; i++ {
				if responses[i] == answers[i] {
					numCorrect++
				}
			}

			// Write score to CSV
			score := float64(numCorrect) / float64(numQuestions) * 100
			writeCSV("score.csv", [][]string{{fmt.Sprintf("%.2f%%", score)}})
		} else {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("question")
				v.Clear()
				v.Title = "Test complete"
				fmt.Fprintf(v, "You have completed the test.\n")
				return nil
			})
		}
	} else {
		// Next question
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("question")
			v.Title = fmt.Sprintf("Question %d", questionNum)
			fmt.Fprintf(v, "This is question %d\n", questionNum)
			return nil
		})
	}
	return nil
}

func setAnswerKey() {
	// Manually set the answer key
	for i := 0; i < numQuestions; i++ {
		fmt.Printf("Enter the answer for question %d: ", i+1)
		fmt.Scanf("%s", &answers[i])
	}
	writeCSV("answer_key.csv", answers)
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
