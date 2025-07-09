package manager

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/homily707/mock-llm/internal/types"
)

type TokenizerManager struct {
	reqStates map[string]*ReqState
}

func (tm *TokenizerManager) tokenizeGenerateReqInput(input *types.GenerateReqInput) TokenizedGenerateReqInput {
	rid := uuid.New().String()
	tm.reqStates[rid] = &ReqState{
		rid:      rid,
		Finished: false,
		Output:   input.Output,
		// For metrics
		CreatedTime:    time.Now(),
		FinishedTime:   time.Now(),
		FirstTokenTime: time.Now(),
		LastTokenTime:  time.Now(),
	}
	return TokenizedGenerateReqInput{}
}

type DetokenizerManager struct {
	recvIdsOutCh chan *TokenIDOut
	sendStrOutCh chan *StrOut
}

func (dm *DetokenizerManager) eventLoop(ctx context.Context) {
	for {
		select {
		case out := <-dm.recvIdsOutCh:
			dm.sendStrOutCh <- dm.detokenizeTokenIDOut(out)
		case <-ctx.Done():
			return
		}
	}
}

func (dm *DetokenizerManager) detokenizeTokenIDOut(out *TokenIDOut) *StrOut {
	// mock detokenization
	str := strings.Builder{}
	for _, id := range out.decodeIds {
		ch := rune(id%26 + 'a')
		str.WriteRune(ch)
	}

	return &StrOut{
		rid:              out.rid,
		str:              str.String(),
		finishedReason:   out.finishedReason,
		promptTokens:     out.promptTokens,
		completionTokens: out.completionTokens,
		cachedTokens:     out.cachedTokens,
	}
}

type TokenIDOut struct {
	rid              string
	decodeIds        []int
	finishedReason   string
	promptTokens     int
	completionTokens int
	cachedTokens     int
}

type StrOut struct {
	rid              string
	str              string
	finishedReason   string
	promptTokens     int
	completionTokens int
	cachedTokens     int
}

type TokenizedGenerateReqInput struct {
	InputText      string
	InputIds       []int
	SamplingParams types.SamplingParams
	Stream         bool
}

type ReqState struct {
	rid      string
	Finished bool
	Output   chan any
	// For metrics
	CreatedTime    time.Time
	FinishedTime   time.Time
	FirstTokenTime time.Time
	LastTokenTime  time.Time
}
