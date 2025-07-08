package manager

import "github.com/homily707/mock-llm/internal/types"

type TpWorker interface {
	ForwardBatchGeneration(batch *ModelWorkBatch) (types.LogitsProcessorOutput, types.Tensor)
}

type TpModelWorker struct {
}
