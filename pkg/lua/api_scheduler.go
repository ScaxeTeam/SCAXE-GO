package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func registerSchedulerAPI(L *lua.LState, p *Plugin, server ServerAPI) {
	mod := L.NewTable()

	mod.RawSetString("delayed", L.NewFunction(func(L *lua.LState) int {
		delay := int64(L.CheckNumber(1))
		callback := L.CheckFunction(2)

		p.nextTaskID++
		task := &schedulerTask{
			id:       p.nextTaskID,
			callback: callback,
			interval: 0,
			delay:    delay,
			nextRun:  server.GetCurrentTick() + delay,
			repeat:   false,
		}
		p.schedulerTasks = append(p.schedulerTasks, task)

		L.Push(lua.LNumber(task.id))
		return 1
	}))

	mod.RawSetString("repeating", L.NewFunction(func(L *lua.LState) int {
		interval := int64(L.CheckNumber(1))
		callback := L.CheckFunction(2)

		p.nextTaskID++
		task := &schedulerTask{
			id:       p.nextTaskID,
			callback: callback,
			interval: interval,
			delay:    interval,
			nextRun:  server.GetCurrentTick() + interval,
			repeat:   true,
		}
		p.schedulerTasks = append(p.schedulerTasks, task)

		L.Push(lua.LNumber(task.id))
		return 1
	}))

	mod.RawSetString("cancel", L.NewFunction(func(L *lua.LState) int {
		id := L.CheckInt(1)
		for _, task := range p.schedulerTasks {
			if task.id == id {
				task.cancel = true
				break
			}
		}
		return 0
	}))

	L.SetGlobal("scheduler", mod)
}
