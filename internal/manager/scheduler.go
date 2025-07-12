package manager

import (
	"github.com/homily707/mock-llm/internal/types"
)

type ScheduleReq struct {
	rid            string
	originInputIds []int
	samplingParams types.SamplingParams

	outputIds        []int
	finishedReason   string
	promptTokens     int
	completionTokens int
	cachedTokens     int
}

func (r *ScheduleReq) finished() bool {
	panic("not implemented")
}

type Scheduler struct {
	recvRequests chan any
	tpWorker     TpWorker

	lastBatch    *ScheduleBatch
	runningBatch *ScheduleBatch
	waitingQueue []*ScheduleReq
}

func (s *Scheduler) eventLoopNormal() {
	for {
		select {
		case req := <-s.recvRequests:
			switch req := req.(type) {
			case *TokenizedGenerateReqInput:
				r := ScheduleReq{
					rid:            req.rid,
					originInputIds: req.inputIds,
					promptTokens:   len(req.inputIds),
				}
				s.waitingQueue = append(s.waitingQueue, &r)
			}
		}
		batch := s.getNextBatchToRun()
		if batch != nil {
			result := s.runBatch(batch)
			s.processBatchResult(batch, result)
		}
		s.lastBatch = batch
	}
}

func (s *Scheduler) getNextBatchToRun() *ScheduleBatch {
	if s.lastBatch != nil {
		s.runningBatch.merge(s.lastBatch)
	}
	newBatch := s.getNewBatchPrefill()
	// prefill first
	if newBatch != nil {
		return newBatch
	}
	// run decode
	s.updateRunningBatch()
	return s.runningBatch
}

func (s *Scheduler) getNewBatchPrefill() *ScheduleBatch {
	s.sortWaitingQueue()
	adder := PrefillAdder{}
	left := []*ScheduleReq{}
	for _, req := range s.waitingQueue {
		success := adder.add(req)
		if !success {
			left = append(left, req)
		}
	}
	s.waitingQueue = left
	return NewScheduleBatch()
}

func (s *Scheduler) sortWaitingQueue() {
	// FIFO do nothing
}

func (s *Scheduler) runBatch(batch *ScheduleBatch) *ScheduleBatchResult {
	modelWorkBatch := batch.getModelWorkBatch()
	logitsOutput, nextTokenIds := s.tpWorker.ForwardBatchGeneration(modelWorkBatch)
	return &ScheduleBatchResult{
		logitsOutput: logitsOutput,
		nextTokenIds: nextTokenIds,
	}
}

func (s *Scheduler) processBatchResult(batch *ScheduleBatch, result *ScheduleBatchResult) {
	if batch.forwardMode == "EXTEND" {
		// todo: chunked

		// non chunked, just act like decode first token
		for i, req := range batch.reqs {
			req.outputIds = append(req.outputIds, result.nextTokenIds.Index(i).ToInt())
		}
	}
	if batch.forwardMode == "DECODE" {
		for i, req := range batch.reqs {
			req.outputIds = append(req.outputIds, result.nextTokenIds.Index(i).ToInt())
		}
	}

	s.streamOutput(batch.reqs)
}

func (s *Scheduler) streamOutput(reqs []*ScheduleReq) {
	for _, req := range reqs {
		if req.finished() {
		}
	}
}

func (s *Scheduler) updateRunningBatch() {
	return
}

type ScheduleBatch struct {
	forwardMode string
	reqs        []*ScheduleReq
}

func NewScheduleBatch() *ScheduleBatch {
	return &ScheduleBatch{}
}

func (s *ScheduleBatch) merge(other *ScheduleBatch) {

}

func (s *ScheduleBatch) getModelWorkBatch() *ModelWorkBatch {
	return &ModelWorkBatch{}
}

type ScheduleBatchResult struct {
	logitsOutput types.LogitsProcessorOutput
	nextTokenIds types.Tensor
}

type PrefillAdder struct {
	RemTotalTokens int
}

func NewPerfillAdder(batch *ScheduleBatch) *PrefillAdder {
	remTokens := 0
	for _, req := range batch.reqs {
		remTokens += len(req.originInputIds)
	}
	return &PrefillAdder{
		RemTotalTokens: remTokens,
	}
}

func kvCacheAvailableTokenSize() int {
	panic("not implemented")
}

func (p PrefillAdder) add(req *ScheduleReq) bool {
	totalTokens := len(req.originInputIds) + req.samplingParams.MaxNewTokens
	if p.RemTotalTokens < totalTokens {
		return false
	}
	p.RemTotalTokens -= totalTokens
	return true
}

type ModelWorkBatch struct {
}
