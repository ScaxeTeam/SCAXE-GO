package scheduler

type Task interface {
	Name() string
	OnRun(currentTick int64)
	OnCancel()
}

type BaseTask struct {
	name    string
	handler *TaskHandler
}

func NewBaseTask(name string) *BaseTask {
	return &BaseTask{name: name}
}

func (t *BaseTask) Name() string {
	return t.name
}

func (t *BaseTask) GetHandler() *TaskHandler {
	return t.handler
}

func (t *BaseTask) SetHandler(handler *TaskHandler) {
	if t.handler == nil || handler == nil {
		t.handler = handler
	}
}

func (t *BaseTask) OnCancel() {}

type TaskHandler struct {
	task      Task
	taskID    int
	delay     int
	period    int
	nextRun   int64
	cancelled bool
}

func NewTaskHandler(task Task, taskID, delay, period int) *TaskHandler {
	return &TaskHandler{
		task:      task,
		taskID:    taskID,
		delay:     delay,
		period:    period,
		nextRun:   0,
		cancelled: false,
	}
}

func (h *TaskHandler) GetTask() Task {
	return h.task
}

func (h *TaskHandler) GetTaskID() int {
	return h.taskID
}

func (h *TaskHandler) GetDelay() int {
	return h.delay
}

func (h *TaskHandler) GetPeriod() int {
	return h.period
}

func (h *TaskHandler) GetNextRun() int64 {
	return h.nextRun
}

func (h *TaskHandler) SetNextRun(tick int64) {
	h.nextRun = tick
}

func (h *TaskHandler) IsDelayed() bool {
	return h.delay > 0
}

func (h *TaskHandler) IsRepeating() bool {
	return h.period > 0
}

func (h *TaskHandler) IsCancelled() bool {
	return h.cancelled
}

func (h *TaskHandler) Cancel() {
	h.cancelled = true
	if h.task != nil {
		h.task.OnCancel()
	}
}

func (h *TaskHandler) Remove() {
	h.cancelled = true
}

func (h *TaskHandler) Run(currentTick int64) {
	if !h.cancelled && h.task != nil {
		h.task.OnRun(currentTick)
	}
}

type Scheduler struct {
	tasks       map[int]*TaskHandler
	queue       []*TaskHandler
	nextID      int
	currentTick int64
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks:  make(map[int]*TaskHandler),
		queue:  make([]*TaskHandler, 0),
		nextID: 1,
	}
}

var globalScheduler = NewScheduler()

func GetGlobalScheduler() *Scheduler {
	return globalScheduler
}

func (s *Scheduler) ScheduleTask(task Task) *TaskHandler {
	return s.addTask(task, 0, -1)
}

func (s *Scheduler) ScheduleDelayedTask(task Task, delay int) *TaskHandler {
	return s.addTask(task, delay, -1)
}

func (s *Scheduler) ScheduleRepeatingTask(task Task, period int) *TaskHandler {
	return s.addTask(task, 0, period)
}

func (s *Scheduler) ScheduleDelayedRepeatingTask(task Task, delay, period int) *TaskHandler {
	return s.addTask(task, delay, period)
}

func (s *Scheduler) addTask(task Task, delay, period int) *TaskHandler {
	if delay < 0 {
		delay = 0
	}
	if period < 1 && period != -1 {
		period = 1
	}

	handler := NewTaskHandler(task, s.nextID, delay, period)
	s.nextID++

	if handler.IsDelayed() {
		handler.SetNextRun(s.currentTick + int64(delay))
	} else {
		handler.SetNextRun(s.currentTick)
	}

	s.tasks[handler.GetTaskID()] = handler
	s.insertIntoQueue(handler)

	return handler
}

func (s *Scheduler) insertIntoQueue(handler *TaskHandler) {

	inserted := false
	for i, h := range s.queue {
		if handler.GetNextRun() < h.GetNextRun() {

			s.queue = append(s.queue[:i], append([]*TaskHandler{handler}, s.queue[i:]...)...)
			inserted = true
			break
		}
	}
	if !inserted {
		s.queue = append(s.queue, handler)
	}
}

func (s *Scheduler) CancelTask(taskID int) {
	if handler, ok := s.tasks[taskID]; ok {
		handler.Cancel()
		delete(s.tasks, taskID)
	}
}

func (s *Scheduler) CancelAllTasks() {
	for _, handler := range s.tasks {
		handler.Cancel()
	}
	s.tasks = make(map[int]*TaskHandler)
	s.queue = make([]*TaskHandler, 0)
	s.nextID = 1
}

func (s *Scheduler) IsQueued(taskID int) bool {
	_, ok := s.tasks[taskID]
	return ok
}

func (s *Scheduler) MainThreadHeartbeat(currentTick int64) {
	s.currentTick = currentTick

	for s.isReady(currentTick) {

		handler := s.queue[0]
		s.queue = s.queue[1:]

		if handler.IsCancelled() {
			delete(s.tasks, handler.GetTaskID())
			continue
		}

		handler.Run(currentTick)

		if handler.IsRepeating() {

			handler.SetNextRun(currentTick + int64(handler.GetPeriod()))
			s.insertIntoQueue(handler)
		} else {

			handler.Remove()
			delete(s.tasks, handler.GetTaskID())
		}
	}
}

func (s *Scheduler) isReady(currentTick int64) bool {
	return len(s.queue) > 0 && s.queue[0].GetNextRun() <= currentTick
}

func (s *Scheduler) GetTaskCount() int {
	return len(s.tasks)
}

type ClosureTask struct {
	*BaseTask
	fn func(int64)
}

func NewClosureTask(name string, fn func(currentTick int64)) *ClosureTask {
	return &ClosureTask{
		BaseTask: NewBaseTask(name),
		fn:       fn,
	}
}

func (c *ClosureTask) OnRun(currentTick int64) {
	if c.fn != nil {
		c.fn(currentTick)
	}
}

func ScheduleTask(task Task) *TaskHandler {
	return globalScheduler.ScheduleTask(task)
}

func ScheduleDelayed(task Task, delay int) *TaskHandler {
	return globalScheduler.ScheduleDelayedTask(task, delay)
}

func ScheduleRepeating(task Task, period int) *TaskHandler {
	return globalScheduler.ScheduleRepeatingTask(task, period)
}

func RunLater(name string, delay int, fn func(int64)) *TaskHandler {
	return globalScheduler.ScheduleDelayedTask(NewClosureTask(name, fn), delay)
}

func RunRepeating(name string, period int, fn func(int64)) *TaskHandler {
	return globalScheduler.ScheduleRepeatingTask(NewClosureTask(name, fn), period)
}
