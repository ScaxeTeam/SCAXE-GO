package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func registerEventAPI(L *lua.LState, p *Plugin) {
	mod := L.NewTable()

	mod.RawSetString("listen", L.NewFunction(func(L *lua.LState) int {
		eventName := L.CheckString(1)
		handler := L.CheckFunction(2)

		if p.eventHandlers == nil {
			p.eventHandlers = make(map[string][]*lua.LFunction)
		}
		p.eventHandlers[eventName] = append(p.eventHandlers[eventName], handler)
		return 0
	}))

	L.SetGlobal("events", mod)
}
