package types

import "time"

type GenerateReqInput struct {
	Text string `json:"text"`
}

type TokenizedGenerateReqInput struct {
	InputText      string
	InputIds       []int
	SamplingParams SamplingParams
	Stream         bool
}

type ReqState struct {
	Finished bool
	Output   chan any
	Input    any
	// For metrics
	CreatedTime    time.Time
	FinishedTime   time.Time
	FirstTokenTime time.Time
	LastTokenTime  time.Time
	// For streaming output
	LastOutputOffset int
}
