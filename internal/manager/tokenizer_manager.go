package manager

import "github.com/homily707/mock-llm/internal/types"

type TokenizerManager struct {
}

func (tm *TokenizerManager) tokenizeGenerateReqInput(input *types.GenerateReqInput) types.TokenizedGenerateReqInput {
	return types.TokenizedGenerateReqInput{}
}
