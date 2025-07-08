package manager

import (
	"reflect"

	"github.com/homily707/mock-llm/internal/types"
)

type ScheduleReq struct {
	OutputIds        []int
	FinishedReason   string
	PromptTokens     int
	CompletionTokens int
	CachedTokens     int
}

func (r *ScheduleReq) finished() bool {
	panic("not implemented")
}

type Scheduler struct {
	requestDispatcher map[string]func(any)
	recvRequests      chan any
	tpWorker          TpWorker

	lastBatch    *ScheduleBatch
	runningBatch *ScheduleBatch
	waitingQueue []*ScheduleReq
}

func (s *Scheduler) eventLoopNormal() {
	for {
		select {
		case req := <-s.recvRequests:
			typ := reflect.TypeOf(req)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			if _, ok := s.requestDispatcher[typ.String()]; !ok {
				panic("no dispatcher for " + typ.String())
			}
			s.requestDispatcher[typ.String()](req)
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

	}
	if batch.forwardMode == "DECODE" {
		for i, req := range batch.reqs {
			req.OutputIds = append(req.OutputIds, result.nextTokenIds.Index(i).ToInt())
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
}

func (p PrefillAdder) add(req *ScheduleReq) bool {
	return true
}

type ModelWorkBatch struct {
}
