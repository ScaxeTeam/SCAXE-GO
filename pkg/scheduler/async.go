package scheduler

import (
	"sync"
)

type AsyncTask interface {
	OnRun()
	OnCompletion()
	OnProgressUpdate(progress interface{})
}

type BaseAsyncTask struct {
	result    interface{}
	resultMu  sync.RWMutex
	finished  bool
	cancelled bool
	progress  chan interface{}
}

func NewBaseAsyncTask() *BaseAsyncTask {
	return &BaseAsyncTask{
		progress: make(chan interface{}, 100),
	}
}

func (t *BaseAsyncTask) GetResult() interface{} {
	t.resultMu.RLock()
	defer t.resultMu.RUnlock()
	return t.result
}

func (t *BaseAsyncTask) SetResult(result interface{}) {
	t.resultMu.Lock()
	defer t.resultMu.Unlock()
	t.result = result
}

func (t *BaseAsyncTask) IsFinished() bool {
	return t.finished
}

func (t *BaseAsyncTask) IsCancelled() bool {
	return t.cancelled
}

func (t *BaseAsyncTask) Cancel() {
	t.cancelled = true
}

func (t *BaseAsyncTask) PublishProgress(progress interface{}) {
	select {
	case t.progress <- progress:
	default:

	}
}

func (t *BaseAsyncTask) OnCompletion() {}

func (t *BaseAsyncTask) OnProgressUpdate(progress interface{}) {}

type AsyncPool struct {
	pending    chan *asyncTaskWrapper
	results    chan *asyncTaskWrapper
	workerSize int
	wg         sync.WaitGroup
	closed     bool
	mu         sync.Mutex
}

type asyncTaskWrapper struct {
	task AsyncTask
	base *BaseAsyncTask
}

func NewAsyncPool(workers int) *AsyncPool {
	if workers < 1 {
		workers = 4
	}
	pool := &AsyncPool{
		pending:    make(chan *asyncTaskWrapper, 1000),
		results:    make(chan *asyncTaskWrapper, 1000),
		workerSize: workers,
	}
	pool.startWorkers()
	return pool
}

func (p *AsyncPool) startWorkers() {
	for i := 0; i < p.workerSize; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for wrapper := range p.pending {
				if wrapper.base != nil && wrapper.base.IsCancelled() {
					continue
				}

				wrapper.task.OnRun()
				if wrapper.base != nil {
					wrapper.base.finished = true
				}

				p.results <- wrapper
			}
		}()
	}
}

func (p *AsyncPool) SubmitTask(task AsyncTask) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}

	var base *BaseAsyncTask
	type hasBase interface {
		getBaseAsyncTask() *BaseAsyncTask
	}
	if hb, ok := task.(hasBase); ok {
		base = hb.getBaseAsyncTask()
	}

	p.pending <- &asyncTaskWrapper{task: task, base: base}
}

func (p *AsyncPool) CollectTasks() {
	for {
		select {
		case wrapper := <-p.results:

			if wrapper.base != nil {
				for {
					select {
					case prog := <-wrapper.base.progress:
						wrapper.task.OnProgressUpdate(prog)
					default:
						goto doneProgress
					}
				}
			doneProgress:
			}

			wrapper.task.OnCompletion()
		default:
			return
		}
	}
}

func (p *AsyncPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.closed {
		p.closed = true
		close(p.pending)
		p.wg.Wait()
		close(p.results)
	}
}

func (p *AsyncPool) GetSize() int {
	return p.workerSize
}

var globalAsyncPool *AsyncPool
var asyncPoolOnce sync.Once

func GetAsyncPool() *AsyncPool {
	asyncPoolOnce.Do(func() {
		globalAsyncPool = NewAsyncPool(4)
	})
	return globalAsyncPool
}

func SubmitAsync(task AsyncTask) {
	GetAsyncPool().SubmitTask(task)
}

type ClosureAsyncTask struct {
	*BaseAsyncTask
	runFn        func() interface{}
	completionFn func(result interface{})
}

func NewClosureAsyncTask(runFn func() interface{}, completionFn func(result interface{})) *ClosureAsyncTask {
	return &ClosureAsyncTask{
		BaseAsyncTask: NewBaseAsyncTask(),
		runFn:         runFn,
		completionFn:  completionFn,
	}
}

func (c *ClosureAsyncTask) OnRun() {
	if c.runFn != nil {
		c.SetResult(c.runFn())
	}
}

func (c *ClosureAsyncTask) OnCompletion() {
	if c.completionFn != nil {
		c.completionFn(c.GetResult())
	}
}

func (c *ClosureAsyncTask) getBaseAsyncTask() *BaseAsyncTask {
	return c.BaseAsyncTask
}

func RunAsync(runFn func() interface{}, completionFn func(result interface{})) {
	SubmitAsync(NewClosureAsyncTask(runFn, completionFn))
}
