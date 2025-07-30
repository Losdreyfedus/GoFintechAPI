package transaction

import (
	"backend_path/internal/domain"
	"sync"
	"sync/atomic"
)

type TransactionJob struct {
	Tx *domain.Transaction
	// Additional fields can be added as needed
}

type Processor struct {
	queue     chan TransactionJob
	wg        sync.WaitGroup
	workerNum int
	processed int64 // atomic counter
}

func NewProcessor(workerNum, queueSize int) *Processor {
	return &Processor{
		queue:     make(chan TransactionJob, queueSize),
		workerNum: workerNum,
	}
}

func (p *Processor) Start(processFunc func(TransactionJob)) {
	for i := 0; i < p.workerNum; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.queue {
				processFunc(job)
				atomic.AddInt64(&p.processed, 1)
			}
		}()
	}
}

func (p *Processor) Enqueue(job TransactionJob) {
	p.queue <- job
}

func (p *Processor) Stop() {
	close(p.queue)
	p.wg.Wait()
}

func (p *Processor) ProcessedCount() int64 {
	return atomic.LoadInt64(&p.processed)
}

// Batch processing example
func (p *Processor) ProcessBatch(jobs []TransactionJob, processFunc func(TransactionJob)) {
	for _, job := range jobs {
		p.Enqueue(job)
	}
	p.Stop()
}
