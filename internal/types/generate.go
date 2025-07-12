package types

type GenerateReqInput struct {
	Text           string         `json:"text"`
	SamplingParams SamplingParams `json:"sampling_params"`
	Stream         bool           `json:"stream"`

	Output chan any
}
