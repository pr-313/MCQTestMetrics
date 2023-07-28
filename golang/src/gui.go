package src

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/pr-313/MCQTestMetrics/utils"
)

var (
	StartTime    time.Time
	LastAnsTime  time.Time
	CurrIdx      int
	QuestionBank = make([]utils.QuestionData, 1)
	Args         utils.CmdlineArgs
)

func SetupGlobalVars() {
	QuestionBank = Args.SetupCmdlineArgs()
	CurrIdx = Args.StartIdx - 1
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func checkEOT(g *gocui.Gui, v *gocui.View) error {
	file_prefix := fmt.Sprintf("Q_Start_%d_End_%d", Args.StartIdx, Args.StopIdx)
	if CurrIdx > Args.StopIdx {
		// End test
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("question")
			v.Clear()
			v.Title = "Test complete"
			fmt.Fprintf(v, "\nYou have completed the test.\n")
			return nil
		})
		if Args.Key_mode {
			utils.PushQuestionBankToCsv(QuestionBank, fmt.Sprintf("%s_key.csv", file_prefix))
		} else {
			utils.PushQuestionBankToCsv(QuestionBank, fmt.Sprintf("%s_responses.csv", file_prefix))
		}
		// Write response times and responses to CSV
		return os.ErrProcessDone
	}
	return nil
}

func nextQuestion(g *gocui.Gui, v *gocui.View) error {

	//Wait for user to start the test
	if CurrIdx == Args.StartIdx-1 {
		CurrIdx++
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("question")
			v.Title = fmt.Sprintf("Question")
			fmt.Fprintf(v, "This is question %d\n", CurrIdx)
			return nil
		})
		StartTime = time.Now()
		LastAnsTime = StartTime
		return nil
	}
	if err := checkEOT(g, v); err != nil {
		if err == os.ErrProcessDone {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("response")
				fmt.Fprintf(v, "\nTest Complete: Press q to exit")
				return nil
			})
			return nil
		}
	}
	// Record response time for previous question
	QuestionBank[CurrIdx-Args.StartIdx].ResponseTime = time.Since(LastAnsTime).Seconds()

	// Record response for current question
	response := strings.TrimSpace(v.Buffer())
	if response == "" {
		if Args.Key_mode {
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("response")
				fmt.Fprintf(v, "\nAnswer Key answer cannot be blank")
				return nil
			})
			v.Clear()
			v.EditDelete(true)
			return nil
		} else {
			response = "Not answered"
		}
	} else if err := utils.CheckValidResponse(response); err != nil {
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
	fmt.Fprintf(responseView, "\nQuestion %d  Answer: %s  Resp Time: %0f", CurrIdx, response, QuestionBank[CurrIdx-Args.StartIdx].ResponseTime)
	QuestionBank[CurrIdx-Args.StartIdx].Response = response
	QuestionBank[CurrIdx-Args.StartIdx].QuestionNum = CurrIdx
	v.Clear()
	v.EditDelete(true)

	LastAnsTime = time.Now()
	// Move to next question or end test
	CurrIdx++
	if err := checkEOT(g, v); err != nil {
		if err == os.ErrProcessDone {
			return nil
		}
	}
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("question")
		v.Title = fmt.Sprintf("Question")
		fmt.Fprintf(v, "This is question %d\n", CurrIdx)
		return nil
	})
	return nil
}

func RunTimer(g *gocui.Gui) {
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
				if CurrIdx >= Args.StartIdx {
					if Args.Key_mode {
						fmt.Fprintf(v, "Answer Key Mode")
					} else {
						fmt.Fprintf(v, "Time spent on this question: %ds\n", int(time.Since(LastAnsTime).Seconds()))
						fmt.Fprintf(v, "Total Time remaining: %dm %ds\n",
							(Args.TestDuration_m - 1 - int(time.Since(StartTime).Minutes())),
							(Args.TestDuration_s-int(time.Since(StartTime).Seconds()))%60)
					}
				}
				return nil
			})
		}
	}
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("textBox", gocui.KeyEnter, gocui.ModNone, nextQuestion); err != nil {
		return err
	}
	return nil
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("question", 1, 1, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("Questions")
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "Press Enter to Start the Test")
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
