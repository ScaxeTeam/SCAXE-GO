package entity

type AITaskEntry struct {
	Action   AIGoal
	Priority int
	Using    bool
}

type AITasks struct {
	taskEntries      []*AITaskEntry
	executingEntries []*AITaskEntry
	tickCount        int
	tickRate         int
	disabledFlags    int
}

func NewAITasks() *AITasks {
	return &AITasks{
		taskEntries:      make([]*AITaskEntry, 0),
		executingEntries: make([]*AITaskEntry, 0),
		tickRate:         3,
	}
}

func (t *AITasks) AddTask(priority int, task AIGoal) {
	t.taskEntries = append(t.taskEntries, &AITaskEntry{
		Action:   task,
		Priority: priority,
		Using:    false,
	})
}

func (t *AITasks) RemoveTask(task AIGoal) {
	for i, entry := range t.taskEntries {
		if entry.Action == task {
			if entry.Using {
				entry.Using = false
				entry.Action.ResetTask()
				t.removeExecuting(entry)
			}
			t.taskEntries = append(t.taskEntries[:i], t.taskEntries[i+1:]...)
			return
		}
	}
}
func (t *AITasks) OnUpdateTasks() {
	t.tickCount++

	if t.tickCount%t.tickRate == 0 {
		for _, entry := range t.taskEntries {
			if entry.Using {
				if !t.canUse(entry) || !entry.Action.ShouldContinueExecuting() {
					entry.Using = false
					entry.Action.ResetTask()
					t.removeExecuting(entry)
				}
			} else {
				if t.canUse(entry) && entry.Action.ShouldExecute() {
					entry.Using = true
					entry.Action.StartExecuting()
					t.addExecuting(entry)
				}
			}
		}
	} else {
		newExecuting := make([]*AITaskEntry, 0, len(t.executingEntries))
		for _, entry := range t.executingEntries {
			if !entry.Action.ShouldContinueExecuting() {
				entry.Using = false
				entry.Action.ResetTask()
			} else {
				newExecuting = append(newExecuting, entry)
			}
		}
		t.executingEntries = newExecuting
	}

	for _, entry := range t.executingEntries {
		entry.Action.UpdateTask()
	}
}
func (t *AITasks) canUse(taskEntry *AITaskEntry) bool {
	if len(t.executingEntries) == 0 {
		return true
	}

	if t.isControlFlagDisabled(taskEntry.Action.GetMutexBits()) {
		return false
	}

	for _, executing := range t.executingEntries {
		if executing == taskEntry {
			continue
		}

		if taskEntry.Priority >= executing.Priority {
			if !t.areTasksCompatible(taskEntry, executing) {
				return false
			}
		} else {
			if !executing.Action.IsInterruptible() {
				return false
			}
		}
	}

	return true
}
func (t *AITasks) areTasksCompatible(entry1, entry2 *AITaskEntry) bool {
	return (entry1.Action.GetMutexBits() & entry2.Action.GetMutexBits()) == 0
}

func (t *AITasks) isControlFlagDisabled(flags int) bool {
	return (t.disabledFlags & flags) > 0
}

func (t *AITasks) DisableControlFlag(flag int) {
	t.disabledFlags |= flag
}

func (t *AITasks) EnableControlFlag(flag int) {
	t.disabledFlags &= ^flag
}

func (t *AITasks) addExecuting(entry *AITaskEntry) {
	for _, e := range t.executingEntries {
		if e == entry {
			return
		}
	}
	t.executingEntries = append(t.executingEntries, entry)
}

func (t *AITasks) removeExecuting(entry *AITaskEntry) {
	for i, e := range t.executingEntries {
		if e == entry {
			t.executingEntries = append(t.executingEntries[:i], t.executingEntries[i+1:]...)
			return
		}
	}
}
