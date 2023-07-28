package utils

import ()

type QuestionData struct {
	QuestionNum  int
	Answer       string
	Response     string
	ResponseTime float64
	Correct      string
}

type CmdlineArgs struct {
	StartIdx       int
	StopIdx        int
	NumQuestions   int
	TestDuration_m int
	TestDuration_s int
	Key_mode       bool
	CheckResponses bool
}
