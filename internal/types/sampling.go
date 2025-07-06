package types

type SamplingParams struct {
	MaxNewTokens int
	Stop []string
	StopTokenIds []int
	Temperature float64
	TopP float64
	TopK int

	// need understand
	FrequencyPenalty float64
	PresencePenalty float64
	RepetitionPenalty float64
	SkipSpecialTokens bool
	SpacesBetweenTokens bool
}