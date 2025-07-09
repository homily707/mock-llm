package types

type GenerateReqInput struct {
	Text   string `json:"text"`
	Output chan any
}
