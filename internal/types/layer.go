package types

type LogitsProcessorOutput struct {
	NextTokenLogits Tensor
}

type Tensor struct {
}

func (t *Tensor) Index(i int) *Tensor {
	return &Tensor{}
}

func (t *Tensor) ToInt() int {
	return 0
}
