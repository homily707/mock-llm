package manager

import "reflect"

type ScheduleReq struct {
}

type Scheduler struct {
	requestDispatcher map[string]func(any)
	recvRequests      chan any

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
			s.processBatchResult(result)
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

}

func (s *Scheduler) processBatchResult(result *ScheduleBatchResult) {

}

func (s *Scheduler) updateRunningBatch() {

}

type ScheduleBatch struct {
	batchIsFull bool
}

func NewScheduleBatch() *ScheduleBatch {
	return &ScheduleBatch{}
}

func (s *ScheduleBatch) merge(other *ScheduleBatch) {
	s.batchIsFull = s.batchIsFull || other.batchIsFull
}

type ScheduleBatchResult struct {
}

type PrefillAdder struct {
}

func (p PrefillAdder) add(req *ScheduleReq) bool {

}
