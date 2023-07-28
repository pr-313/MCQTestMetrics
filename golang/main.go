package main

import (
	"os"

	"github.com/jroimartin/gocui"
	"github.com/pr-313/MCQTestMetrics/src"
	"github.com/pr-313/MCQTestMetrics/utils"
)

func main() {
	src.SetupGlobalVars()

	if src.Args.CheckResponses {
		utils.EvalResponses(src.Args)
		os.Exit(0)
	}
	// Initialize gocui
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	// Set the view dimensions
	g.SetManagerFunc(src.Layout)

	// Set the keybindings
	err = src.Keybindings(g)
	if err != nil {
		panic(err)
	}

	go src.RunTimer(g)
	// Start the test
	g.MainLoop()
}
