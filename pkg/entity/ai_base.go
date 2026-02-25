package entity
type AIGoal interface {
	ShouldExecute() bool
	ShouldContinueExecuting() bool
	IsInterruptible() bool
	StartExecuting()
	ResetTask()
	UpdateTask()
	GetMutexBits() int
}
type BaseAIGoal struct {
	mutexBits int
}

func (b *BaseAIGoal) ShouldContinueExecuting() bool {
	return false
}

func (b *BaseAIGoal) IsInterruptible() bool {
	return true
}

func (b *BaseAIGoal) StartExecuting() {}

func (b *BaseAIGoal) ResetTask() {}

func (b *BaseAIGoal) UpdateTask() {}

func (b *BaseAIGoal) SetMutexBits(bits int) {
	b.mutexBits = bits
}

func (b *BaseAIGoal) GetMutexBits() int {
	return b.mutexBits
}
